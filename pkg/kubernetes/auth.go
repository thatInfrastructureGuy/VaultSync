package kubernetes

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Kubeconfig contains the auth creds
type Kubeconfig struct {
	clientset *kubernetes.Clientset
}

// Authenticate authenticates client.
func (k *Kubeconfig) Authenticate() (err error) {
	var config *rest.Config
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if len(kubeconfigPath) > 0 {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		fmt.Printf("Cannot connect to Kubernetes: %v\n", err)
		return err
	}
	k.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
