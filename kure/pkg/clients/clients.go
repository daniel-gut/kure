package clients

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// GetConfig provides a configuration based on the environment
func GetConfig() *rest.Config {

	env := GetEnv()
	switch env {
	case "non-k8s":
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		return config
	case "k8s":
		// creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		return config
	default:
		fmt.Errorf("Unknown error \n")
	}
	return nil
}

// GetEnv returns the environment as string the client running in; "non-k8s" or "k8s"
func GetEnv() string {
	_, err := net.LookupHost("kube-dns.kube-system")
	if err != nil {
		fmt.Println("Assuming running outside K8S Cluster")
		env := "non-k8s"
		return env
	}
	fmt.Println("Assuming running inside K8S Cluster")
	env := "k8s"
	return env

}

// GetNamespaceFromKubeconfig fetches current namespace from Kubeconfig
func GetNamespaceFromKubeconfig() (string, error) {

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	namespace, _, err := config.Namespace()
	if err != nil {
		return "", err
	}

	return namespace, err
}
