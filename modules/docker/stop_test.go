package docker

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStop(t *testing.T) {
	t.Parallel()

	// appending timestamp to container name to run tests in parallel
	name := "test-nginx" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// for testing the stopping of a docker container
	// we got to run a container first and then stop it
	runOpts := &RunOptions{
		Detach: true,
		Name:   name,
		Remove: true,
	}
	Run(t, "nginx:1.17-alpine", runOpts)

	// try to stop it now
	stopOpts := &StopOptions{}
	out := Stop(t, []string{name}, stopOpts)
	require.Contains(t, out, name)
}
