package docker

import (
	"os"
	"strings"
)

// GetDockerHost returns the name or address of the host on which the Docker engine is running.
func GetDockerHost() string {
	// Parses the DOCKER_HOST environment variable to find the address
	//
	// For valid formats see:
	// https://github.com/docker/cli/blob/6916b427a0b07e8581d121967633235ced6db9a1/opts/hosts.go#L69
	dockerUrl := strings.Split(os.Getenv("DOCKER_HOST"), ":")
	switch dockerUrl[0] {
	case "tcp", "ssh", "fd":
		return strings.TrimPrefix(dockerUrl[1], "//")
	default:
		// if DOCKER_HOST is empty, or not in one of the formats listed above
		return "localhost"
	}
}
