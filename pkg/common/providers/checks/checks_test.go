package checks

import (
	"testing"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
)

var UpdatedDate time.Time = time.Now()

var testDataSet = []struct {
	title              string
	env                *data.Env
	sourceDate         time.Time
	destinationDate    time.Time
	originalSecretName string
	updatedSecretName  string
	skipUpdate         bool
}{
	{
		title:              "Do Nothing",
		env:                &data.Env{},
		sourceDate:         UpdatedDate,
		destinationDate:    UpdatedDate.Add(1 * time.Second),
		originalSecretName: "test-1",
		updatedSecretName:  "test-1",
		skipUpdate:         true,
	},
	{
		title:              "Do Nothing. Same date.",
		env:                &data.Env{},
		sourceDate:         UpdatedDate,
		destinationDate:    UpdatedDate,
		originalSecretName: "test-1",
		updatedSecretName:  "test-1",
		skipUpdate:         true,
	},
	{
		title: "Convert Hyphen to Underscore",
		env: &data.Env{
			ConvertHyphensToUnderscores: true,
		},
		sourceDate:         UpdatedDate,
		destinationDate:    UpdatedDate,
		originalSecretName: "Test-2",
		updatedSecretName:  "Test_2",
		skipUpdate:         true,
	},
	{
		title:              "Update Secret.",
		env:                &data.Env{},
		sourceDate:         UpdatedDate.Add(1 * time.Second),
		destinationDate:    UpdatedDate,
		originalSecretName: "Test_3",
		updatedSecretName:  "Test_3",
		skipUpdate:         false,
	},
	{
		title: "Update Secret. Convert Hyphens.",
		env: &data.Env{
			ConvertHyphensToUnderscores: true,
		},
		sourceDate:         UpdatedDate.Add(1 * time.Second),
		destinationDate:    UpdatedDate,
		originalSecretName: "Test-4",
		updatedSecretName:  "Test_4",
		skipUpdate:         false,
	},
}

func TestCommonProviderChecks(t *testing.T) {
	for _, testData := range testDataSet {
		testData := testData
		t.Run(testData.title, func(t *testing.T) {
			t.Parallel()
			gotUpdatedSecretName, gotSkipUpdate := CommonProviderChecks(testData.env, testData.originalSecretName, testData.sourceDate, testData.destinationDate)
			compareData(t, gotUpdatedSecretName, gotSkipUpdate, testData.updatedSecretName, testData.skipUpdate)
		})
	}
}

func compareData(t *testing.T, gotUpdatedSecretName string, gotSkipUpdate bool, wantUpdatedSecretName string, wantSkipUpdate bool) {
	if gotUpdatedSecretName != wantUpdatedSecretName {
		t.Errorf("SecretName: got %v want %v", gotUpdatedSecretName, wantUpdatedSecretName)
	}
	if gotSkipUpdate != wantSkipUpdate {
		t.Errorf("SkipUpdate: got %v want %v", gotSkipUpdate, wantSkipUpdate)
	}
}
