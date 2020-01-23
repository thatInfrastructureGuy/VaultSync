package main

import (
	"log"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/vaultsync"
)

var env *data.Env

func main() {
	err := env.Getenv()
	if err != nil {
		log.Fatal(err)
	}
	err, destinationUpdated := vaultsync.Synchronize(env)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(destinationUpdated)
}
