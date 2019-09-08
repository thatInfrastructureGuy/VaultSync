package kubernetes

import (
	"fmt"
	"os"

	"github.com/thatInfrastructureGuy/VaultSync/azure/keyvault"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// GetSecretObject gets secrets from Azure KeyVault.
// Creates a Kubernetes Secret Object in memory.
func GetSecretObject(secretName, namespace string) apiv1.Secret {
	secret := apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{},
		Type: "Opaque",
	}

	// Poll secrets from keyvault
	secretList := keyvault.ListSecrets(keyvault.Initializer())
	for secretKey, secretAttributes := range secretList {
		secret.Data[secretKey] = []byte(secretAttributes.Value)
	}

	return secret
}

// SecretsUpdater puts secrets into specified Kubernetes Secret
// If secret name not specified; secret with same name as vault is created.
// Errors out if namespace is not present. 
// It does not create namespace as serviceaccount generally should not have permissions to create namespaces.
func SecretsUpdater(secretName, namespace string) {
	secret := GetSecretObject(secretName, namespace)
	clientset := Authenticator()
	//CreateNamespace(clientset, namespace)

	secretOut, err := clientset.CoreV1().Secrets(namespace).Update(&secret)
	if err != nil {
		fmt.Printf("Cannot update secret: %v\nCreating new secret! \n", err)
		secretOut, err = clientset.CoreV1().Secrets(namespace).Create(&secret)
		if err != nil {
			fmt.Println("Error creating secret: ", err)
		} else {
			fmt.Printf("Created secret %q.\n", secretOut.GetObjectMeta().GetName())
		}
	} else {
		fmt.Printf("Updated secret %q.\n", secretOut.GetObjectMeta().GetName())
	}
}
