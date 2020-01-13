package keyvault

import (
	"fmt"
	"os"
	"strings"
)

var vaultName, setDebug, convertHyphenToUnderscores string

func init() {
	if os.Getenv("AZURE_TENANT_ID") == "" || os.Getenv("AZURE_CLIENT_ID") == "" || os.Getenv("AZURE_CLIENT_SECRET") == "" || os.Getenv("VAULT_NAME") == "" {
		fmt.Println("env vars not set, exiting...")
		os.Exit(1)
	}

	vaultName = os.Getenv("VAULT_NAME")
	convertHyphenToUnderscores = strings.ToLower(os.Getenv("CONVERT_HYPHENS_TO_UNDERSCORES"))
}
