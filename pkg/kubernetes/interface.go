package kubernetes

import "github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"

type Destination interface {
	Authenticate() error
	SecretsUpdater(secretList map[string]data.SecretAttribute) error
}
