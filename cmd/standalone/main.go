package main

import (
	"fmt"
	"os"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.0/pkg/azure/keyvault"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.0/pkg/kubernetes"
)

var vaultName, namespace, secretName string

func main() {
	vaultName = os.Getenv("KVAULT")
	if len(vaultName) == 0 {
		fmt.Println("Azure KeyVault not provided. Exiting!")
		os.Exit(1)
	}
	secretName = os.Getenv("SECRET_NAME")
	if len(secretName) == 0 {
		secretName = vaultName
	}
	namespace = os.Getenv("SECRETS_NAMESPACE")
	if len(namespace) == 0 {
		fmt.Println("Namespace not provided. Exiting!")
		os.Exit(1)
	}

	// Poll secrets from keyvault
	secretsList := keyvault.ListSecrets(keyvault.Initializer())

	// Update kuberenetes secrets
	kubernetes.Authenticate().SecretsUpdater(secretsList, secretName, namespace)
}
