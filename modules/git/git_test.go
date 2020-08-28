package git

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentBranchNameReturnsBranchName(t *testing.T) {
	err := exec.Command("git", "checkout", "master").Run()
	assert.Nil(t, err)

	name := GetCurrentBranchName(t)

	assert.Equal(t, "master", name)
}

func TestGetCurrentBranchNameReturnsEmptyForDetachedState(t *testing.T) {
	err := exec.Command("git", "checkout", "v0.0.1").Run()
	assert.Nil(t, err)

	name := GetCurrentBranchName(t)

	assert.Empty(t, name)
}

func TestGetCurrentRefReturnsBranchName(t *testing.T) {
	err := exec.Command("git", "checkout", "master").Run()
	assert.Nil(t, err)

	name := GetCurrentGitRef(t)

	assert.Equal(t, "master", name)
}

func TestGetCurrentRefReturnsTagValue(t *testing.T) {
	err := exec.Command("git", "checkout", "v0.0.1").Run()
	assert.Nil(t, err)

	name := GetCurrentGitRef(t)

	assert.Equal(t, "v0.0.1", name)
}

func TestGetCurrentRefReturnsLightTagValue(t *testing.T) {
	err := exec.Command("git", "checkout", "58d3ea8").Run()
	assert.Nil(t, err)

	name := GetCurrentGitRef(t)

	assert.Equal(t, "v0.0.1-1-g58d3ea8", name)
}

func setup(tempDirName string) {
	url := "https://github.com/gruntwork-io/terratest.git"
	err := exec.Command("git", "clone", url, tempDirName).Run()
	if err != nil {
		fmt.Println("Test setup - Terratest git repository clone failed")
		os.Exit(1)
	}

	err = os.Chdir(tempDirName)
	if err != nil {
		fmt.Println("Test setup - Change directory failed")
		os.Exit(1)
	}
}

func teardown(tempDirName string) {
	err := os.RemoveAll("../" + tempDirName + "/")
	if err != nil {
		fmt.Println("Test teardown - remove temp directory failed")
	}
}

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	tempDirName := "temp_terratest_clone"
	setup(tempDirName)
	defer teardown(tempDirName)

	return m.Run()
}
