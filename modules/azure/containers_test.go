// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure MySQL server and database, these tests can be extended
*/

func TestContainerRegistryExistsE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	registryName := ""
	subscriptionID := ""

	_, err := ContainerRegistryExistsE(registryName, resGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetContainerRegistryE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	registryName := ""
	subscriptionID := ""

	_, err := GetContainerRegistryE(registryName, resGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetContainerRegistryClientE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""

	_, err := GetContainerRegistryClientE(subscriptionID)
	require.NoError(t, err)
}
