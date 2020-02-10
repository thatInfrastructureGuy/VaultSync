// +build mage

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
	"github.com/magefile/mage/sh"
	"log"
)

// Teri maa ki aankh
func Build() {
	log.Println("build")
}

// Install nothing. do nothing
func InstallBuildUtils() error {
	// 1. dockle
	// 2. scopelint
	// 3. kubeval
	// 4. security check kube
	// 5. integration test `kind`
	return nil
}

// Run tests
func Test() error {
	return sh.Run("go", "test", "./...", "-tags", "CI")
}
