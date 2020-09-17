package azure

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/require"
)

func TestResourceGroupExists(t *testing.T) {
	t.Parallel()

	resourceGroupName := "fakeResourceGroupName"
	_, err := azure.ResourceGroupExistsE(resourceGroupName, "")
	require.Error(t, err)
}
