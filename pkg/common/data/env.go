package data

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Env struct {
	Provider                    string
	VaultName                   string
	ConsumerType                string
	DeploymentList              []string
	StatefulsetList             []string
	SecretName                  string
	Namespace                   string
	RefreshRate                 int
	ConvertHyphensToUnderscores bool
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
	e.SecretName = os.Getenv("SECRET_NAME")
	if len(e.SecretName) == 0 {
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

	deployments := os.Getenv("DEPLOYMENT_LIST")
	if len(deployments) > 0 {
		e.DeploymentList = strings.Split(deployments, ",")
	}

	statefulsets := os.Getenv("STATEFULSET_LIST")
	if len(statefulsets) > 0 {
		e.StatefulsetList = strings.Split(statefulsets, ",")
	}

	convertHyphensToUnderscores := os.Getenv("CONVERT_HYPHENS_TO_UNDERSCORES")
	if convertHyphensToUnderscores == "true" {
		e.ConvertHyphensToUnderscores = true
	}
	return nil
}
