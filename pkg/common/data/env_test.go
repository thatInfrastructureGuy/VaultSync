package data

import (
	"os"
	"reflect"
	"testing"
)

var testDataSets = []struct {
	title string
	in    map[string]string
	out   Env
	err   string
}{
	{
		title: "Set No Environment Variables",
		in:    nil,
		out:   Env{},
		err:   "PROVIDER env not present",
	},
	{
		title: "Vault Name environment variable not set",
		in: map[string]string{
			"PROVIDER": "provider",
		},
		out: Env{
			Provider: "provider",
		},
		err: "VAULT_NAME env not present",
	},
	{
		title: "Set Only Required Environment Variables",
		in: map[string]string{
			"PROVIDER":   "provider",
			"VAULT_NAME": "myVault",
		},
		out: Env{
			Provider:                    "provider",
			VaultName:                   "myVault",
			ConsumerType:                "kubernetes",
			DeploymentList:              nil,
			StatefulsetList:             nil,
			SecretName:                  "myVault",
			Namespace:                   "default",
			RefreshRate:                 60,
			ConvertHyphensToUnderscores: false,
		},
		err: "",
	},
	{
		title: "Invalid Refresh Rate",
		in: map[string]string{
			"PROVIDER":     "provider",
			"VAULT_NAME":   "myVault",
			"REFRESH_RATE": "48.95",
		},
		out: Env{
			Provider:                    "provider",
			VaultName:                   "myVault",
			ConsumerType:                "kubernetes",
			DeploymentList:              nil,
			StatefulsetList:             nil,
			SecretName:                  "myVault",
			Namespace:                   "default",
			RefreshRate:                 60,
			ConvertHyphensToUnderscores: false,
		},
		err: "",
	},
	{
		title: "Set All Environment Variables",
		in: map[string]string{
			"PROVIDER":                       "provider",
			"CONSUMER":                       "kubernetes",
			"VAULT_NAME":                     "myVault",
			"SECRET_NAME":                    "mySecret",
			"SECRET_NAMESPACE":               "ns",
			"REFRESH_RATE":                   "75",
			"CONVERT_HYPHENS_TO_UNDERSCORES": "true",
			"DEPLOYMENT_LIST":                "deploy1,deploy2,deploy3",
			"STATEFULSET_LIST":               "sts1,sts2,sts3,sts4",
		},
		out: Env{
			Provider:                    "provider",
			VaultName:                   "myVault",
			ConsumerType:                "kubernetes",
			DeploymentList:              []string{"deploy1", "deploy2", "deploy3"},
			StatefulsetList:             []string{"sts1", "sts2", "sts3", "sts4"},
			SecretName:                  "mySecret",
			Namespace:                   "ns",
			RefreshRate:                 75,
			ConvertHyphensToUnderscores: true,
		},
		err: "",
	},
}

func TestGetenv(t *testing.T) {
	for _, testData := range testDataSets {
		t.Run(testData.title, func(t *testing.T) {
			if testData.in != nil {
				for key, value := range testData.in {
					os.Setenv(key, value)
				}
			}
			var envVars Env
			err := envVars.Getenv()
			if testData.in != nil {
				for key, _ := range testData.in {
					os.Unsetenv(key)
				}
			}
			compareData(t, envVars, testData.out)
			if err != nil {
				if err.Error() != testData.err {
					t.Errorf("Error Got '%v' Error Want '%v'", err.Error(), testData.err)
				}
			}
		})
	}
}

func compareData(t *testing.T, got Env, want Env) {
	if got.Provider != want.Provider {
		t.Errorf("Got '%v' Want '%v'", got.Provider, want.Provider)
	}
	if got.VaultName != want.VaultName {
		t.Errorf("Got '%v' Want '%v'", got.VaultName, want.VaultName)
	}
	if got.ConsumerType != want.ConsumerType {
		t.Errorf("Got '%v' Want '%v'", got.ConsumerType, want.ConsumerType)
	}
	if got.SecretName != want.SecretName {
		t.Errorf("Got '%v' Want '%v'", got.SecretName, want.SecretName)
	}
	if got.Namespace != want.Namespace {
		t.Errorf("Got '%v' Want '%v'", got.Namespace, want.Namespace)
	}
	if got.RefreshRate != want.RefreshRate {
		t.Errorf("Got '%v' Want '%v'", got.RefreshRate, want.RefreshRate)
	}
	if got.ConvertHyphensToUnderscores != want.ConvertHyphensToUnderscores {
		t.Errorf("Got '%v' Want '%v'", got.ConvertHyphensToUnderscores, want.ConvertHyphensToUnderscores)
	}
	if !reflect.DeepEqual(got.DeploymentList, want.DeploymentList) {
		t.Errorf("Got '%v' Want '%v'", got.DeploymentList, want.DeploymentList)
	}
	if !reflect.DeepEqual(got.StatefulsetList, want.StatefulsetList) {
		t.Errorf("Got '%v' Want '%v'", got.StatefulsetList, want.StatefulsetList)
	}
}
