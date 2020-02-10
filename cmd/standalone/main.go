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

package main

import (
	"log"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/consumer/kubernetes"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/posthook"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vaultsync"
)

var (
	CodeVersion string
	BuildTime   string
	GoVersion   string
)

func main() {
	log.Printf("[VaultSync] CodeVersion = [%s], BuildTime = [%s], GoVersion = [%s]", CodeVersion, BuildTime, GoVersion)
	env := &data.Env{}
	err := env.Getenv()
	if err != nil {
		log.Fatal(err)
	}
	PeriodicSynchronize(env)
}

func PeriodicSynchronize(env *data.Env) {
	for {
		err, destinationUpdated := vaultsync.Synchronize(env)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(destinationUpdated)
		if env.RefreshRate == 0 {
			break
		}
		if destinationUpdated {
			p := &posthook.PostHook{&kubernetes.Config{
				Namespace:       env.Namespace,
				DeploymentList:  env.DeploymentList,
				StatefulsetList: env.StatefulsetList,
			}}
			err = p.PostExec()
			if err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(time.Duration(env.RefreshRate) * time.Second)
	}
}
