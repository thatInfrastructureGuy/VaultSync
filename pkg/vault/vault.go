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

package vault

import (
	"errors"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vault/providers/aws/secretsmanager"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vault/providers/azure/keyvault"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vault/providers/local"
)

type Vaults interface {
	GetSecrets(env *data.Env) (map[string]data.SecretAttribute, error)
}

type Vault struct {
	Provider Vaults
}

func (v *Vault) GetSecrets(env *data.Env) (map[string]data.SecretAttribute, error) {
	return v.Provider.GetSecrets(env)
}

func SelectProvider(env *data.Env, lastUpdated time.Time) (v *Vault, err error) {
	switch env.Provider {
	case "azure":
		v = &Vault{&keyvault.Keyvault{DestinationLastUpdated: lastUpdated, VaultName: env.VaultName}}
	case "aws":
		v = &Vault{&secretsmanager.SecretsManager{DestinationLastUpdated: lastUpdated, VaultName: env.VaultName}}
	case "gcp":
		return nil, errors.New("Google Secrets Manager: Not implemented yet!")
	case "hashicorp":
		return nil, errors.New("Hashicorp Vault: Not implemented yet!")
	case "local":
		v = &Vault{&local.Local{DestinationLastUpdated: lastUpdated}}
	default:
		return nil, errors.New("Please specify valid vault provider: azure, aws. (Coming soon: gcp, hashicorp)")
	}
	return v, nil
}
