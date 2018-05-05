package terraform

import (
	"testing"
)

// Call terraform init and return stdout/stderr
func Init(t *testing.T, options *Options) string {
	out, err := InitE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Call terraform init and return stdout/stderr
func InitE(t *testing.T, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, "init")
}

