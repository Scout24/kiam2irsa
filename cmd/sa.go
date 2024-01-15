package cmd

import (
	"fmt"
	"github.com/Scout24/kiam2irsa/pkg/k8s"
	"github.com/Scout24/kiam2irsa/pkg/logging"
	"os"

	"github.com/spf13/cobra"
)

var saCmd = &cobra.Command{
	Use:   "sa",
	Short: "Find Kubernetes ServiceAccounts with certain annotations",
	Run: func(cmd *cobra.Command, args []string) {
		k8s.CheckAllServiceAccounts(cmd)
	},
}

func saInit() {
	sugar := logging.SugarLogger()

	// Setting a default value for kubeconfig
	homeDir, err := os.UserHomeDir()
	if err != nil {
		sugar.Error(err.Error())
		return
	}
	kubeconfig, exist := os.LookupEnv("KUBECONFIG")
	if !exist {
		kubeconfig = fmt.Sprintf("%s/.kube/config", homeDir)
	}

	saCmd.PersistentFlags().StringP("kubeconfig", "f", kubeconfig, "Full path to the kubeconfig file")
}
