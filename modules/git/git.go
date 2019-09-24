// Package git allows to interact with Git.
package git

import (
	"os/exec"
	"strings"

	_ "github.com/gruntwork-io/terratest/modules/testing"
)

// GetCurrentBranchName retrieves the current branch name.
func GetCurrentBranchName(t TestingT) string {
	out, err := GetCurrentBranchNameE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetCurrentBranchNameE retrieves the current branch name.
func GetCurrentBranchNameE(t TestingT) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytes)), nil
}
