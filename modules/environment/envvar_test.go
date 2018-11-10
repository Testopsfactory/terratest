package environment

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests can not run in parallel, since they manipulate env vars
var envvarList = []string{
	"TERRATEST_TEST_ENVIRONMENT",
	"TERRATESTTESTENVIRONMENT",
	"TERRATESTENVIRONMENT",
}

func TestGetFirstNonEmptyEnvVarOrEmptyStringChecksInOrder(t *testing.T) {

	os.Setenv("TERRATESTTESTENVIRONMENT", "test")
	os.Setenv("TERRATESTENVIRONMENT", "circleCI")
	defer os.Setenv("TERRATESTTESTENVIRONMENT", "")
	defer os.Setenv("TERRATESTENVIRONMENT", "")
	value := GetFirstNonEmptyEnvVarOrEmptyString(t, envvarList)
	assert.Equal(t, value, "test")
}

func TestGetFirstNonEmptyEnvVarOrEmptyStringReturnsEmpty(t *testing.T) {
	value := GetFirstNonEmptyEnvVarOrEmptyString(t, envvarList)
	assert.Equal(t, value, "")
}
