package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/stretchr/testify/assert"
)

func TestDockerComposeStdoutExample(t *testing.T) {
	dockerComposeFile := "../examples/docker-compose-stdout-example/docker-compose.yml"

	// Build the Docker image.
	docker.RunDockerCompose(
		t,
		&docker.Options{},
		"-f",
		dockerComposeFile,
		"build",
	)

	// Run the Docker image, read the text file from it, and make sure it contains the expected output.
	output := docker.RunDockerComposeAndGetStdOut(
		t,
		&docker.Options{},
		"-f",
		dockerComposeFile,
		"run",
		"bash_script",
	)

	assert.Contains(t, output, "stdout: message")
	assert.NotContains(t, output, "stderr: error")
}
