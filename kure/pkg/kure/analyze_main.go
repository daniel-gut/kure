package kure

import (
	"fmt"
	"strings"

	"github.com/daniel-gut/kure/pkg/clients"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
)

const (
	// Genereal Standads Definition
	podNameConst        = "pod"
	deploymentNameConst = "deployment"
	stsNameConst        = "statefulset"
	allResourcesConst   = "all"
)

var (
	k8sConfig *rest.Config
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
func Analyze(args []string) error {

	var (
		podList []string
	)

	if len(args) < 2 {
		args = append(args, allResourcesConst)
	}

	//  initialize Kubeconfig
	k8sConfig = clients.GetConfig()

	resourceType := strings.ToLower(args[0])
	resourceName := strings.ToLower(args[1])

	switch resourceType {
	case deploymentNameConst:
		if resourceName == allResourcesConst {
			return fmt.Errorf("error: no Deployment name provided. All deployments not yet supported")
		}

		podList, err := getDeploymentsPods(resourceName)
		if err != nil {
			return fmt.Errorf("couldn't fetch the pods of the deployment: %w", err)
		}

		err = analyzeLog(podList)
		if err != nil {
			return fmt.Errorf("error during log analysis analyzeLog(): %w", err)
		}

	case stsNameConst:
		if resourceName == allResourcesConst {
			return fmt.Errorf("error: no Statefulset name provided. All Statefulset not yet supported")
		}

		podList, err := getStatefulsetPods(resourceName)
		if err != nil {
			return fmt.Errorf("couldn't fetch the pods of the statefulset: %w", err)
		}

		err = analyzeLog(podList)
		if err != nil {
			return fmt.Errorf("error during log analysis analyzeLog(): %w", err)
		}
	case podNameConst:

		if resourceName == allResourcesConst {
			return fmt.Errorf("error: no Pod name provided")
		}
		podList = append(podList, resourceName)

		err := analyzeLog(podList)

		if err != nil {
			return fmt.Errorf("error during log analysis analyzeLog(): %w", err)
		}

	default:
		return fmt.Errorf("Unexpected argument %v", args)
	}
	return nil
}

//getPods returns all pods in the cluster
func getPods() {

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}

func getDeploymentsPods(resourceName string) ([]string, error) {

	var (
		podLabelSelector string
		podList          []string
	)

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		panic(err.Error())
	}

	namespace, err := clients.GetNamespaceFromKubeconfig()
	if err != nil {
		namespace = ""
	}

	deployment, err := clientset.AppsV1().Deployments(namespace).List(metav1.ListOptions{FieldSelector: "metadata.name=" + resourceName})
	if err != nil {
		return nil, err
	}
	if len(deployment.Items) == 0 {
		return nil, fmt.Errorf("no deployment found with name: %s", resourceName)
	}

	for _, d := range deployment.Items {
		podLabelSelector = labels.Set(d.Spec.Selector.MatchLabels).String()
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: podLabelSelector})
	if err != nil {
		return nil, err
	}

	for _, p := range pods.Items {
		podList = append(podList, p.ObjectMeta.Name)
	}

	return podList, nil
}

func getStatefulsetPods(resourceName string) ([]string, error) {

	var (
		podLabelSelector string
		podList          []string
	)

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		panic(err.Error())
	}

	namespace, err := clients.GetNamespaceFromKubeconfig()
	if err != nil {
		namespace = ""
	}

	sts, err := clientset.AppsV1().StatefulSets(namespace).List(metav1.ListOptions{FieldSelector: "metadata.name=" + resourceName})
	if err != nil {
		return nil, err
	}
	if len(sts.Items) == 0 {
		return nil, fmt.Errorf("no statefulset found with name: %s", resourceName)
	}

	for _, sts := range sts.Items {
		podLabelSelector = labels.Set(sts.Spec.Selector.MatchLabels).String()
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: podLabelSelector})
	if err != nil {
		return nil, err
	}

	for _, p := range pods.Items {
		podList = append(podList, p.ObjectMeta.Name)
	}

	return podList, nil
}
