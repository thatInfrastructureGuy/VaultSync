package vault

import "github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"

type Vaults interface {
	Initializer() error
	ListSecrets() (map[string]data.SecretAttribute, error)
}

type Vault struct {
	Provider Vaults
}

func (v *Vault) Initialize() error {
	return v.Provider.Initializer()
}

func (v *Vault) GetSecrets() (map[string]data.SecretAttribute, error) {
	return v.Provider.ListSecrets()
}
