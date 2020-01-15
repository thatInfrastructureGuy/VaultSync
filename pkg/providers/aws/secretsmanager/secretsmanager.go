package secretsmanager

import "github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"

type SecretsManager struct{}

func (s *SecretsManager) Initializer() (err error) {
	return nil
}

func (s *SecretsManager) ListSecrets() (secretList map[string]data.SecretAttribute, err error) {
	return secretList, nil
}
