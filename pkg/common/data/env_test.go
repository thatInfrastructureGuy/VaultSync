package data

import (
	"os"
	"testing"
)

func ExecuteGetenv(envKeyList []string, envValueList []string) error {
	if envKeyList != nil && envValueList != nil && len(envKeyList) == len(envValueList) {
		for index, key := range envKeyList {
			value := envValueList[index]
			os.Setenv(key, value)
		}
	}

	var envVars Env
	err := envVars.Getenv()

	if envKeyList != nil && envValueList != nil && len(envKeyList) == len(envValueList) {
		for _, key := range envKeyList {
			os.Unsetenv(key)
		}
	}

	return err
}

func TestGetenvVaultName(t *testing.T) {
	keys := []string{"PROVIDER"}
	vals := []string{"something"}
	got := ExecuteGetenv(keys, vals)
	want := "VAULT_NAME env not present"

	if got.Error() != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetenvProvider(t *testing.T) {
	got := ExecuteGetenv(nil, nil)
	want := "PROVIDER env not present"

	if got.Error() != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetenvRequired(t *testing.T) {
	keys := []string{"PROVIDER", "VAULT_NAME"}
	vals := []string{"something", "myVault"}
	err := ExecuteGetenv(keys, vals)
	if err != nil {
		t.Errorf("Errors after inputing Required keys: %v", err)
	}
	gotConsumer := os.Getenv("CONSUMER")
	gotSecretName := os.Getenv("SECRET_NAME")
	gotSecretNamespace := os.Getenv("SECRET_NAMESPACE")
	gotRefreshRate := os.Getenv("REFRESH_RATE")
	gotConvertHyphensToUnderscores := os.Getenv("CONVERT_HYPHENS_TO_UNDERSCORES")
	want := "PROVIDER env not present"

	if got.Error() != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func checkEnvVars(envKeyList, envValueList) bool {
	if envKeyList != nil && envValueList != nil && len(envKeyList) == len(envValueList) {
		for index, key := range envKeyList {
			value == envValueList[index]
			if key != value {
				return error.Error()
			}
		}
	}
	return nil
}
