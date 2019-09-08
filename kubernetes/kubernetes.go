package kubernetes

import (
	"fmt"
	"os"

	"github.com/thatInfrastructureGuy/VaultSync/azure/keyvault"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Authenticator authenticates client.
// Returns pointer to Clientset
func Authenticator() *kubernetes.Clientset {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
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

// CreateNamespace creates namespace if not already present.
// errors out if already present. Does not update.
func CreateNamespace(clientset *kubernetes.Clientset, namespace string) {
	namespaceObject := apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	// Create namespace
	_, err := clientset.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		_, err := clientset.CoreV1().Namespaces().Create(&namespaceObject)
		if err != nil {
			fmt.Println(err)
		}
	}

}

// SecretsUpdater puts secrets into specified Kubernetes Secret
// If secret name not specified; secret with same name as vault is created.
func SecretsUpdater(secretName, namespace string) {
	secret := GetSecretObject(secretName, namespace)
	clientset := Authenticator()
	CreateNamespace(clientset, namespace)

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
