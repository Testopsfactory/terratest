package docker

import (
	"crypto/tls"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
)

func TestStop(t *testing.T) {
	t.Parallel()

	// appending timestamp to container name to run tests in parallel
	name := "test-nginx" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// Parse the DOCKER_HOST environment variable so we know where to connect
	//
	// For valid formats see:
	// https://github.com/docker/cli/blob/6916b427a0b07e8581d121967633235ced6db9a1/opts/hosts.go#L69
	var host string

	dockerUrl := strings.Split(os.Getenv("DOCKER_HOST"), ":")
	switch dockerUrl[0] {
	case "tcp", "ssh", "fd":
		host = strings.TrimPrefix(dockerUrl[1], "//")
	default:
		// if DOCKER_HOST is empty, or not in one of the formats listed above
		host = "localhost"
	}

	// choosing a unique port since 80 may not fly well on test machines
	port := "13030"
	testURL := "http://" + host + ":" + port

	// for testing the stopping of a docker container
	// we got to run a container first and then stop it
	runOpts := &RunOptions{
		Detach:       true,
		Name:         name,
		Remove:       true,
		OtherOptions: []string{"-p", port + ":80"},
	}
	Run(t, "nginx:1.17-alpine", runOpts)

	// verify nginx is running
	http_helper.HttpGetWithRetryWithCustomValidation(t, testURL, &tls.Config{}, 60, 2*time.Second, verifyNginxIsUp)

	// try to stop it now
	out := Stop(t, []string{name}, &StopOptions{})
	require.Contains(t, out, name)

	// verify nginx is down
	// run a docker ps with name filter
	command := shell.Command{
		Command: "docker",
		Args:    []string{"ps", "-q", "--filter", "name=" + name},
	}
	output := shell.RunCommandAndGetStdOut(t, command)
	require.Empty(t, output)
}

func verifyNginxIsUp(statusCode int, body string) bool {
	return statusCode == 200 && strings.Contains(body, "nginx!")
}
