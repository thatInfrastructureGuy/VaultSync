package keyvault

import (
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"
)

type Keyvault struct {
	basicClient keyvault.BaseClient
}

type Source interface {
	Initializer()
	ListSecrets() map[string]data.SecretAttribute
}
