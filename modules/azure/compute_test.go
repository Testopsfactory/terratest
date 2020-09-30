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
If/when CRUD methods are introduced for Azure Virtual Machines, these tests can be extended
(see AWS S3 tests for reference).
*/

func TestGetVirtualMachineE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineInstanceE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineInstanceE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetResourceGroupVirtualMachinesObjectsE(t *testing.T) {
	t.Parallel()

	rgName := ""
	subID := ""

	_, err := GetResourceGroupVirtualMachinesObjectsE(t, rgName, subID)

	require.Error(t, err)
}

func TestGetResourceGroupVirtualMachinesE(t *testing.T) {
	t.Parallel()

	rgName := ""
	subID := ""

	_, err := GetResourceGroupVirtualMachinesE(t, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineTagsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineTagsE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineSizeE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineSizeE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineAdminUserE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineAdminUserE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineImageE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineImageE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineAvailabilitySetIDE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineAvailabilitySetIDE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineStateE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineStateE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineOsDiskNameE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineOsDiskNameE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineManagedDiskCountE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineManagedDiskCountE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineManagedDisksE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineManagedDisksE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineNicCountE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineNicCountE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineNicsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineNicsE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestVirtualMachineExistsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := VirtualMachineExistsE(t, vmName, rgName, subID)

	require.Error(t, err)
}
