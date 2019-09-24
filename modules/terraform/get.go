package terraform

import (
	_ "github.com/gruntwork-io/terratest/modules/testing"
)

// Get calls terraform get and return stdout/stderr.
func Get(t TestingT, options *Options) string {
	out, err := GetE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetE calls terraform get and return stdout/stderr.
func GetE(t TestingT, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, "get", "-update")
}
