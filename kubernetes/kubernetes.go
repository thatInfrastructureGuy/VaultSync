package main

import (
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/play/.kube/config", "kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("The kubeconfig cannot be loaded: %v\n", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
	}
	pod, err := clientset.CoreV1().Pods("pintdb").Get("internaldb-0", metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pod)
}
