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

func TestGetTagsForVirtualMachineE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetTagsForVirtualMachineE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetSizeOfVirtualMachineE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetSizeOfVirtualMachineE(t, vmName, rgName, subID)

	require.Error(t, err)
}
