package main

import (
	"fmt"
	"os"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/consumer"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/kubernetes"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/providers/aws/secretsmanager"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/providers/azure/keyvault"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/providers/local"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/vault"
)

func main() {
	if os.Getenv("PROVIDER") == "" || os.Getenv("VAULT_NAME") == "" || os.Getenv("SECRETS_NAMESPACE") == "" {
		fmt.Println("Required Env Vars not set, exiting...")
		os.Exit(1)
	}

	// Get lastUpdated date timestamp from consumer
	destination := selectConsumer()
	destinationlastUpdated, err := destination.GetLastUpdatedDate()
	if err != nil {
		fmt.Println(err)
	}

	// Poll secrets from vault which were updated since lastUpdated value
	vault := selectProvider(destinationlastUpdated)
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

func selectProvider(lastUpdated time.Time) (vaultInstance vault.Vault) {
	provider := os.Getenv("PROVIDER")
	switch provider {
	case "azure":
		vaultInstance = vault.Vault{&keyvault.Keyvault{DestinationLastUpdated: lastUpdated}}
	case "aws":
		vaultInstance = vault.Vault{&secretsmanager.SecretsManager{DestinationLastUpdated: lastUpdated}}
	case "gcp":
		os.Exit(1)
	case "hashicorp":
		os.Exit(1)
	case "local":
		vaultInstance = vault.Vault{&local.Local{DestinationLastUpdated: lastUpdated}}
	default:
		fmt.Println("Please specify valid vault provider: azure, aws, gcp, hashicorp")
		os.Exit(1)
	}
	return vaultInstance
}

func selectConsumer() (destination consumer.Consumer) {
	consumerType := os.Getenv("CONSUMER")
	switch consumerType {
	case "kubernetes":
		namespace := os.Getenv("SECRETS_NAMESPACE")
		vaultName := os.Getenv("VAULT_NAME")
		secretName := os.Getenv("SECRET_NAME")
		if len(secretName) == 0 {
			secretName = vaultName
		}
		destination = consumer.Consumer{&kubernetes.Config{
			SecretName: secretName,
			Namespace:  namespace,
		}}
	default:
		fmt.Println("No consumer provided.")
		os.Exit(1)
	}

	return destination
}
