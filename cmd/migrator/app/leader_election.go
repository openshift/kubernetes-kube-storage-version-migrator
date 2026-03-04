package app

import (
	"context"
	"fmt"
	"os"
	"time"

	flag "github.com/spf13/pflag"
	"sigs.k8s.io/kube-storage-version-migrator/pkg/controller"

	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
)

var (
	leaderElectionEnabled = flag.Bool("leader-election", false, "enable leader election.")
	leaseHolderId         = flag.String("lease-holder-id", "", "lease lock holder identity name")
	leaseLockName         = flag.String("lease-lock-name", "storage-version-migration-migrator-lock", "the lease lock resource name")
	leaseLockNamespace    = flag.String("lease-lock-namespace", "", "the lease lock resource namespace")
	leaseDuration         = flag.Duration("lease-duration", 137*time.Second, "how long to wait before forcefully attempting to acquire lock")
	leaseRenewDeadline    = flag.Duration("lease-renew-deadline", 107*time.Second, "how long to wait before giving up trying to refresh a lease")
	leaseRetryPeriod      = flag.Duration("lease-retry-period", 26*time.Second, "how long to wait between any lease actions")
)

func newResourceLock(config *rest.Config) (resourcelock.Interface, error) {
	if len(*leaseHolderId) == 0 {
		var i string
		if i = os.Getenv("POD_NAME"); i == "" {
			if hostname, err := os.Hostname(); err != nil {
				i = string(uuid.NewUUID())
			} else {
				i = hostname + "_" + string(uuid.NewUUID())
			}
		}
		leaseHolderId = &i
	}
	if len(*leaseLockNamespace) == 0 {
		n, err := determineLeaseLockNamespace()
		if err != nil {
			return nil, fmt.Errorf("error determining lease lock namespace: %v", err)
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
		*leaseRenewDeadline,
	)
	return lock, err
}

func determineLeaseLockNamespace() (string, error) {
	// cli flag overrides all
	if len(*leaseLockNamespace) > 0 {
		return *leaseLockNamespace, nil
	}
	// use the pod namespace from the env if available
	if ns := os.Getenv("POD_NAMESPACE"); len(ns) > 0 {
		return ns, nil
	}
	return "", fmt.Errorf("lease lock namespace must be provided explicitly (--lease-lock-namespace) or via POD_NAMESPACE env var")
}

func newLeaderElectionConfig(lock resourcelock.Interface) leaderelection.LeaderElectionConfig {
	leaderElectionConfig := leaderelection.LeaderElectionConfig{
		Lock:            lock,
		LeaseDuration:   *leaseDuration,
		RenewDeadline:   *leaseRenewDeadline,
		RetryPeriod:     *leaseRetryPeriod,
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

func runWithLeaderElection(ctx context.Context, config *rest.Config, c *controller.KubeMigrator) error {
	lock, err := newResourceLock(config)
	if err != nil {
		return err
	}
	leaderElectionConfig := newLeaderElectionConfig(lock)
	leaderElectionConfig.Callbacks.OnStartedLeading = func(ctx context.Context) { c.Run(ctx) }
	leaderelection.RunOrDie(ctx, leaderElectionConfig)
	return nil
}
