// +build azure azureslim,network

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
If/when methods can be mocked or Create/Delete APIs are added, these tests can be extended.
*/

func TestGetVirtualNetworksClientE(t *testing.T) {
	t.Parallel()

	subID := ""

	_, err := azure.GetVirtualNetworksClientE(subID)

	require.Error(t, err)
}

func TestGetVirtualNetworkE(t *testing.T) {
	t.Parallel()

	vnetName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualNetworkE(t, vnetName, rgName, subID)

	require.Error(t, err)
}

func TestGetSubnetClientE(t *testing.T) {
	t.Parallel()

	subID := ""

	_, err := azure.GetSubnetClientE(subID)

	require.Error(t, err)
}

func TestGetSubnetE(t *testing.T) {
	t.Parallel()

	subnetName := ""
	vnetName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetSubnetE(t, subnetName, vnetName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualNetworkDNSServerIPsE(t *testing.T) {
	t.Parallel()

	vnetName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualNetworkDNSServerIPsE(t, vnetName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualNetworkSubnetsE(t *testing.T) {
	t.Parallel()

	vnetName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetVirtualNetworkSubnetsE(t, vnetName, rgName, subID)

	require.Error(t, err)
}

func TestGetSubnetIPRangeE(t *testing.T) {
	t.Parallel()

	subnetName := ""
	vnetName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetSubnetIPRangeE(t, subnetName, vnetName, rgName, subID)

	require.Error(t, err)
}

func TestCheckSubnetContainsIPE(t *testing.T) {
	t.Parallel()

	ipAddress := ""
	subnetName := ""
	vnetName := ""
	rgName := ""
	subID := ""

	_, err := azure.CheckSubnetContainsIPE(t, ipAddress, subnetName, vnetName, rgName, subID)

	require.Error(t, err)
}

func TestSubnetExistsE(t *testing.T) {
	t.Parallel()

	subnetName := ""
	vnetName := ""
	rgName := ""
	subID := ""

	_, err := azure.SubnetExistsE(t, subnetName, vnetName, rgName, subID)

	require.Error(t, err)
}

func TestVirtualNetworkExistsE(t *testing.T) {
	t.Parallel()

	vnetName := ""
	rgName := ""
	subID := ""

	_, err := azure.VirtualNetworkExistsE(t, vnetName, rgName, subID)

	require.Error(t, err)
}
