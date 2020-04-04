package docker

import (
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInspect(t *testing.T) {
	t.Parallel()

	cName := "foobar" // @TODO: Generate random string to avoid name conflicts
	cImage := "alpine:3.6"

	options := &RunOptions{
		Detach: true,
		Name: cName,
		OtherOptions:         []string{"--expose=80", "--publish-all"},
	}

	id := Run(t, cImage, options)
	defer removeContainer(t, id)

	container := Inspect(t, id)

	require.Equal(t, "foobar", container.Name)
}

func removeContainer(t *testing.T, id string) {
	cmd := shell.Command{
		Command: "docker",
		Args:    []string{"container", "rm", "-f", id},
	}

	shell.RunCommand(t, cmd)
}