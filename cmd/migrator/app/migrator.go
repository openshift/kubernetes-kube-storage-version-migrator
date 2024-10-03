package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
	migrationclient "sigs.k8s.io/kube-storage-version-migrator/pkg/clients/clientset"
	"sigs.k8s.io/kube-storage-version-migrator/pkg/controller"
)

const (
	migratorUserAgent = "storage-version-migration-migrator"
)

var (
	kubeconfigPath     = flag.String("kubeconfig", "", "absolute path to the kubeconfig file specifying the apiserver instance. If unspecified, fallback to in-cluster configuration")
	kubeAPIQPS         = flag.Float32("kube-api-qps", 40.0, "QPS to use while talking with kubernetes apiserver.")
	kubeAPIBurst       = flag.Int("kube-api-burst", 1000, "Burst to use while talking with kubernetes apiserver.")
	leaseHolderId      = flag.String("lease-holder-id", "", "lease lock holder identity name")
	leaseLockName      = flag.String("lease-lock-name", "storage-version-migration-migrator-lock", "the lease lock resource name")
	leaseLockNamespace = flag.String("lease-lock-namespace", "", "the lease lock resource namespace")
)

func NewMigratorCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "kube-storage-migrator",
		Long: `The Kubernetes storage migrator migrates resources based on the StorageVersionMigrations APIs.`,
		Run: func(cmd *cobra.Command, args []string) {

			ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer done()

			if err := Run(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				done()
				os.Exit(1)
			}
		},
	}
}

func Run(ctx context.Context) error {
	http.Handle("/metrics", promhttp.Handler())
	livenessHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
	http.HandleFunc("/healthz", livenessHandler)
	go func() {
		if err := http.ListenAndServe(":2112", nil); err != nil {
			klog.Fatal(err)
		}
		klog.Info("server exited")
	}()

	var err error
	var config *rest.Config
	if *kubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfigPath)
		if err != nil {
			log.Fatalf("Error initializing client config: %v for kubeconfig: %v", err.Error(), *kubeconfigPath)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
	}
	config.QPS = *kubeAPIQPS
	config.Burst = *kubeAPIBurst
	client, err := dynamic.NewForConfig(rest.AddUserAgent(config, migratorUserAgent))
	if err != nil {
		return err
	}
	migration, err := migrationclient.NewForConfig(config)
	if err != nil {
		return err
	}
	c := controller.NewKubeMigrator(
		client,
		migration,
	)
	lock, err := newResourceLock(config)
	if err != nil {
		return err
	}
	leaderElectionConfig := newLeaderElectionConfig(lock, config)
	leaderElectionConfig.Callbacks.OnStartedLeading = func(ctx context.Context) { c.Run(ctx) }
	leaderelection.RunOrDie(ctx, leaderElectionConfig)
	panic("unreachable")
}

func newResourceLock(config *rest.Config) (resourcelock.Interface, error) {
	if len(*leaseHolderId) == 0 {
		var i string
		if hostname, err := os.Hostname(); err != nil {
			i = string(uuid.NewUUID())
		} else {
			i = hostname + "_" + string(uuid.NewUUID())
		}
		leaseHolderId = &i
	}
	if len(*leaseLockNamespace) == 0 {
		var n string
		data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			n = "default"
			klog.Warningf("Error reading service account namespace: %v", err)
		}
		if n = strings.TrimSpace(string(data)); len(n) == 0 {
			n = "default"
		}
		leaseLockNamespace = &n
	}
	lock, err := resourcelock.NewFromKubeconfig(
		resourcelock.LeasesResourceLock,
		*leaseLockNamespace,
		*leaseLockName,
		resourcelock.ResourceLockConfig{
			Identity:      *leaseHolderId,
			EventRecorder: nil,
		},
		config,
		107*time.Second,
	)
	return lock, err
}

func newLeaderElectionConfig(lock resourcelock.Interface, config *rest.Config) leaderelection.LeaderElectionConfig {
	leaderElectionConfig := leaderelection.LeaderElectionConfig{
		Lock:            lock,
		LeaseDuration:   137 * time.Second,
		RenewDeadline:   107 * time.Second,
		RetryPeriod:     26 * time.Second,
		ReleaseOnCancel: true,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStoppedLeading: func() {
				klog.Warningf("leader election lost")
				os.Exit(0)
			},
		},
		Name: *leaseLockName,
	}
	return leaderElectionConfig
}
