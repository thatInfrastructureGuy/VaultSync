package keyvault

import (
	"fmt"
	"os"

	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
)

// Initializer creates KeyVault instance
func (k *Keyvault) Initializer() {
	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Printf("unable to create vault authorizer: %v\n", err)
		os.Exit(1)
	}

	k.basicClient = keyvault.New()
	k.basicClient.Authorizer = authorizer
}
