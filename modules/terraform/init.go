package terraform

import (
	"testing"
)

// Init calls terraform init and return stdout/stderr.
func Init(t *testing.T, options *Options) string {
	out, err := InitE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// InitE calls terraform init and return stdout/stderr.
func InitE(t *testing.T, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, "init")
}
