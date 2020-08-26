// Package git allows to interact with Git.
package git

import (
	"os/exec"
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// GetCurrentBranchName retrieves the current branch name.
func GetCurrentBranchName(t testing.TestingT) string {
	out, err := GetCurrentBranchNameE(t)
	if err != nil {
		t.Fatal(err)
	}
	if out == "HEAD" {
		return ""
	}
	return out
}

// GetCurrentBranchNameE retrieves the current branch name. Created as variable
// to enable mocking within the tests.
var GetCurrentBranchNameE = func(t testing.TestingT) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytes)), nil
}

// GetCurrentGitRef retrieves current branch name or most recent tag from a commit. If the tag
// points to the commit, then only tag is returned.
func GetCurrentGitRef(t testing.TestingT) string {
	out, err := GetCurrentGitRefE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetCurrentGitRefE retrieves current branch name or most recent tag from a commit. If the tag
// points to the commit, then only tag is returned.
func GetCurrentGitRefE(t testing.TestingT) (string, error) {
	out, err := GetCurrentBranchNameE(t)

	if err != nil {
		return "", err
	}

	if out != "" {
		return out, nil
	}

	out, err = GetMostRecentTagE(t)
	if err != nil {
		return "", err
	}
	return out, nil
}

// GetMostRecentTagE retrieves most recent tag that is reachable from a commit. If the tag points
// to the commit, then only the tag is returned. Created as variable to enable mocking within tests.
var GetMostRecentTagE = func(t testing.TestingT) (string, error) {
	cmd := exec.Command("git", "describe", "--tags")
	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytes)), nil
}
