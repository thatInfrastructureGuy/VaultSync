package consumer

import "github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"

type Consumers interface {
	PostSecrets(secretList map[string]data.SecretAttribute) error
}

type Consumer struct {
	Destination Consumers
}

func (c *Consumer) PostSecrets(secretList map[string]data.SecretAttribute) (err error) {
	return c.Destination.PostSecrets(secretList)
}
