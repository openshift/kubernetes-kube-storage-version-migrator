package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	migrationclient "sigs.k8s.io/kube-storage-version-migrator/pkg/clients/clientset"
	"sigs.k8s.io/kube-storage-version-migrator/pkg/controller"
)

const (
	migratorUserAgent = "storage-version-migration-migrator"
)

var (
	kubeconfigPath = flag.String("kubeconfig", "", "absolute path to the kubeconfig file specifying the apiserver instance. If unspecified, fallback to in-cluster configuration")
	kubeAPIQPS     = flag.Float32("kube-api-qps", 40.0, "QPS to use while talking with kubernetes apiserver.")
	kubeAPIBurst   = flag.Int("kube-api-burst", 1000, "Burst to use while talking with kubernetes apiserver.")
)

func NewMigratorCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "kube-storage-migrator",
		Long: `The Kubernetes storage migrator migrates resources based on the StorageVersionMigrations APIs.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer done() // give leader election a chance to release the lock

			if err := Run(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				done() // os.Exit does not call deferred functions
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
	go func() { http.ListenAndServe(":2112", nil) }()

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
	dynamic, err := dynamic.NewForConfig(rest.AddUserAgent(config, migratorUserAgent))
	if err != nil {
		return err
	}
	migration, err := migrationclient.NewForConfig(config)
	if err != nil {
		return err
	}
	c := controller.NewKubeMigrator(
		dynamic,
		migration,
	)
	if *leaderElectionEnabled {
		return runWithLeaderElection(ctx, config, c)
	}
	c.Run(ctx)
	return nil // reachable if signal cancels ctx
}
