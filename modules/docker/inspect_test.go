package docker

import (
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

const image = "nginx:1.17-alpine"

func TestInspect(t *testing.T) {
	t.Parallel()

	// append timestamp to container name to allow running tests in parallel
	name := "inspect-test-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// running the container detached to allow inspection while it is running
	options := &RunOptions{
		Detach: true,
		Name:   name,
	}

	id := Run(t, image, options)
	defer removeContainer(t, id)

	c := Inspect(t, id)

	require.Equal(t, id, c.ID)
	require.Equal(t, name, c.Name)
	require.IsType(t, time.Time{}, c.Created)
	require.Equal(t, true, c.Running)
}

func removeContainer(t *testing.T, id string) {
	cmd := shell.Command{
		Command: "docker",
		Args:    []string{"container", "rm", "-f", id},
	}

	shell.RunCommand(t, cmd)
}