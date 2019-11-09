package docker

import (
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
	"testing"
)

// BuildOptions defines options that can be passed to the 'docker build' command.
type BuildOptions struct {
	File      string
	Tags      []string
	BuildArgs []string
}

// Build runs the 'docker build' command at the given path with the given options and fails the test if there are any
// errors.
func Build(t *testing.T, path string, options *BuildOptions) {
	require.NoError(t, BuildE(t, path, options))
}

// BuildE runs the 'docker build' command at the given path with the given options and returns any errors.
func BuildE(t *testing.T, path string, options *BuildOptions) error {
	logger.Logf(t, "Running Docker build on Dockerfile %s", options.File)

	args, err := formatDockerBuildArgs(path, options)
	if err != nil {
		return err
	}

	cmd := shell.Command{
		Command: "docker",
		Args:    args,
	}

	_, buildErr := shell.RunCommandAndGetOutputE(t, cmd)
	return buildErr
}

// formatDockerBuildArgs formats the arguments for the 'docker build' command.
func formatDockerBuildArgs(path string, options *BuildOptions) ([]string, error) {
	args := []string{"build"}

	if options.File != "" {
		args = append(args, "--file", options.File)
	}

	for _, tag := range options.Tags {
		args = append(args, "--tag", tag)
	}

	for _, arg := range options.BuildArgs {
		args = append(args, "--build-arg", arg)
	}

	args = append(args, path)

	return args, nil
}
