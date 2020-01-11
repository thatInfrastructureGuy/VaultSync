package keyvault

import "github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"

type Source interface {
	Initializer()
	ListSecrets() map[string]data.SecretAttribute
}
