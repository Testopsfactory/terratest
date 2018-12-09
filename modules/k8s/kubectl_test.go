package k8s

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test that RunKubectlAndGetOutputE will run kubectl and return the output by running a can-i command call.
func TestRunKubectlAndGetOutputReturnsOutput(t *testing.T) {
	options := NewKubectlOptions("", "")
	output, err := RunKubectlAndGetOutputE(t, options, "auth", "can-i", "get", "pods")
	require.NoError(t, err)
	require.Equal(t, output, "yes")
}
