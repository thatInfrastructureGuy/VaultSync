package vault

import "github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"

type Vaults interface {
	GetSecrets() (map[string]data.SecretAttribute, error)
}

type Vault struct {
	Provider Vaults
}

func (v *Vault) GetSecrets() (map[string]data.SecretAttribute, error) {
	return v.Provider.GetSecrets()
}
