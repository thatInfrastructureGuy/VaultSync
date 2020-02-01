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

package keyvault

import (
	"time"

	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
)

type Keyvault struct {
	basicClient            keyvault.BaseClient
	DestinationLastUpdated time.Time
	VaultName              string
}

// initialize creates KeyVault instance
func (k *Keyvault) initialize() (err error) {
	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	k.logger()
	k.basicClient = keyvault.New()
	k.basicClient.Authorizer = authorizer
	return nil
}
