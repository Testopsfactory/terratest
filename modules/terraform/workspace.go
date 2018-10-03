package terraform

import (
	"strings"
	"testing"
)

// WorkspaceSelectOrNew runs terraform workspace with the given options and returns workspace name.
// It tries to select a workspace with the given name, or it creates a new one if it doesn't exist.
func WorkspaceSelectOrNew(t *testing.T, options *Options, name string) string {
	out, err := WorkspaceSelectOrNewE(t, options, name)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// WorkspaceSelectOrNewE runs terraform workspace with the given options and returns workspace name.
// It tries to select a workspace with the given name, or it creates a new one if it doesn't exist.
func WorkspaceSelectOrNewE(t *testing.T, options *Options, name string) (string, error) {
	out, err := RunTerraformCommandE(t, options, "workspace", "list")
	if err != nil {
		return out, nil
	}

	if strings.Contains(out, name) {
		_, err = RunTerraformCommandE(t, options, "workspace", "select", name)
	} else {
		_, err = RunTerraformCommandE(t, options, "workspace", "new", name)
	}
	if err != nil {
		return out, nil
	}

	return RunTerraformCommandE(t, options, "workspace", "show")
}
