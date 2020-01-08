package docker

import (
	"crypto/tls"
	"strconv"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/stretchr/testify/require"
)

func TestStop(t *testing.T) {
	t.Parallel()

	// appending timestamp to container name to run tests in parallel
	name := "test-nginx" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// choosing a unique port since 80 may not fly well on test machines
	port := "13030"
	testURL := strings.Join([]string{"http://localhost", port}, ":")

	// for testing the stopping of a docker container
	// we got to run a container first and then stop it
	runOpts := &RunOptions{
		Detach:       true,
		Name:         name,
		Remove:       true,
		OtherOptions: []string{"-p", strings.Join([]string{port, "80"}, ":")},
	}
	Run(t, "nginx:1.17-alpine", runOpts)

	// verify nginx is running
	tlsConfig := &tls.Config{}
	statusCode, _ := http_helper.HttpGet(t, testURL, tlsConfig)
	require.Equal(t, 200, statusCode)

	// try to stop it now
	stopOpts := &StopOptions{}
	out := Stop(t, []string{name}, stopOpts)
	require.Contains(t, out, name)

	// verify nginx is down
	statusCode, _, err := http_helper.HttpGetE(t, testURL, tlsConfig)
	require.NotEmpty(t, err)

}
