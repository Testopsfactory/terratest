package docker

import (
	"encoding/json"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type InspectOutput struct {
	Id string
	Created string
	Name string
}


type ContainerInspect struct {
	ID string
	Name string
}

func Inspect(t *testing.T, id string) ContainerInspect {
	// @TODO: Validate if id is a valid containerID

	cmd := shell.Command{
		Command: "docker",
		Args:    []string{"container", "inspect", id},
	}

	out, err := shell.RunCommandAndGetOutputE(t, cmd)
	require.NoError(t, err)

	var containers []InspectOutput
	err = json.Unmarshal([]byte(out), &containers)
	require.NoError(t, err)

	if len(containers) == 0 {
		return ContainerInspect{}
	}

	container := containers[0]

	inspect := ContainerInspect{
		ID: container.Id,
		Name: strings.TrimLeft(container.Name, "/"),
	}

	return inspect
}
