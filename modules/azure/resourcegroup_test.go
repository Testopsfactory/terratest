package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceGroupExists(t *testing.T) {
	t.Parallel()

	resourceGroupName := "fakeResourceGroupName"
	_, err := ResourceGroupExistsE(resourceGroupName, "")
	require.Error(t, err)
}
