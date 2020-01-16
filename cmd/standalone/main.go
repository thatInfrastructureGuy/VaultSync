package main

import (
	"fmt"
	"os"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/consumer"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/kubernetes"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/providers/aws/secretsmanager"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/providers/azure/keyvault"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/vault"
)

func main() {
	if os.Getenv("PROVIDER") == "" || os.Getenv("VAULT_NAME") == "" || os.Getenv("SECRETS_NAMESPACE") == "" {
		fmt.Println("Required Env Vars not set, exiting...")
		os.Exit(1)
	}
	provider := os.Getenv("PROVIDER")
	vaultName := os.Getenv("VAULT_NAME")
	namespace := os.Getenv("SECRETS_NAMESPACE")

	consumerType := os.Getenv("CONSUMER")
	secretName := os.Getenv("SECRET_NAME")
	if len(secretName) == 0 {
		secretName = vaultName
	}

	var vault vault.Vaults
	var destination consumer.Consumers

	switch provider {
	case "azure":
		vault = &keyvault.Keyvault{}
	case "aws":
		vault = &secretsmanager.SecretsManager{}
	case "gcp":
		return
	case "hashicorp":
		return
	default:
		fmt.Println("Please specify valid vault provider: azure, aws, gcp, hashicorp")
		return
	}

	switch consumerType {
	case "kubernetes":
		destination = kubernetes.Config{
			SecretName: secretName,
			Namespace:  namespace,
		}
	default:
		fmt.Println("No consumer provided. Hence using kubernetes as default.")
		destination = kubernetes.Config{
			SecretName: secretName,
			Namespace:  namespace,
		}
	}

	// Poll secrets from keyvault
	secretList, err := vault.GetSecrets()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Update kuberenetes secrets
	err = destination.PostSecrets(secretList)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
