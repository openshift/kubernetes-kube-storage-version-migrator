package app

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	migrationclient "sigs.k8s.io/kube-storage-version-migrator/pkg/clients/clientset"
	"sigs.k8s.io/kube-storage-version-migrator/pkg/trigger"
)

const (
	triggerUserAgent = "storage-version-migration-trigger"
)

var (
	kubeconfigPath = flag.String("kubeconfig", "", "absolute path to the kubeconfig file specifying the apiserver instance. If unspecified, fallback to in-cluster configuration")
)

func NewTriggerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "kube-storage-migrator-trigger",
		Long: `The Kubernetes storage migrator triggering controller
		detects storage version changes and creates migration requests.
		It also records the status of the storage via the storageState
		API.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := Run(wait.NeverStop); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
}

func Run(stopCh <-chan struct{}) error {
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
	migration, err := migrationclient.NewForConfig(rest.AddUserAgent(config, triggerUserAgent))
	if err != nil {
		return err
	}
	c := trigger.NewMigrationTrigger(migration)
	c.Run(stopCh)
	panic("unreachable")
}
