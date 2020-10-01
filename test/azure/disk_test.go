// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/require"
)

func TestGetDiskE(t *testing.T) {
	t.Parallel()

	diskName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetDiskE(t, diskName, rgName, subID)

	require.Error(t, err)
}

func TestGetDiskTypeE(t *testing.T) {
	t.Parallel()

	diskName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetDiskTypeE(t, diskName, rgName, subID)

	require.Error(t, err)
}
