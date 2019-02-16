package helm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Delete(t *testing.T, options *Options, releaseName string, purge bool) {
	require.NoError(t, DeleteE(t, options, releaseName, purge))
}

func DeleteE(t *testing.T, options *Options, releaseName string, purge bool) error {
	args := []string{}
	if purge {
		args = append(args, "--purge")
	}
	args = append(args, releaseName)
	_, err := RunHelmCommandAndGetOutputE(t, options, "delete", args...)
	return err
}
