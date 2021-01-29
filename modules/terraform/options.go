package terraform

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/hashicorp/hcl"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
)

var (
	DefaultRetryableTerraformErrors = map[string]string{
		// Helm related terraform calls may fail when too many tests run in parallel. While the exact cause is unknown,
		// this is presumably due to all the network contention involved. Usually a retry resolves the issue.
		".*read: connection reset by peer.*": "Failed to reach helm charts repository.",
		".*transport is closing.*":           "Failed to reach Kubernetes API.",

		// `terraform init` frequently fails in CI due to network issues accessing plugins. The reason is unknown, but
		// eventually these succeed after a few retries.
		".*unable to verify signature.*":                  "Failed to retrieve plugin due to transient network error.",
		".*unable to verify checksum.*":                   "Failed to retrieve plugin due to transient network error.",
		".*no provider exists with the given name.*":      "Failed to retrieve plugin due to transient network error.",
		".*registry service is unreachable.*":             "Failed to retrieve plugin due to transient network error.",
		".*Error installing provider.*":                   "Failed to retrieve plugin due to transient network error.",
		".*Failed to query available provider packages.*": "Failed to retrieve plugin due to transient network error.",
		".*timeout while waiting for plugin to start.*":   "Failed to retrieve plugin due to transient network error.",
		".*timed out waiting for server handshake.*":      "Failed to retrieve plugin due to transient network error.",
		"could not query provider registry for":           "Failed to retrieve plugin due to transient network error.",

		// Provider bugs where the data after apply is not propagated. This is usually an eventual consistency issue, so
		// retrying should self resolve it.
		// See https://github.com/terraform-providers/terraform-provider-aws/issues/12449 for an example.
		".*Provider produced inconsistent result after apply.*": "Provider eventual consistency error.",
	}
)

// Options for running Terraform commands
type Options struct {
	TerraformBinary string // Name of the binary that will be used
	TerraformDir    string // The path to the folder where the Terraform code is defined.

	// The vars to pass to Terraform commands using the -var option. Note that terraform does not support passing `null`
	// as a variable value through the command line. That is, if you use `map[string]interface{}{"foo": nil}` as `Vars`,
	// this will translate to the string literal `"null"` being assigned to the variable `foo`. However, nulls in
	// lists and maps/objects are supported. E.g., the following var will be set as expected (`{ bar = null }`:
	// map[string]interface{}{
	//     "foo": map[string]interface{}{"bar": nil},
	// }
	Vars map[string]interface{}

	VarFiles                 []string               // The var file paths to pass to Terraform commands using -var-file option.
	Targets                  []string               // The target resources to pass to the terraform command with -target
	Lock                     bool                   // The lock option to pass to the terraform command with -lock
	LockTimeout              string                 // The lock timeout option to pass to the terraform command with -lock-timeout
	EnvVars                  map[string]string      // Environment variables to set when running Terraform
	BackendConfig            map[string]interface{} // The vars to pass to the terraform init command for extra configuration for the backend
	RetryableTerraformErrors map[string]string      // If Terraform apply fails with one of these (transient) errors, retry. The keys are a regexp to match against the error and the message is what to display to a user if that error is matched.
	MaxRetries               int                    // Maximum number of times to retry errors matching RetryableTerraformErrors
	TimeBetweenRetries       time.Duration          // The amount of time to wait between retries
	Upgrade                  bool                   // Whether the -upgrade flag of the terraform init command should be set to true or not
	NoColor                  bool                   // Whether the -no-color flag will be set for any Terraform command or not
	SshAgent                 *ssh.SshAgent          // Overrides local SSH agent with the given in-process agent
	NoStderr                 bool                   // Disable stderr redirection
	OutputMaxLineSize        int                    // The max size of one line in stdout and stderr (in bytes)
	Logger                   *logger.Logger         // Set a non-default logger that should be used. See the logger package for more info.
	Parallelism              int                    // Set the parallelism setting for Terraform
	PlanFilePath             string                 // The path to output a plan file to (for the plan command) or read one from (for the apply command)
}

// Clone makes a deep copy of most fields on the Options object and returns it.
//
// NOTE: options.SshAgent and options.Logger CANNOT be deep copied (e.g., the SshAgent struct contains channels and
// listeners that can't be meaningfully copied), so the original values are retained.
func (options *Options) Clone() (*Options, error) {
	newOptions := &Options{}
	if err := copier.Copy(newOptions, options); err != nil {
		return nil, err
	}
	return newOptions, nil
}

// WithDefaultRetryableErrors makes a copy of the Options object and returns an updated object with sensible defaults
// for retryable errors. The included retryable errors are typical errors that most terraform modules encounter during
// testing, and are known to self resolve upon retrying.
// This will fail the test if there are any errors in the cloning process.
func WithDefaultRetryableErrors(t *testing.T, originalOptions *Options) *Options {
	newOptions, err := originalOptions.Clone()
	require.NoError(t, err)

	if newOptions.RetryableTerraformErrors == nil {
		newOptions.RetryableTerraformErrors = map[string]string{}
	}
	for k, v := range DefaultRetryableTerraformErrors {
		newOptions.RetryableTerraformErrors[k] = v
	}

	// These defaults for retry configuration are arbitrary, but have worked well in practice across Gruntwork
	// modules.
	newOptions.MaxRetries = 3
	newOptions.TimeBetweenRetries = 5 * time.Second

	return newOptions
}

