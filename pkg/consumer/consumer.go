package consumer

import (
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"
)

type Consumers interface {
	GetLastUpdatedDate() (date time.Time, err error)
	PostSecrets(secretList map[string]data.SecretAttribute) error
}

type Consumer struct {
	Destination Consumers
}

func (c *Consumer) PostSecrets(secretList map[string]data.SecretAttribute) (err error) {
	return c.Destination.PostSecrets(secretList)
}

func (c *Consumer) GetLastUpdatedDate() (date time.Time, err error) {
	return c.Destination.GetLastUpdatedDate()
}
