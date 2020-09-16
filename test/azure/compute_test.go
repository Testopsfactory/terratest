// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
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

	_, err := azure.GetVirtualMachineE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineInstanceE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineInstanceE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetResourceGroupVirtualMachinesObjectsE(t *testing.T) {
	t.Parallel()

	rgName := ""
	subID := ""

	_, err := azure.GetResourceGroupVirtualMachinesObjectsE(t, rgName, subID)

	require.Error(t, err)
}

func TestGetResourceGroupVirtualMachinesE(t *testing.T) {
	t.Parallel()

	rgName := ""
	subID := ""

	_, err := azure.GetResourceGroupVirtualMachinesE(t, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineTagsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineTagsE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineSizeE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineSizeE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineAdminUserE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineAdminUserE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineImageE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineImageE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineAvailabilitySetIDE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineAvailabilitySetIDE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineStateE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineStateE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineOsDiskNameE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineOsDiskNameE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineManagedDiskCountE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineManagedDiskCountE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineManagedDisksE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineManagedDisksE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineNicCountE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineNicCountE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineNicsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualMachineNicsE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestVirtualMachineExistsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.VirtualMachineExistsE(t, vmName, rgName, subID)

	require.Error(t, err)
}
