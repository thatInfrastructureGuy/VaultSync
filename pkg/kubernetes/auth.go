package kubernetes

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type kubeconfig struct {
	clientset *kubernetes.Clientset
}

// Authenticate authenticates client.
func (k *kubeconfig) Authenticate() {
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
	k.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}
}
