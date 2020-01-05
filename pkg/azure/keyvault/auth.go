package keyvault

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
)

// Initializer creates KeyVault instance
func Initializer() keyvault.BaseClient {
	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Printf("unable to create vault authorizer: %v\n", err)
		os.Exit(1)
	}

	basicClient := keyvault.New()
	basicClient.Authorizer = authorizer

	if strings.ToLower(setDebug) == "true" {
		basicClient.RequestInspector = logRequest()
		basicClient.ResponseInspector = logResponse()
	}

	return basicClient
}
