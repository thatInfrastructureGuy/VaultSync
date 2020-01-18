package keyvault

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/thatInfrastructureGuy/VaultSync/v0.0.1/pkg/common/data"
)

// listSecrets Get all the secrets from specified keyvault
func (k *Keyvault) listSecrets() (secretList map[string]data.SecretAttribute, err error) {
	ctx := context.Background()
	secretItr, err := k.basicClient.GetSecrets(ctx, "https://"+vaultName+".vault.azure.net", nil)
	if err != nil {
		fmt.Printf("unable to get list of secrets: %v\n", err)
		os.Exit(1)
	}

	secretList = make(map[string]data.SecretAttribute)

	for {
		if secretItr.Values() == nil {
			break
		}

		for _, secretProperties := range secretItr.Values() {
			secretName := path.Base(*secretProperties.ID)
			dateUpdated := time.Time(*secretProperties.Attributes.Updated)

			//Checks against key metadata
			secretName, skipUpdate := CommonProviderChecks(secretName, dateUpdated, k.DestinationLastUpdated)
			if skipUpdate {
				continue
			}
			skipUpdate = customProviderChecks(secretProperties)
			if skipUpdate {
				continue
			}

			//Get Secret Values
			secretValue, err := k.getSecretValue(secretName)
			if err != nil {
				return nil, err
			}

			//Create Key-Value map
			secretList[secretName] = data.SecretAttribute{
				DateUpdated: dateUpdated,
				Value:       secretValue,
			}
		}

		err = secretItr.NextWithContext(ctx)
		if err != nil {
			fmt.Println("unable to get next page for list of secrets.")
			return nil, err
		}
	}

	return secretList, nil
}

func customProviderChecks(secretProperties keyvault.SecretItem) (skipUpdate bool) {
	secretName := path.Base(*secretProperties.ID)
	currentTimeUTC := time.Now().UTC()
	// Check Activation date
	if secretProperties.Attributes.NotBefore != nil {
		activationDate := time.Time(*secretProperties.Attributes.NotBefore)
		if activationDate.After(currentTimeUTC) {
			fmt.Printf("%v key is not activated yet\n", secretName)
			skipUpdate = true
		}
	}

	// Check Expiry date
	if secretProperties.Attributes.Expires != nil {
		expiryDate := time.Time(*secretProperties.Attributes.Expires)
		if expiryDate.Before(currentTimeUTC) {
			fmt.Printf("%v key has expired\n", secretName)
			skipUpdate = true
		}
	}

	// Check if secret is disabled.
	isEnabled := *secretProperties.Attributes.Enabled
	if !isEnabled {
		skipUpdate = true
	}
	return skipUpdate
}

func CommonProviderChecks(originalSecretName string, sourceDate time.Time, destinationDate time.Time) (updatedSecretName string, skipUpdate bool) {
	// Set updatedName as original name
	updatedSecretName = originalSecretName
	// Check if destination keys are outdated.
	if !sourceDate.After(destinationDate) {
		fmt.Printf("%v key is not updated since %v . Skipping update.", originalSecretName, sourceDate)
		skipUpdate = true
	}
	// Check if ALL hyphers should be converted to underscores
	if convertHyphenToUnderscores == "true" {
		updatedSecretName = strings.ReplaceAll(originalSecretName, "-", "_")
	}
	return updatedSecretName, skipUpdate
}

// Get SecretValue from KeyVault if Secret is enabled.
// If secret is disabled, return empty string.
func (k *Keyvault) getSecretValue(secretName string) (value string, err error) {
	secretResp, err := k.basicClient.GetSecret(context.Background(), "https://"+vaultName+".vault.azure.net", secretName, "")
	if err != nil {
		fmt.Println("unable to get value for secret")
		return "", err
	}

	return *secretResp.Value, nil
}

func (k *Keyvault) GetSecrets() (secretList map[string]data.SecretAttribute, err error) {
	err = k.initialize()
	if err != nil {
		return nil, err
	}
	secretList, err = k.listSecrets()
	if err != nil {
		return nil, err
	}
	return secretList, nil
}
