package checks

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func CommonProviderChecks(originalSecretName string, sourceDate time.Time, destinationDate time.Time) (updatedSecretName string, skipUpdate bool) {
	// Set updatedName as original name
	updatedSecretName = originalSecretName
	// Check if destination keys are outdated.
	if !sourceDate.After(destinationDate) {
		fmt.Printf("%v key is not updated since %v . Skipping update.\n", originalSecretName, sourceDate)
		skipUpdate = true
	}
	// Check if ALL hyphers should be converted to underscores
	convertHyphenToUnderscores := strings.ToLower(os.Getenv("CONVERT_HYPHENS_TO_UNDERSCORES"))
	if convertHyphenToUnderscores == "true" {
		updatedSecretName = strings.ReplaceAll(originalSecretName, "-", "_")
	}
	return updatedSecretName, skipUpdate
}
