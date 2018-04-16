package test

import (
	"testing"
	"github.com/gruntwork-io/terratest/packer"
	"github.com/gruntwork-io/terratest/aws"
	"fmt"
	"github.com/gruntwork-io/terratest/util"
	"github.com/gruntwork-io/terratest/docker"
	"time"
	"github.com/gruntwork-io/terratest/http"
	"strconv"
)

// An example of how to test the Packer template in examples/packer-docker-example completely locally using Terratest
// and Docker.
func PackerDockerExampleLocalTest(t *testing.T)  {
	t.Parallel()

	packerOptions := packer.Options {
		// The path to where the Packer template is located
		Template: "../examples/packer-docker-example/build.json",

		// Only build the Docker image for local testing
		Only: "ubuntu-docker",
	}

	// Build the Docker image using Packer
	packer.BuildAmi(t, packerOptions)

	// Configure the port the web app will listen on and the text it will return using environment variables
	serverPort := 8080
	expectedServerText := fmt.Sprintf("Hello, %s!", util.UniqueId())
	envVars := map[string]string{
		"SERVER_PORT": strconv.Itoa(serverPort),
		"SERVER_TEXT": expectedServerText,
	}

	// Run Docker Compose to fire up the web app. We run it in the background (-d) so it doesn't block this test.
	docker.RunDockerCompose(t, "../examples-packer-docker-example", envVars, "up", "-d")

	// Make sure to shut down the Docker container at the end of the test
	defer docker.RunDockerCompose(t, "../examples-packer-docker-example", envVars, "down")

	// It can take a few seconds for the Docker container boot up, so retry a few times
	maxRetries := 5
	timeBetweenRetries := 2 * time.Second
	url := fmt.Sprint("http://localhost:%d", serverPort)

	// Verify that we get back a 200 OK with the expected text
	http_helper.HttpGetWithRetry(t, url, 200, expectedServerText, maxRetries, timeBetweenRetries)
}
