package docker

import (
	"fmt"
	"os"
	"testing"
)

func TestGetDockerHost(t *testing.T) {

	tests := []struct {
		Input    string
		Expected string
	}{
		{
			"unix:///var/run/docker.sock",
			"localhost",
		},
		{
			"npipe:////./pipe/docker_engine",
			"localhost",
		},
		{
			"tcp://1.2.3.4:1234",
			"1.2.3.4",
		},
		{
			"tcp://1.2.3.4",
			"1.2.3.4",
		},
		{
			"ssh://1.2.3.4:22",
			"1.2.3.4",
		},
		{
			"fd://1.2.3.4:1234",
			"1.2.3.4",
		},
		{
			"",
			"localhost",
		},
		{
			"invalidValue",
			"localhost",
		},
		{
			"invalid::value::with::semicolons",
			"localhost",
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("GetDockerHost: %s", test.Input), func(t *testing.T) {
			// GetHost() uses the DOCKER_HOST environment variable, so we need to set
			// it to our test input and then reset it afterwards for other tests
			defer os.Setenv("DOCKER_HOST", os.Getenv("DOCKER_HOST"))
			os.Setenv("DOCKER_HOST", test.Input)

			host := GetDockerHost()

			if host != test.Expected {
				t.Fatalf("Error: expected %s, got %s", test.Expected, host)
			}
		})
	}
}
