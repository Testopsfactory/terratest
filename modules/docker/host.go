package docker

import (
	"os"
	"strings"
)

// GetDockerHost returns the name or address of the host on which the Docker engine is running.
func GetDockerHost() string {
	return getDockerHostFromEnv(os.Environ())
}

func getDockerHostFromEnv(env []string) string {
	// Parses the DOCKER_HOST environment variable to find the address
	//
	// For valid formats see:
	// https://github.com/docker/cli/blob/6916b427a0b07e8581d121967633235ced6db9a1/opts/hosts.go#L69
	var dockerUrl []string

	for _, item := range env {
		envVar := strings.Split(item, "=")
		if envVar[0] == "DOCKER_HOST" && len(envVar) == 2 {
			dockerUrl = strings.Split(envVar[1], ":")
			if len(dockerUrl) < 2 {
				// DOCKER_HOST did not split, is empty or ended with ":"
				return "localhost"
			}
			break
		}
	}

	if len(dockerUrl) == 0 {
		// no DOCKER_HOST variable, return default value
		return "localhost"
	}

	switch dockerUrl[0] {
	case "tcp", "ssh", "fd":
		return strings.TrimPrefix(dockerUrl[1], "//")
	default:
		// if DOCKER_HOST is not in one of the formats listed above, return default
		return "localhost"
	}
}
