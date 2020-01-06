package keyvault

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.0/pkg/common/data"
)

// ListSecrets Get all the secrets from specified keyvault
func ListSecrets(basicClient keyvault.BaseClient) map[string]data.SecretAttribute {
	currentTimeUTC := time.Now().UTC()
	ctx := context.Background()
	secretItr, err := basicClient.GetSecrets(ctx, "https://"+vaultName+".vault.azure.net", nil)
	if err != nil {
		fmt.Printf("unable to get list of secrets: %v\n", err)
		os.Exit(1)
	}

	secretList := make(map[string]data.SecretAttribute)

	for {
		if secretItr.Values() == nil {
			break
		}

		for _, secretProperties := range secretItr.Values() {
			var activationDate, expiryDate time.Time

			lastUpdated := time.Time(*secretProperties.Attributes.Updated)
			secretName := path.Base(*secretProperties.ID)

			// Check Activation date
			if secretProperties.Attributes.NotBefore != nil {
				activationDate = time.Time(*secretProperties.Attributes.NotBefore)
				if activationDate.After(currentTimeUTC) {
					fmt.Printf("%v key is not activated yet\n", secretName)
					continue
				}
			}

			// Check Expiry date
			if secretProperties.Attributes.Expires != nil {
				expiryDate = time.Time(*secretProperties.Attributes.Expires)
				if expiryDate.Before(currentTimeUTC) {
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
			secretList[secretName] = data.SecretAttribute{
				LastUpdated:    lastUpdated,
				ActivationDate: activationDate,
				ExpiryDate:     expiryDate,
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
