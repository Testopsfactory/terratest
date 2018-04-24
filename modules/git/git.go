package git

import (
	"testing"
	"os/exec"
	"strings"
)

func GetCurrentBranchName(t *testing.T) string {
	out, err := GetCurrentBranchNameE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func GetCurrentBranchNameE(t *testing.T) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytes)), nil
}
