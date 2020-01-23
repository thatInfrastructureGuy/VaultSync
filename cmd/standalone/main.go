package main

import (
	"log"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.8/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.8/pkg/vaultsync"
)

func main() {
	env := &data.Env{}
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
