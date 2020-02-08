/*
 * Copyright 2020 Kulkarni, Ashish <thatInfrastructureGuy@gmail.com>
 * Author: Ashish Kulkarni
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package data

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// Env: This struct is instantiated by environment variables.
type Env struct {
	Provider                    string   // [Required] Cloud providers eg: aws, azure, gcp
	VaultName                   string   // [Required] Vault from which secrets will be pulled
	ConsumerType                string   // [Optional] Consumer Name eg: kubernetes
	DeploymentList              []string // [Optional] Comma separated list of deployments which need to be restarted on secret update.
	StatefulsetList             []string // [Optional] Comma separated list of statefulsets which need to be restarted on secret update.
	SecretName                  string   // [Optional] Name of the secret to be created/updated. Defaults to vaultName value.
	Namespace                   string   // [Optional] Kubernetes namespace where the secret is created/updated.
	RefreshRate                 int      // [Optional] Rate at which check for updated secret should be done. Defaults to 60.
	ConvertHyphensToUnderscores bool     // [Optional] Converts secret keys with - to _. Eg: MY-KEY ==> MY_KEY . Defaults to false.
}

// Getenv is wrapper function which instantiates Env struct
// from Environment Variables. Some sane defaults are set here.
func (e *Env) Getenv() (err error) {
	e.Provider = getenv("PROVIDER", "")
	if len(e.Provider) == 0 {
		return errors.New("PROVIDER env not present")
	}
	e.VaultName = getenv("VAULT_NAME", "")
	if len(e.VaultName) == 0 {
		return errors.New("VAULT_NAME env not present")
	}
	e.ConsumerType = getenv("CONSUMER", "kubernetes")
	e.Namespace = getenv("SECRET_NAMESPACE", "default")
	e.SecretName = getenv("SECRET_NAME", e.VaultName)
	deployments := getenv("DEPLOYMENT_LIST", "")
	if len(deployments) > 0 {
		e.DeploymentList = strings.Split(deployments, ",")
	}
	statefulsets := getenv("STATEFULSET_LIST", "")
	if len(statefulsets) > 0 {
		e.StatefulsetList = strings.Split(statefulsets, ",")
	}
	convertHyphensToUnderscores := getenv("CONVERT_HYPHENS_TO_UNDERSCORES", "false")
	if convertHyphensToUnderscores == "true" {
		e.ConvertHyphensToUnderscores = true
	}

	refreshRate := getenv("REFRESH_RATE", "60")
	e.RefreshRate, err = strconv.Atoi(refreshRate)
	if err != nil {
		e.RefreshRate = 60
		return err
	}
	return nil
}

// getenv is the internal function which does the grunt work
// of getting environment variables. Error handling is left to the wrapper function.
// Input1: Environment Variable to look up
// Input2: Default Value to set if envVar not found.
// Output1: Returns Value of the envVar
func getenv(envVar, defaultValue string) (value string) {
	value = os.Getenv(envVar)
	if len(value) == 0 {
		value = defaultValue
	}
	return value
}
