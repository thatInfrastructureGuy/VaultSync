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

package consumer

import (
	"errors"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/consumer/kubernetes"
)

type Consumers interface {
	GetLastUpdatedDate() (date time.Time, err error)
	PostSecrets(secretList map[string]data.SecretAttribute) error
}

type Consumer struct {
	Destination Consumers
}

func (c *Consumer) PostSecrets(secretList map[string]data.SecretAttribute) (err error, updatedDestination bool) {
	if len(secretList) == 0 {
		return nil, false
	}
	return c.Destination.PostSecrets(secretList), true
}

func (c *Consumer) GetLastUpdatedDate() (date time.Time, err error) {
	return c.Destination.GetLastUpdatedDate()
}

func SelectConsumer(env *data.Env) (c *Consumer, err error) {
	switch env.ConsumerType {
	case "kubernetes":
		if env.SecretName == "" {
			return nil, errors.New("Invalid secret name!")
		}
		c = &Consumer{&kubernetes.Config{
			SecretName: env.SecretName,
			Namespace:  env.Namespace,
		}}
	default:
		return nil, errors.New("No consumer provided.")
	}
	return c, nil
}
