package main

import (
	"fmt"
	"os"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/azure/keyvault"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/kubernetes"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/vault"
)

var vaultName, namespace, secretName string

func main() {
	vaultName = os.Getenv("VAULT_NAME")
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
	var azure vault.Vaults = &keyvault.Keyvault{}
	err := azure.Initializer()
	if err != nil {
		fmt.Println(err)
	}
	secretList := azure.ListSecrets()

	// Update kuberenetes secrets
	var destination kubernetes.Destination = kubernetes.Config{
		SecretName: secretName,
		Namespace:  namespace,
	}

	// Use destination interface methods
	err = destination.Authenticate()
	if err != nil {
		fmt.Println(err)
	}
	err = destination.SecretsUpdater(secretList)
	if err != nil {
		fmt.Println(err)
	}
}
