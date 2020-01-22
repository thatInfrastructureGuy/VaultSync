package main

import (
	"errors"
	"os"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/vaultsync"
)

func (e *Env) Getenv() error {
	/////////vault.go
	var ok bool
	e.Provider, ok = os.LookupEnv("PROVIDER")
	if !ok {
		return errors.New("PROVIDER env not present")
	}
	e.VaultName, ok = os.LookupEnv("VAULT_NAME")
	if !ok {
		return errors.New("VAULT_NAME env not present")
	}
	e.ConsumerType, ok = os.LookupEnv("CONSUMER")
	if !ok {
		return errors.New("CONSUMER env var not present")
	}

	e.Namespace, ok = os.LookupEnv("SECRET_NAMESPACE")
	if !ok {
		e.Namespace = "default"
	}
	e.SecretName, ok = os.LookupEnv("SECRET_NAME")
	if !ok {
		e.SecretName = e.VaultName
	}

	// checks.go
	_, ok = os.LookupEnv("CONVERT_HYPHENS_TO_UNDERSCORES")
	if ok {
		e.ConvertHyphenToUnderscores = true
	}
	return nil
}

func main() {
	vaultsync.Synchronize()
}
