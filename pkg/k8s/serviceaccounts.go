package k8s

import (
	"context"
	"regexp"

	"github.com/Scout24/kiam2irsa/pkg/logging"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CheckAllServiceAccounts(cmd *cobra.Command) {
	sugar := logging.SugarLogger()

	kubeconfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		sugar.Error(err.Error())
		return
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		sugar.Panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		sugar.Panic(err.Error())
	}

	serviceAccounts, err := clientset.CoreV1().ServiceAccounts("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		sugar.Panic(err.Error())
	}

	for _, sa := range serviceAccounts.Items {
		name := sa.Name
		ns := sa.Namespace
		annotations := sa.Annotations
		hasFavorable := false
		hasUndesirable := false
		for annoKey, annoValue := range annotations {
			if annoKey == RoleArnAnnotationName {
				hasFavorable = true
			}
			if annoKey == RegionalStsAnnotationName && annoValue == RegionalStsAnnotationValue {
				hasUndesirable = true
			}
		}
		if hasFavorable && !hasUndesirable {
			sugar.Infof("Service account %s in namespace %s is not yet migrated to IRSA", name, ns)
		}
	}
}

func HasServiceAccountAnnotationForIRSA(name string, namespace string, saList *v1.ServiceAccountList) (bool, error) {
	hasRoleArn := false

	for _, sa := range saList.Items {
		if name == sa.Name && namespace == sa.Namespace {
			for key, val := range sa.Annotations {
				matchVal, _ := regexp.Match("arn:aws:iam::\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d\\d:role/", []byte(val))
				if key == RoleArnAnnotationName && matchVal {
					hasRoleArn = true
				}
			}
		}
	}

	if hasRoleArn {
		return true, nil
	}
	return false, nil
}

func GetAllServiceAccounts(clientset *kubernetes.Clientset) (*v1.ServiceAccountList, error) {
	sugar := logging.SugarLogger()
	serviceAccounts, err := clientset.CoreV1().ServiceAccounts("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		sugar.Panic(err.Error())
	}

	return serviceAccounts, err
}
