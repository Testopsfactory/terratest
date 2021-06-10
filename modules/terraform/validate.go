package terraform

import (
	// We alias Golang's native testing package to go_test to avoid naming conflicts with terratest's own testing module
	go_test "testing"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// ValidateAllTerraformModules rutomatically finds all folders specified in RootDir that contain .tf files and runs
// InitAndValidate in all of them.
// Filters down to only those paths passed in ValidationOptions.IncludeDirs, if passed.
// Excludes any folders specified in the ValidationOptions.ExcludeDirs. IncludeDirs will take precedence over ExcludeDirs
// Use the NewValidationOptions method to pass relative paths for either of these options to have the full paths built
// Note that go_test is an alias to Golang's native testing package created to avoid naming conflicts with Terratest's
// own testing package
func ValidateAllTerraformModules(t *go_test.T, opts ValidationOptions) {
	dirsToValidate, readErr := FindTerraformModulePathsInRootE(opts)
	require.NoError(t, readErr)

	for _, dir := range dirsToValidate {
		dir := dir
		t.Run(dir, func(t *go_test.T) {
			t.Parallel()
			tfOpts := &Options{TerraformDir: dir}
			InitAndValidate(t, tfOpts)
		})
	}
}

// Validate calls terraform validate and returns stdout/stderr.
func Validate(t testing.TestingT, options *Options) string {
	out, err := ValidateE(t, options)
	require.NoError(t, err)
	return out
}

// ValidateE calls terraform validate and returns stdout/stderr.
func ValidateE(t testing.TestingT, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, FormatArgs(options, "validate")...)
}

// InitAndValidate runs terraform init and validate with the given options and returns stdout/stderr from the validate command.
// This will fail the test if there is an error in the command.
func InitAndValidate(t testing.TestingT, options *Options) string {
	out, err := InitAndValidateE(t, options)
	require.NoError(t, err)
	return out
}

// InitAndValidateE runs terraform init and validate with the given options and returns stdout/stderr from the validate command.
func InitAndValidateE(t testing.TestingT, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}

	return ValidateE(t, options)
}
