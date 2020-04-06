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

// ContainerInspect defines the output of the Inspect method, with the options returned by 'docker inspect'
// converted into a more friendly and testable interface
type ContainerInspect struct {
	// ID of the inspected container
	ID       string

	// Name of the inspected container
	Name     string

	// time.Time that the container was created
	Created  time.Time

	// String representing the container's status
	Status   string

	// Whether the container is currently running or not
	Running  bool

	// Container's exit code
	ExitCode uint8

	// String with the container's error message, if there is any
	Error    string

	// Ports exposed by the container
	Ports    []Port

	// Volume bindings made to the container
	Binds    []VolumeBind
}

// Port represents a single port mapping exported by the container
type Port struct {
	HostPort      uint16
	ContainerPort uint16
	Protocol      string
}

// VolumeBind represents a single volume binding made to the container
type VolumeBind struct {
	Source      string
	Destination string
}

// inspectOutput defines options that will be returned by 'docker inspect', in JSON format.
// Not all options are included here, only the ones that we might need
type inspectOutput struct {
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
	HostConfig struct {
		Binds []string
	}
}

func Inspect(t *testing.T, id string) ContainerInspect {
	out, err := InspectE(t, id)
	require.NoError(t, err)
	return out
}

// InspectE runs the 'docker inspect {container id}' command and returns a ContainerInspect
// struct, converted from the output JSON, along with any errors
func InspectE(t *testing.T, id string) (ContainerInspect, error) {
	cmd := shell.Command{
		Command: "docker",
		Args:    []string{"container", "inspect", id},
	}

	out, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		return ContainerInspect{}, err
	}

	var containers []inspectOutput
	err = json.Unmarshal([]byte(out), &containers)
	if err != nil {
		return ContainerInspect{}, err
	}

	if len(containers) == 0 {
		return ContainerInspect{}, nil
	}

	c := containers[0]

	return transformContainer(t, c), nil
}

// transformContainerPorts converts 'docker inspect' output JSON into a more friendly and testable format
func transformContainer(t *testing.T, c inspectOutput) ContainerInspect {
	name := strings.TrimLeft(c.Name, "/")
	ports := transformContainerPorts(t, c)
	volumes := transformContainerVolumes(t, c)

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
		Binds:    volumes,
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

// transformContainerVolumes converts Docker's volume bindings from the
// format "/foo/bar:/foo/baz" into a more testable one
func transformContainerVolumes(t *testing.T, c inspectOutput) []VolumeBind {
	binds := c.HostConfig.Binds
	volumes := make([]VolumeBind, 0, len(binds))

	for _, b := range binds {
		var source, dest string

		split := strings.Split(b, ":")

		// Considering it a named volume
		source = split[0]
		dest = split[1]

		volumes = append(volumes, VolumeBind{
			Source:      source,
			Destination: dest,
		})
	}

	return volumes
}
