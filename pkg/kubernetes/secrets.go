package kubernetes

import (
	"fmt"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var kubeSecretExists bool = true

// createSecretObject creates a Kubernetes Secret Object in memory.
func createSecretObject(secretName, namespace string) (secretObject *apiv1.Secret) {
	secretObject = &apiv1.Secret{
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

	return
}

// secretUpdater updates secrets into specified Kubernetes Secret
// If secret name not specified; secret with same name as vault is created.
// Errors out if namespace is not present.
func (k *Kubeconfig) secretUpdater(secretObject *apiv1.Secret) error {
	namespace := secretObject.GetNamespace()
	secretOut, err := k.clientset.CoreV1().Secrets(namespace).Update(secretObject)
	if err != nil {
		fmt.Println("Error updating secret: ", err)
		return err
	}
	fmt.Printf("Updated secret %q.\n", secretOut.GetObjectMeta().GetName())
	return nil
}

// secretCreator creates secrets into specified Kubernetes Secret
// If secret name not specified; secret with same name as vault is created.
// Errors out if namespace is not present.
func (k *Kubeconfig) secretCreator(secretObject *apiv1.Secret) (err error) {
	namespace := secretObject.GetNamespace()
	secretOut, err := k.clientset.CoreV1().Secrets(namespace).Create(secretObject)
	if err != nil {
		fmt.Println("Error creating secret: ", err)
		return err
	}
	fmt.Printf("Created secret %q.\n", secretOut.GetObjectMeta().GetName())
	return nil
}

// SecretsUpdater is wrapper over kube secret operations.
func (k *Kubeconfig) SecretsUpdater(secretName, namespace string, secretList map[string]data.SecretAttribute) error {
	// Get the secret object
	secretObject, err := k.clientset.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		kubeSecretExists = false

		// Create kube secret empty object
		secretObject = createSecretObject(secretName, namespace)
	}

	// Instantiate secret data
	for secretKey, secretAttributes := range secretList {
		secretObject.Data[secretKey] = []byte(secretAttributes.Value)
	}

	if !kubeSecretExists {
		return k.secretCreator(secretObject)
	}
	return k.secretUpdater(secretObject)
}