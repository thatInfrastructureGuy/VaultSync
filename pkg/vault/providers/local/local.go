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

package local

import (
	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"time"
)

type Local struct {
	DestinationLastUpdated time.Time
}

func (l *Local) GetSecrets(env *data.Env) (map[string]data.SecretAttribute, error) {
	sampleData := map[string]data.SecretAttribute{
		"key1": {
			DateUpdated: time.Now(),
			Value:       "value1",
		},
	}

	return sampleData, nil
}
