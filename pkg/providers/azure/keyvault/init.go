package keyvault

import (
	"fmt"
	"os"
)

func init() {
	if os.Getenv("AZURE_TENANT_ID") == "" || os.Getenv("AZURE_CLIENT_ID") == "" || os.Getenv("AZURE_CLIENT_SECRET") == "" || os.Getenv("VAULT_NAME") == "" {
		fmt.Println("azure credentials env vars not set, exiting...")
		os.Exit(1)
	}
}
