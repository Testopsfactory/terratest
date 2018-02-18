// This package contains methods that help us structure our test code.
//
// Motivation: many of our tests involve:
//
// 1. Setup: e.g., building an AMI or Docker image, configuring TerratestOptions, deploying infrastructure with Terraform
// 2. Teardown: e.g., undeploying infrastructure with Terraform
// 3. Validation: e.g., checking the infrastructure works as expected
//
// Typical test case:
//
// func TestExample(t *testing.T) {
//   testPath := "../examples/foo"
//   logger := terralog.NewLogger("TestExample")
//
//   amiId := buildAmi(t)                                                 // setup
//   terratestOptions := createOptions(t, amiId, testPath)                // setup
//   deployInfrastructureWithTerraform(t, terratestOptions)               // setup
//
//   defer undeployInfrastructureWithTerraform(t, terratestOptions)       // teardown
//
//   testInfrastructureWorks(t, terratestOptions)                         // validation
// }
//
// The setup and teardown steps can be very slow (on the order of 1-10 minutes), whereas the validation steps are
// typically fast. In the dev environment, you often want to be able to iterate quickly on the validation step without
// having to go through the entire setup and teardown process from scratch every time.
//
// To make this possible, you can use the methods in this package to structure your test case as follows:
//
// func TestExample(t *testing.T) {
//   testPath := "../examples/foo"
//   logger := terralog.NewLogger("TestExample")
//
//   test_structure.Setup(logger, func() {
//     amiId := buildAmi(t)                                                 // setup
//     terratestOptions := createOptions(t, amiId, testPath)                // setup
//     deployInfrastructureWithTerraform(t, terratestOptions)               // setup
//     test_structure.SaveTerratestOptions(t, testPath, terratestOptions)   // save TerratestOptions for later steps
//   })
//
//   test_structure.Teardown(logger, func() {
//     terratestOptions := test_structure.LoadTerratestOptions(t, testPath) // load TerratestOptions from earlier setup
//     defer undeployInfrastructureWithTerraform(t, terratestOptions)       // teardown
//   })
//
//   test_structure.Validation(logger, func() {
//     terratestOptions := test_structure.LoadTerratestOptions(t, testPath) // load TerratestOptions from earlier setup
//     testInfrastructureWorks(t, terratestOptions)                         // validation
//   })
// }
//
// Now, in the dev environment, the workflow you can use is:
//
// 1. Do the initial setup (just once): SKIP_VALIDATION=true SKIP_TEARDOWN=true go test -run TestExample
// 2. Do your validation (as many times as you want): SKIP_SETUP=true SKIP_TEARDOWN=true go test -run TestExample
// 3. Do the teardown (just once): SKIP_SETUP=true SKIP_VALIDATION=true go test -run TestExample
//
// This way, you only pay the cost of setup and teardown once and you can do as many iterations on validation in
// between as you want.
//
// In the CI environment, none of the SKIP_XXX env vars will be set, so all steps will execute from start to finish.
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
