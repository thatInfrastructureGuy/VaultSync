package consumer

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/consumer/kubernetes"
)

type Consumers interface {
	GetLastUpdatedDate() (date time.Time, err error)
	PostSecrets(secretList map[string]data.SecretAttribute) error
}

type Consumer struct {
	Destination Consumers
}

func (c *Consumer) PostSecrets(secretList map[string]data.SecretAttribute) (err error) {
	if len(secretList) == 0 {
		log.Println("Nothing to update!")
		return nil
	}
	return c.Destination.PostSecrets(secretList)
}

func (c *Consumer) GetLastUpdatedDate() (date time.Time, err error) {
	return c.Destination.GetLastUpdatedDate()
}

func SelectConsumer() (c *Consumer, err error) {
	consumerType, ok := os.LookupEnv("CONSUMER")
	if !ok {
		return nil, errors.New("CONSUMER env var not present")
	}
	vaultName := os.Getenv("VAULT_NAME")
	switch consumerType {
	case "kubernetes":
		namespace, ok := os.LookupEnv("SECRET_NAMESPACE")
		if !ok {
			namespace = "default"
		}
		secretName, ok := os.LookupEnv("SECRET_NAME")
		if !ok {
			secretName = vaultName
		}
		if secretName == "" {
			return nil, errors.New("Invalid secret name!")
		}
		c = &Consumer{&kubernetes.Config{
			SecretName: secretName,
			Namespace:  namespace,
		}}
	default:
		return nil, errors.New("No consumer provided.")
	}
	return c, nil
}
