package test_structure

import (
	"os"
	"testing"
	"encoding/json"
	"io/ioutil"
	"github.com/gruntwork-io/terratest"
	"path/filepath"
	"log"
)

const SkipSetupEnvVar = "SKIP_SETUP"
const SkipTeardownEnvVar = "SKIP_TEARDOWN"
const SkipValidationEnvVar = "SKIP_VALIDATION"

// Run the given function to perform setup steps for this test unless the SKIP_SETUP env var is set, in which case
// setup will be skipped.
func Setup(logger *log.Logger, setup func()) {
	runFunctionIfEnvVarNotSet(logger, SkipSetupEnvVar, "setup", setup)
}

// Run the given function to perform validation steps for this test unless the SKIP_VALIDATION env var is set, in which
// case validation will be skipped.
func Validation(logger *log.Logger, validation func()) {
	runFunctionIfEnvVarNotSet(logger, SkipValidationEnvVar, "validation", validation)
}

// Run the given function to perform teardown steps for this test unless the SKIP_TEARDOWN var is set, in which case
// teardown will be skipped.
func Teardown(logger *log.Logger, teardown func()) {
	runFunctionIfEnvVarNotSet(logger, SkipTeardownEnvVar, "teardown", teardown)
}

// If the given environment variable is not set, run the given function
func runFunctionIfEnvVarNotSet(logger *log.Logger, envVarName string, functionName string, function func()) {
	if os.Getenv(envVarName) == "" {
		logger.Printf("The %s environment variable is not set, so running %s.", envVarName, functionName)
		function()
	} else {
		logger.Printf("The %s environment variable is set, so skipping %s.", envVarName, functionName)
	}
}

// Serialize and save TerratestOptions into the given folder. This allows you to create TerratestOptions during setup
// and to reuse that TerratestOptions later during validation and teardown.
func SaveTerratestOptions(t *testing.T, testFolder string, terratestOptions *terratest.TerratestOptions) {
	SaveTestData(t, formatTerratestOptionsPath(testFolder), terratestOptions)
}

// Load and unserialize TerratestOptions from the given folder. This allows you to reuse a TerratestOptions that was
// created during an earlier setup steps in later validation and teardown steps.
func LoadTerratestOptions(t *testing.T, testFolder string) *terratest.TerratestOptions {
	var terratestOptions *terratest.TerratestOptions
	LoadTestData(t, formatTerratestOptionsPath(testFolder), terratestOptions)
	return terratestOptions
}

// Format a path to save TerratestOptions in the given folder
func formatTerratestOptionsPath(testFolder string) string {
	return filepath.Join(testFolder, "TerratestOptions.json")
}

// Serialize and save a value used at test time to the given path. This allows you to create some sort of test data
// (e.g., TerratestOptions) during setup and to reuse this data later during validation and teardown.
func SaveTestData(t *testing.T, path string, value interface{}) {
	bytes, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Failed to convert value %s to JSON: %v", path, err)
	}

	if err := ioutil.WriteFile(path, bytes, 0644); err != nil {
		t.Fatalf("Failed to save value %s: %v", path, err)
	}
}

// Load and unserialize a value stored at the given path. The value should be a pointer to a struct into which the
// value will be deserialized. This allows you to reuse some sort of test data (e.g., TerratestOptions) from earlier
// setup steps in later validation and teardown steps.
func LoadTestData(t *testing.T, path string, value interface{}) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to load value from %s: %v", path, err)
	}

	if err := json.Unmarshal(bytes, value); err != nil {
		t.Fatalf("Failed to parse JSON for value %s: %v", path, err)
	}
}
