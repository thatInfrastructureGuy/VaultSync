package kubernetes

import (
	"github.com/thatInfrastructureGuy/VaultSync/azure/keyvault"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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