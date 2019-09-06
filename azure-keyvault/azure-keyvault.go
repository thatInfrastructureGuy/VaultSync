package main

/*
 * You need to set four environment variables before using the app:
 * AZURE_TENANT_ID: Your Azure tenant ID
 * AZURE_CLIENT_ID: Your Azure client ID. This will be an app ID from your AAD.
 * AZURE_CLIENT_SECRET: The secret for the client ID above.
 * KVAULT: The name of your vault (just the name, not the full URL/path)
 *
 * Usage
 * List the secrets currently in the vault (not the values though):
 * kv-pass
 *
 * Get the value for a secret in the vault:
 * kv-pass YOUR_SECRETS_NAME
 *
 * Add or Update a secret in the vault:
 * kv-pass -edit YOUR_NEW_VALUE YOUR_SECRETS_NAME
 *
 * Delete a secret in the vault:
 * kv-pass -delete YOUR_SECRETS_NAME
 */

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest"
)

var (
	setDebug  = flag.Bool("debug", false, "debug")
	vaultName string
)

func main() {
	flag.Parse()

	//if os.Getenv("AZURE_TENANT_ID") == "" || os.Getenv("AZURE_CLIENT_ID") == "" || os.Getenv("AZURE_CLIENT_SECRET") == "" || os.Getenv("KVAULT") == "" {
	//	fmt.Println("env vars not set, exiting...")
	//	os.Exit(1)
	//}

	os.Setenv("AZURE_TENANT_ID", "ad4b7142-cac0-4b8b-9573-c02311c68a26")
	os.Setenv("AZURE_CLIENT_ID", "d70a6390-35a9-443e-9e8e-805cd93aa68b")
	os.Setenv("AZURE_CLIENT_SECRET", "a40c02e7-2dab-47b3-a470-65d6451d6f13")
	os.Setenv("KVAULT", "kvtestin")

	vaultName = os.Getenv("KVAULT")

	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		fmt.Printf("unable to create vault authorizer: %v\n", err)
		os.Exit(1)
	}

	basicClient := keyvault.New()
	basicClient.Authorizer = authorizer

	if *setDebug {
		basicClient.RequestInspector = logRequest()
		basicClient.ResponseInspector = logResponse()
	}

	listSecrets(basicClient)
}

// Get all the secrets from specified keyvault
func listSecrets(basicClient keyvault.BaseClient) {
	ctx := context.Background()
	secretItr, err := basicClient.GetSecrets(ctx, "https://"+vaultName+".vault.azure.net", nil)
	if err != nil {
		fmt.Printf("unable to get list of secrets: %v\n", err)
		os.Exit(1)
	}

	type secretAttribute struct {
		dateUpdated float64
		Activates   float64
		Expires     float64
		Enabled     bool
		value       string
	}

	secretList := make(map[string]secretAttribute)

	for {
		if secretItr.Values() == nil {
			break
		}

		for _, secretValue := range secretItr.Values() {
			secretList[*secretValue.ID] = secretAttribute{
				dateUpdated: secretValue.Attributes.Updated.Duration().Seconds(),
				//Activates:   secretValue.Attributes.NotBefore.Duration().Seconds(),
				//Expires:     secretValue.Attributes.Expires.Duration().Seconds(),
				Enabled: *secretValue.Attributes.Enabled,
				value:   getSecret(basicClient, path.Base(*secretValue.ID), *secretValue.Attributes.Enabled),
			}
		}

		err = secretItr.NextWithContext(ctx)
		if err != nil {
			fmt.Printf("unable to get next page for list of secrets: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println(secretList)
}

func getSecret(basicClient keyvault.BaseClient, secname string, isEnabled bool) (value string) {
	if !isEnabled {
		return ""
	}
	secretResp, err := basicClient.GetSecret(context.Background(), "https://"+vaultName+".vault.azure.net", secname, "")
	if err != nil {
		fmt.Printf("unable to get value for secret: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(*secretResp.Value)
	return *secretResp.Value
}

func logRequest() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				log.Println(err)
			}
			dump, _ := httputil.DumpRequestOut(r, true)
			log.Println(string(dump))
			return r, err
		})
	}
}

func logResponse() autorest.RespondDecorator {
	return func(p autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(r *http.Response) error {
			err := p.Respond(r)
			if err != nil {
				log.Println(err)
			}
			dump, _ := httputil.DumpResponse(r, true)
			log.Println(string(dump))
			return err
		})
	}
}
