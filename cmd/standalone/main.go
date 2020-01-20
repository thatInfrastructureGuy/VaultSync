package main

import (
	"log"

	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/consumer"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/vault"
)

func main() {
	// Select the destination
	destination := consumer.Consumer{}
	err := destination.SelectConsumer()
	if err != nil {
		log.Fatal(err)
	}

	// Get lastUpdated date timestamp from consumer
	destinationlastUpdated, err := destination.GetLastUpdatedDate()
	if err != nil {
		log.Fatal(err)
	}

	// Select the source
	source := vault.Vault{}
	err = source.SelectProvider(destinationlastUpdated)
	if err != nil {
		log.Fatal(err)
	}
	// Poll secrets from vault which were updated since lastUpdated value
	secretList, err := source.GetSecrets()
	if err != nil {
		log.Fatal(err)
	}

	// Update kuberenetes secrets
	err = destination.PostSecrets(secretList)
	if err != nil {
		log.Fatal(err)
	}
}
