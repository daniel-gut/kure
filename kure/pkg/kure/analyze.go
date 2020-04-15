package kure

import (
	"fmt"
	"strings"

	"github.com/daniel-gut/kure/pkg/clients"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const (
	// Genereal Standads Definition
	podNameConst        = "pod"
	deploymentNameConst = "deployment"
	stsNameConst        = "statefulset"
	allDeploymentsConst = "all"
)

// IsValidAnalyzeObj validates the user input if it's a supported resource
func IsValidAnalyzeObj(arg string) bool {

	resourceArg := strings.ToLower(arg)

	switch resourceArg {
	case deploymentNameConst:
		return true
	case stsNameConst:
		return true
	case podNameConst:
		return true
	default:
		return false
	}
}

// Analyze is handling the different arguments
func Analyze(args []string) {

	if len(args) < 2 {
		args = append(args, allDeploymentsConst)
	}

	resourceType := strings.ToLower(args[0])
	resourceName := strings.ToLower(args[1])

	switch resourceType {
	case deploymentNameConst:
		fmt.Println("Deployments not yet supported")
	case stsNameConst:
		fmt.Println("Statefulsets not yet supported")
	case podNameConst:
		_ = resourceName
		getPods()
	default:
		fmt.Errorf("Unexpected argument %v", args)
	}

}

//getPods returns all pods in the cluster
func getPods() {
	k8sconfig := clients.GetConfig()

	clientset, err := kubernetes.NewForConfig(k8sconfig)
	if err != nil {
		panic(err.Error())
	}
	_ = clientset

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}
