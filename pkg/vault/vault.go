package vault

import (
	"errors"
	"os"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/vault/providers/aws/secretsmanager"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/vault/providers/azure/keyvault"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/vault/providers/local"
)

type Vaults interface {
	GetSecrets() (map[string]data.SecretAttribute, error)
}

type Vault struct {
	Provider Vaults
}

func (v *Vault) GetSecrets() (map[string]data.SecretAttribute, error) {
	return v.Provider.GetSecrets()
}

func SelectProvider(lastUpdated time.Time) (v *Vault, err error) {
	provider, ok := os.LookupEnv("PROVIDER")
	if !ok {
		return nil, errors.New("PROVIDER env not present")
	}
	vaultName, ok := os.LookupEnv("VAULT_NAME")
	if !ok {
		return nil, errors.New("VAULT_NAME env not present")
	}
	switch provider {
	case "azure":
		v = &Vault{&keyvault.Keyvault{DestinationLastUpdated: lastUpdated, VaultName: vaultName}}
	case "aws":
		v = &Vault{&secretsmanager.SecretsManager{DestinationLastUpdated: lastUpdated, VaultName: vaultName}}
	case "gcp":
		return nil, errors.New("Google Secrets Manager: Not implemented yet!")
	case "hashicorp":
		return nil, errors.New("Hashicorp Vault: Not implemented yet!")
	case "local":
		v = &Vault{&local.Local{DestinationLastUpdated: lastUpdated}}
	default:
		return nil, errors.New("Please specify valid vault provider: azure, aws. (Coming soon: gcp, hashicorp)")
	}
	return v, nil
}
