package kubernetes

import (
	"fmt"
)

// SecretsUpdater puts secrets into specified Kubernetes Secret
// If secret name not specified; secret with same name as vault is created.
// Errors out if namespace is not present.
// It does not create namespace as serviceaccount generally should not have permissions to create namespaces.
func SecretsUpdater(secretName, namespace string) {
	secret := GetSecretObject(secretName, namespace)
	clientset := Authenticator()

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
