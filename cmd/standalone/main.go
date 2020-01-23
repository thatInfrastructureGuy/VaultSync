package main

import (
	"log"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vaultsync"
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
