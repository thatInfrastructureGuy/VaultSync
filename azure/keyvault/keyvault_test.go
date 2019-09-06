package keyvault

import (
	"fmt"
	"os"
	"testing"
)

func TestListSecrets(t *testing.T) {
	os.Setenv("AZURE_TENANT_ID", "ad4b7142-cac0-4b8b-9573-c02311c68a26")
	os.Setenv("AZURE_CLIENT_ID", "d70a6390-35a9-443e-9e8e-805cd93aa68b")
	os.Setenv("AZURE_CLIENT_SECRET", "a40c02e7-2dab-47b3-a470-65d6451d6f13")
	os.Setenv("KVAULT", "kvtestin")
	os.Setenv("CONVERT_HYPHENS_TO_UNDERSCORES", "TRUE")

	secretAttribute := ListSecrets(basicClient)
	fmt.Println(secretAttribute)

}
