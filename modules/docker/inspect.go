package docker

import (
	"encoding/json"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
	"time"
)

type InspectOutput struct {
	Id      string
	Created string
	Name    string
	State   struct {
		Status   string
		Running  bool
		ExitCode uint8
		Error    string
	}
	NetworkSettings struct {
		Ports map[string][]struct {
			HostIp   string
			HostPort string
		}
	}
}

type ContainerInspect struct {
	ID       string
	Name     string
	Created  time.Time
	Status   string
	Running  bool
	ExitCode uint8
	Error    string
	Ports    []Port
}

type Port struct {
	HostPort      uint16
	ContainerPort uint16
	Protocol      string
}

// Inspect runs the 'docker inspect {container id} command and returns a ContainerInspect
// struct, converted from the output JSON
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

	c := containers[0]

	return transformContainer(t, c)
}

// transformContainerPorts converts Docker' inspect output JSON into a more friendly and testable format
func transformContainer(t *testing.T, c inspectOutput) ContainerInspect {
	name := strings.TrimLeft(c.Name, "/")
	ports := transformContainerPorts(t, c)
	created, err := time.Parse(time.RFC3339Nano, c.Created)
	require.NoError(t, err)

	inspect := ContainerInspect{
		ID:       c.Id,
		Name:     name,
		Created:  created,
		Status:   c.State.Status,
		Running:  c.State.Running,
		ExitCode: c.State.ExitCode,
		Error:    c.State.Error,
		Ports:    ports,
	}

	return inspect
}

// transformContainerPorts converts Docker's ports from the following json into a more testable format
// {
//   "80/tcp": [
//     {
// 	     "HostIp": ""
//       "HostPort": "8080"
//     }
//   ]
// }
func transformContainerPorts(t *testing.T, c inspectOutput) []Port {
	var ports []Port

	cPorts := c.NetworkSettings.Ports

	for k, pb := range cPorts {
		split := strings.Split(k, "/")

		containerPort, err := strconv.ParseUint(split[0], 10, 16)
		require.NoError(t, err)

		protocol := split[1]

		for _, p := range pb {
			hostPort, err := strconv.ParseUint(p.HostPort, 10, 16)
			require.NoError(t, err)

			ports = append(ports, Port{
				HostPort:      uint16(hostPort),
				ContainerPort: uint16(containerPort),
				Protocol:      protocol,
			})
		}
	}

	return ports
}
