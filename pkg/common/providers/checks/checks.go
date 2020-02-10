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

package checks

import (
	"strings"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
)

// CommonProviderChecks performs checks/tasks which are provider agnostic.
// These tasks act as final sanitization checkpoint before the data is traferred to Consumers.
func CommonProviderChecks(env *data.Env, originalSecretName string, sourceDate time.Time, destinationDate time.Time) (updatedSecretName string, skipUpdate bool) {
	// Set updatedName as original name
	updatedSecretName = originalSecretName
	// Check if destination keys are outdated.
	if !sourceDate.After(destinationDate) {
		skipUpdate = true
	}
	// Check if ALL hyphers should be converted to underscores
	if env.ConvertHyphensToUnderscores {
		updatedSecretName = strings.ReplaceAll(originalSecretName, "-", "_")
	}
	return updatedSecretName, skipUpdate
}
