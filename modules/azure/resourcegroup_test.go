package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure Virtual Machines, these tests can be extended
(see AWS S3 tests for reference).
*/

func TestResourceGroupExists(t *testing.T) {
	t.Parallel()

	resourceGroupName := "fakeResourceGroupName"
	_, err := ResourceGroupExistsE(resourceGroupName, "")
	require.Error(t, err)
}
