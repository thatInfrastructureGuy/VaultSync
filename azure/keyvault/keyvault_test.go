package keyvault

import (
	"testing"
)

func TestListSecrets(t *testing.T) {

	secretAttribute := ListSecrets(Initializer())
	t.Log(secretAttribute)
}
