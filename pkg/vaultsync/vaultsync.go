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

package vaultsync

import (
	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/consumer"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vault"
)

func Synchronize(env *data.Env) (error, bool) {
	// Select the destination
	destination, err := consumer.SelectConsumer(env)
	if err != nil {
		return err, false
	}

	// Get lastUpdated date timestamp from consumer
	destinationlastUpdated, err := destination.GetLastUpdatedDate()
	if err != nil {
		return err, false
	}

	// Select the source
	source, err := vault.SelectProvider(env, destinationlastUpdated)
	if err != nil {
		return err, false
	}
	// Poll secrets from vault which were updated since lastUpdated value
	secretList, err := source.GetSecrets(env)
	if err != nil {
		return err, false
	}

	// Update kuberenetes secrets
	err, updatedDestination := destination.PostSecrets(secretList)
	if err != nil {
		return err, updatedDestination
	}
	return nil, updatedDestination
}
