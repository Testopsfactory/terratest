// +build azure azureslim,network

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods can be mocked or Create/Delete APIs are added, these tests can be extended.
*/

func TestGetNetworkInterfaceE(t *testing.T) {
	t.Parallel()

	nicName := ""
	rgName := ""
	subID := ""

	_, err := GetNetworkInterfaceE(t, nicName, rgName, subID)

	require.Error(t, err)
}

func TestGetNetworkInterfacePrivateIPsE(t *testing.T) {
	t.Parallel()

	nicName := ""
	rgName := ""
	subID := ""

	_, err := GetNetworkInterfacePrivateIPsE(t, nicName, rgName, subID)

	require.Error(t, err)
}

func TestGetNetworkInterfacePublicIPsE(t *testing.T) {
	t.Parallel()

	nicName := ""
	rgName := ""
	subID := ""

	_, err := GetNetworkInterfacePublicIPsE(t, nicName, rgName, subID)

	require.Error(t, err)
}

func TestNetworkInterfaceExistsE(t *testing.T) {
	t.Parallel()

	nicName := ""
	rgName := ""
	subID := ""

	_, err := NetworkInterfaceExistsE(t, nicName, rgName, subID)

	require.Error(t, err)
}
