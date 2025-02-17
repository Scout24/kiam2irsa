package k8s

import (
	"strings"

	"github.com/Scout24/kiam2irsa/pkg/logging"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	RoleAnnotation             string = "iam.amazonaws.com/role"
	RoleArnAnnotationName      string = "eks.amazonaws.com/role-arn"
	RegionalStsAnnotationName  string = "eks.amazonaws.com/sts-regional-endpoints"
	RegionalStsAnnotationValue string = "true"
)

func getFlag(cmd *cobra.Command, name string) (string, error) {
	sugar := logging.SugarLogger()
	status, err := cmd.Flags().GetString(name)
	if err != nil {
		sugar.Error(err.Error())
		return "", err
	}
	return strings.ToUpper(status), nil
}

func k8sClientSet(cmd *cobra.Command) (*kubernetes.Clientset, error) {
	sugar := logging.SugarLogger()

	kubeconfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		sugar.Error(err.Error())
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		sugar.Error(err.Error())
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		sugar.Error(err.Error())
		return nil, err
	}
	return clientset, nil
}
