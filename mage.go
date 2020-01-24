// +build mage

package main

import (
	"fmt"
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
