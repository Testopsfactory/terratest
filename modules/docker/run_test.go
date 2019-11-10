package docker

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()

	options := &RunOptions{
		Command:              []string{"-c", `echo "Hello, $NAME!"`},
		Entrypoint:           "sh",
		EnvironmentVariables: []string{"NAME=World"},
		Remove:               true,
	}

	out := Run(t, "alpine:3.7", options)
	require.Equal(t, "Hello, World!", out)
}
