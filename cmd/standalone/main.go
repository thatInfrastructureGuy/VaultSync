package main

import (
	"log"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vaultsync"
)

func main() {
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
		time.Sleep(time.Duration(env.RefreshRate) * time.Second)
	}
}
