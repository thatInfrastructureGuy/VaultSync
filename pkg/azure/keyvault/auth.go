package keyvault

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
)

type Keyvault struct {
	basicClient keyvault.BaseClient
}

// Initializer creates KeyVault instance
func (k *Keyvault) Initializer() {
	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Printf("unable to create vault authorizer: %v\n", err)
		os.Exit(1)
	}

	k.basicClient = keyvault.New()
	k.basicClient.Authorizer = authorizer

	if strings.ToLower(setDebug) == "true" {
		k.basicClient.RequestInspector = logRequest()
		k.basicClient.ResponseInspector = logResponse()
	}

}
