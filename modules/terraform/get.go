package terraform

import (
	"testing"
)

// Call terraform get and return stdout/stderr
func Get(t *testing.T, options *Options) string {
	out, err := GetE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Call terraform get and return stdout/stderr
func GetE(t *testing.T, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, "get", "-update")
}
