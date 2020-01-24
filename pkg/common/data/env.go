package data

import (
	"errors"
	"os"
	"strconv"
)

type Env struct {
	Provider                   string
	VaultName                  string
	ConsumerType               string
	Namespace                  string
	SecretName                 string
	RefreshRate                int
	ConvertHyphenToUnderscores bool
}

func (e *Env) Getenv() (err error) {
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

	refreshRate, ok := os.LookupEnv("REFRESH_RATE")
	e.RefreshRate = 60
	if ok {
		e.RefreshRate, err = strconv.Atoi(refreshRate)
		if err != nil {
			return err
		}
	}

	_, ok = os.LookupEnv("CONVERT_HYPHENS_TO_UNDERSCORES")
	if ok {
		e.ConvertHyphenToUnderscores = true
	}
	return nil
}
