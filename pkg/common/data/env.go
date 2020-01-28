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
	e.Provider = os.Getenv("PROVIDER")
	if len(e.Provider) == 0 {
		return errors.New("PROVIDER env not present")
	}
	e.VaultName = os.Getenv("VAULT_NAME")
	if len(e.VaultName) == 0 {
		return errors.New("VAULT_NAME env not present")
	}
	e.ConsumerType = os.Getenv("CONSUMER")
	if len(e.ConsumerType) == 0 {
		return errors.New("CONSUMER env var not present")
	}
	e.Namespace = os.Getenv("SECRET_NAMESPACE")
	if len(e.Namespace) == 0 {
		e.Namespace = "default"
	}
	e.SecretName = os.Getenv("SECRET_NAME")
	if len(e.SecretName) == 0 {
		e.SecretName = e.VaultName
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

	e.RefreshRate = 60
	refreshRate := os.Getenv("REFRESH_RATE")
	if len(refreshRate) > 0 {
		e.RefreshRate, err = strconv.Atoi(refreshRate)
		if err != nil {
			return err
		}
	}
	return nil
}
