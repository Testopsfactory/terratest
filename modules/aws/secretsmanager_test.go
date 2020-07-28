package aws

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretsManagerMethods(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	name := random.UniqueId()
	description := "This is just a secrets manager test description."
	secretValue := "This is the secret value."

	secretARN, err := CreateSecretStringWithDefaultKey(t, region, description, name, secretValue)
	require.NoError(t, err)

	defer deleteSecret(t, region, secretARN)

	storedValue, err := GetSecretValue(t, region, secretARN)
	require.NoError(t, err)
	assert.Equal(t, secretValue, storedValue)
}

func deleteSecret(t *testing.T, region, id string) {
	err := DeleteSecret(t, region, id, true)
	require.NoError(t, err)

	_, err = GetSecretValue(t, region, id)
	require.Error(t, err)
}
