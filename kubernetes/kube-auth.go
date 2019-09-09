package kubernetes

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Authenticator authenticates client.
// Returns pointer to Clientset
func Authenticator() *kubernetes.Clientset {
	var config *rest.Config
	var err error
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if len(kubeconfigPath) > 0 {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		fmt.Printf("Cannot connect to Kubernetes: %v\n", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}

	return clientset
}
