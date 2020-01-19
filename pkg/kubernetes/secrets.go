package kubernetes

import (
	"log"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
		Type: "Opaque",
	}

	return
}

// secretUpdater updates secrets into specified Kubernetes Secret
// If secret name not specified; secret with same name as vault is created.
// Errors out if namespace is not present.
func (k *Config) secretUpdater(secretObject *apiv1.Secret) error {
	namespace := secretObject.GetNamespace()
	secretOut, err := k.clientset.CoreV1().Secrets(namespace).Update(secretObject)
	if err != nil {
		log.Println("Error updating secret: ", err)
		return err
	}
	log.Printf("Updated secret %q.\n", secretOut.GetObjectMeta().GetName())
	return nil
}

// secretCreator creates secrets into specified Kubernetes Secret
// If secret name not specified; secret with same name as vault is created.
// Errors out if namespace is not present.
func (k *Config) secretCreator(secretObject *apiv1.Secret) (err error) {
	namespace := secretObject.GetNamespace()
	secretOut, err := k.clientset.CoreV1().Secrets(namespace).Create(secretObject)
	if err != nil {
		log.Println("Error creating secret: ", err)
		return err
	}
	log.Printf("Created secret %q.\n", secretOut.GetObjectMeta().GetName())
	return nil
}

// secretsUpdater is an internal wrapper over kube secret operations.
func (k *Config) secretsUpdater(secretList map[string]data.SecretAttribute) error {
	secretObject, kubeSecretExists := k.getSecretObject()
	// Instantiate secret data
	if len(secretObject.Data) == 0 {
		secretObject.Data = make(map[string][]byte)
	}
	for secretKey, secretAttributes := range secretList {
		secretObject.Data[secretKey] = []byte(secretAttributes.Value)
		//Delete the key if set to empty
		if secretAttributes.Value == "" || secretAttributes.MarkedForDeletion {
			delete(secretObject.Data, secretKey)
		}
	}
	// Set the date updated timestamp
	annotations := secretObject.GetAnnotations()
	if len(annotations) == 0 {
		annotations = make(map[string]string)
	}
	annotations["dateUpdated"] = time.Now().Format(time.RFC3339)
	secretObject.SetAnnotations(annotations)

	if !kubeSecretExists {
		return k.secretCreator(secretObject)
	}
	return k.secretUpdater(secretObject)
}

func (k *Config) getSecretObject() (secretObject *apiv1.Secret, kubeSecretExists bool) {
	// Get the secret object
	secretObject, err := k.clientset.CoreV1().Secrets(k.Namespace).Get(k.SecretName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		kubeSecretExists = false

		// Create kube secret empty object
		secretObject = createSecretObject(k.SecretName, k.Namespace)
		return
	}
	kubeSecretExists = true
	return
}
func (k *Config) GetLastUpdatedDate() (date time.Time, err error) {
	err = k.authenticate()
	if err != nil {
		return date, err
	}
	secretObject, kubeSecretExists := k.getSecretObject()
	if !kubeSecretExists {
		return date, nil
	}
	annotations := secretObject.GetAnnotations()
	value, ok := annotations["dateUpdated"]
	if !ok {
		return date, nil
	}
	date, err = time.Parse(time.RFC3339, value)
	if err != nil {
		return date, err
	}
	return date, nil
}

// PostSecrets is common interface function to post secrets to destination
func (k *Config) PostSecrets(secretList map[string]data.SecretAttribute) (err error) {
	err = k.authenticate()
	if err != nil {
		return err
	}
	err = k.secretsUpdater(secretList)
	if err != nil {
		return err
	}
	return nil
}
