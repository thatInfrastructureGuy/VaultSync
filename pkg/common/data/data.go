package data

import "time"

// SecretAttribute is constructed after querying Vault for each secret.
// It contains various attributes of secret other than values.
type SecretAttribute struct {
	LastUpdated    time.Time
	ActivationDate time.Time
	ExpiryDate     time.Time
	Value          string
	IsEnabled      bool
}
