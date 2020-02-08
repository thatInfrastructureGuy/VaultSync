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

import "time"

// SecretAttribute is constructed after querying Vault for each secret.
// This is core object which is passed from providers to consumers.
type SecretAttribute struct {
	DateUpdated       time.Time // Last time the secret was updated.
	Value             string    // The value of the secret key
	MarkedForDeletion bool      // Identifier which tells consumer if to discard the secret.
}
