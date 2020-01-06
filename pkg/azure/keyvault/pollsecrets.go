package keyvault

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
)

// ListSecrets Get all the secrets from specified keyvault
func ListSecrets(basicClient keyvault.BaseClient) map[string]SecretAttribute {
	currentTimeUTC := time.Now().UTC().Unix()
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
			var activates, expires int64
			dateUpdated := int64(secretProperties.Attributes.Updated.Duration().Seconds())
			secretName := path.Base(*secretProperties.ID)

			// Check if activation date is passed.
			if secretProperties.Attributes.NotBefore != nil {
				activates = int64(secretProperties.Attributes.NotBefore.Duration().Seconds())
				if activates > currentTimeUTC {
					fmt.Printf("%v key is not activated yet\n", secretName)
					continue
				}
			}

			// Check Expiry date
			if secretProperties.Attributes.Expires != nil {
				expires = int64(secretProperties.Attributes.Expires.Duration().Seconds())
				if expires < currentTimeUTC {
					fmt.Printf("%v key has expired\n", secretName)
					continue
				}
			}

			// Check if secret is disabled.
			isEnabled := *secretProperties.Attributes.Enabled
			if !isEnabled {
				continue
			}

			secretValue := getSecret(basicClient, secretName)

			// Check if ALL hyphers should be converted to underscores
			if convertHyphenToUnderscores == "true" {
				secretName = strings.ReplaceAll(secretName, "-", "_")
			}

			//Create Key-Value map
			secretList[secretName] = SecretAttribute{
				DateUpdated:    dateUpdated,
				ActivationDate: activates,
				ExpiryDate:     expires,
				Value:          secretValue,
				IsEnabled:      isEnabled,
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
func getSecret(basicClient keyvault.BaseClient, secretName string) (value string) {
	secretResp, err := basicClient.GetSecret(context.Background(), "https://"+vaultName+".vault.azure.net", secretName, "")
	if err != nil {
		fmt.Printf("unable to get value for secret: %v\n", err)
		os.Exit(1)
	}

	return *secretResp.Value
}
