package main

import (
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/vaultsync"
)

var d data.Env

func main() {
	d.Getenv()
	vaultsync.Synchronize()
}
