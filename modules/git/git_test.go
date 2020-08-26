package git

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var branchNameMock string
var tagMock string

func TestGetCurrentBranchNameReturnsBranchName(t *testing.T) {
	branchNameMock = "master"
	name := GetCurrentBranchName(t)

	assert.Equal(t, "master", name)
}

func TestGetCurrentBranchNameReturnsEmptyForDetachedState(t *testing.T) {
	branchNameMock = "HEAD"
	name := GetCurrentBranchName(t)

	assert.Empty(t, name)
}

func TestGetCurrentRefReturnsBranchName(t *testing.T) {
	branchNameMock = "master"
	name := GetCurrentGitRef(t)

	assert.Equal(t, "master", name)
}

func TestGetCurrentRefReturnsTagValueForEmptyBranchName(t *testing.T) {
	branchNameMock = ""
	tagMock = "v0.0.1"
	name := GetCurrentGitRef(t)

	assert.Equal(t, "v0.0.1", name)
}

// Mock function for branchNameExecCommand
func branchNameExecCommandMock(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestExecCommandHelper", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{
		"GO_WANT_HELPER_PROCESS=1",
		"STDOUT=" + branchNameMock,
	}
	return cmd
}

// Mock function for tagExecCommand
func tagExecCommandMock(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestExecCommandHelper", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{
		"GO_WANT_HELPER_PROCESS=1",
		"STDOUT=" + tagMock,
	}
	return cmd
}

func TestExecCommandHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprintf(os.Stdout, os.Getenv("STDOUT"))
	os.Exit(0)
}

// Setup and teardown
func TestMain(m *testing.M) {
	branchNameExecCommand = branchNameExecCommandMock
	tagExecCommand = tagExecCommandMock
	defer func() { branchNameExecCommand = exec.Command }()
	defer func() { tagExecCommand = exec.Command }()

	code := m.Run()

	os.Exit(code)
}
