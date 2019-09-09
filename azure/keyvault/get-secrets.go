package keyvault

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
)

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
