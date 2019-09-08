package keyvault

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest"
)

var vaultName, setDebug, convertHyphenToUnderscores string

// SecretAttribute is constructed after querying Vault for each secret.
// It contains various attributes of secret other than values.
type SecretAttribute struct {
	DateUpdated    float64
	ActivationDate float64
	ExpiryDate     float64
	Value          string
	IsEnabled      bool
}

func init() {
	if os.Getenv("AZURE_TENANT_ID") == "" || os.Getenv("AZURE_CLIENT_ID") == "" || os.Getenv("AZURE_CLIENT_SECRET") == "" || os.Getenv("KVAULT") == "" {
		fmt.Println("env vars not set, exiting...")
		os.Exit(1)
	}

	vaultName = os.Getenv("KVAULT")
	setDebug = strings.ToLower(os.Getenv("DEBUG"))
	convertHyphenToUnderscores = strings.ToLower(os.Getenv("CONVERT_HYPHENS_TO_UNDERSCORES"))
}

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

// ListSecrets Get all the secrets from specified keyvault
func ListSecrets(basicClient keyvault.BaseClient) map[string]SecretAttribute {
	ctx := context.Background()
	secretItr, err := basicClient.GetSecrets(ctx, "https://"+vaultName+".vault.azure.net", nil)
	if err != nil {
		fmt.Printf("unable to get list of secrets: %v\n", err)
		os.Exit(1)
	}

	secretList := make(map[string]SecretAttribute)

	for {
		if secretItr.Values() == nil {
			break
		}

		for _, secretProperties := range secretItr.Values() {
			var activates, expires float64
			if secretProperties.Attributes.NotBefore != nil {
				activates = secretProperties.Attributes.NotBefore.Duration().Seconds()
			}
			if secretProperties.Attributes.Expires != nil {
				expires = secretProperties.Attributes.Expires.Duration().Seconds()
			}
			secretName := path.Base(*secretProperties.ID)
			secretValue := getSecret(basicClient, secretName, *secretProperties.Attributes.Enabled)

			// Check if ALL hyphers should be converted to underscores
			if convertHyphenToUnderscores == "true" {
				secretName = strings.ReplaceAll(secretName, "-", "_")
			}

			//Create Key-Value map
			secretList[secretName] = SecretAttribute{
				DateUpdated:    secretProperties.Attributes.Updated.Duration().Seconds(),
				ActivationDate: activates,
				ExpiryDate:     expires,
				Value:          secretValue,
				IsEnabled:      *secretProperties.Attributes.Enabled,
			}
		}

		err = secretItr.NextWithContext(ctx)
		if err != nil {
			fmt.Printf("unable to get next page for list of secrets: %v\n", err)
			os.Exit(1)
		}
	}

	return secretList
}

// Get SecretValue from KeyVault if Secret is enabled.
// If secret is disabled, return empty string.
func getSecret(basicClient keyvault.BaseClient, secretName string, isEnabled bool) (value string) {
	if !isEnabled {
		return ""
	}
	secretResp, err := basicClient.GetSecret(context.Background(), "https://"+vaultName+".vault.azure.net", secretName, "")
	if err != nil {
		fmt.Printf("unable to get value for secret: %v\n", err)
		os.Exit(1)
	}

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
