package keyvault

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var vaultName, setDebug, convertHyphenToUnderscores string

// SecretAttribute is constructed after querying Vault for each secret.
// It contains various attributes of secret other than values.
type SecretAttribute struct {
	LastUpdated    time.Time
	ActivationDate time.Time
	ExpiryDate     time.Time
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