// GetVariableAsStringFromVarFile Gets the string represention of a variable from a provided input file found in VarFile
// For list or map, use GetVariableAsListFromVarFile or GetVariableAsMapFromVarFile, respectively.
func GetVariableAsStringFromVarFile(t *testing.T, fileName string, key string) string {
	result, err := GetVariableAsStringFromVarFileE(t, fileName, key)
	require.NoError(t, err)

	return result
}

// GetVariableAsStringFromVarFileE Gets the string represention of a variable from a provided input file found in VarFile
// Will return an error if GetAllVariablesFromVarFileE returns an error or the key provided does not exist in the file.
// For list or map, use GetVariableAsListFromVarFile or GetVariableAsMapFromVarFile, respectively.
func GetVariableAsStringFromVarFileE(t *testing.T, fileName string, key string) (string, error) {
	var variables map[string]interface{}

	err := GetAllVariablesFromVarFileE(t, fileName, &variables)

	if err != nil {
		return "", err
	}

	variable, exists := variables[key]

	if !exists {
		return "", InputFileKeyNotFound{FilePath: fileName, Key: key}
	}

	return fmt.Sprintf("%v", variable), nil
}

// GetVariableAsMapFromVarFile Gets the map represention of a variable from a provided input file found in VarFile
// Note that this returns a map of strings. For maps containing complex types, use GetAllVariablesFromVarFile.
func GetVariableAsMapFromVarFile(t *testing.T, fileName string, key string) map[string]string {
	result, err := GetVariableAsMapFromVarFileE(t, fileName, key)
	require.NoError(t, err)

	return result
}

// GetVariableAsMapFromVarFileE Gets the map represention of a variable from a provided input file found in VarFile.
// Note that this returns a map of strings. For maps containing complex types, use GetAllVariablesFromVarFile
// Returns an error if GetAllVariablesFromVarFileE returns an error, the key provided does not exist, or the value associated with the key is not a map
func GetVariableAsMapFromVarFileE(t *testing.T, fileName string, key string) (map[string]string, error) {
	var variables map[string]interface{}

	resultMap := make(map[string]string)
	err := GetAllVariablesFromVarFileE(t, fileName, &variables)

	if err != nil {
		return nil, err
	}

	variable, exists := variables[key]

	if !exists {
		return nil, InputFileKeyNotFound{FilePath: fileName, Key: key}
	}

	if reflect.TypeOf(variable).String() != "[]map[string]interface {}" {
		return nil, UnexpectedOutputType{Key: key, ExpectedType: "[]map[string]interface {}", ActualType: reflect.TypeOf(variable).String()}
	}

	mapKeys := variable.([]map[string]interface{})

	if len(mapKeys) == 0 {
		return nil, errors.New("no map keys could be found for given map")
	}

	for mapKey, mapVal := range mapKeys[0] {
		resultMap[mapKey] = fmt.Sprintf("%v", mapVal)
	}

	return resultMap, nil
}

// GetVariableAsListFromVarFile Gets the string list represention of a variable from a provided input file found in VarFile
// Note that this returns a list of strings. For lists containing complex types, use GetAllVariablesFromVarFile.
func GetVariableAsListFromVarFile(t *testing.T, fileName string, key string) []string {
	result, err := GetVariableAsListFromVarFileE(t, fileName, key)
	require.NoError(t, err)

	return result
}

// GetVariableAsListFromVarFileE Gets the string list represention of a variable from a provided input file found in VarFile
// Note that this returns a list of strings. For lists containing complex types, use GetAllVariablesFromVarFile.
// Will return error if GetAllVariablesFromVarFileE returns an error, the key provided does not exist, or the value associated with the key is not a list
func GetVariableAsListFromVarFileE(t *testing.T, fileName string, key string) ([]string, error) {
	var variables map[string]interface{}
	resultArray := []string{}
	err := GetAllVariablesFromVarFileE(t, fileName, &variables)

	if err != nil {
		return nil, err
	}

	variable, exists := variables[key]

	if !exists {
		return nil, InputFileKeyNotFound{FilePath: fileName, Key: key}
	}

	if reflect.TypeOf(variable).String() != "[]interface {}" {
		return nil, UnexpectedOutputType{Key: key, ExpectedType: "[]interface {}", ActualType: reflect.TypeOf(variable).String()}
	}

	for _, item := range variable.([]interface{}) {
		resultArray = append(resultArray, fmt.Sprintf("%v", item))
	}

	return resultArray, nil
}

// GetAllVariablesFromVarFile Parses all data from a provided input file found in VarFile and stores the result in the value pointed to by out
func GetAllVariablesFromVarFile(t *testing.T, fileName string, out interface{}) {
	err := GetAllVariablesFromVarFileE(t, fileName, out)
	require.NoError(t, err)
}

// GetAllVariablesFromVarFileE Parses all data from a provided input file found ind in VarFile and stores the result in the value pointed to by out
// Returns an error if the specified file does not exist, the specified file is not readable, or the specified file cannot be decoded from HCL
func GetAllVariablesFromVarFileE(t *testing.T, fileName string, out interface{}) error {
	fileContents, err := ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	err = hcl.Decode(out, string(fileContents))

	if err != nil {
		return HclDecodeError{FilePath: fileName, ErrorText: err.Error()}
	}

	return nil
}
