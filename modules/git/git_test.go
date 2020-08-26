package git

import (
	"testing"

	terratest_testing "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentBranchNameReturnsBranchName(t *testing.T) {
	t.Parallel()

	GetCurrentBranchNameE = func(t terratest_testing.TestingT) (string, error) {
		return "master", nil
	}
	name := GetCurrentBranchName(t)

	assert.Equal(t, name, "master")
}

func TestGetCurrentBranchNameReturnsEmptyForDetachedState(t *testing.T) {
	t.Parallel()

	GetCurrentBranchNameE = func(t terratest_testing.TestingT) (string, error) {
		return "HEAD", nil
	}

	name := GetCurrentBranchName(t)
	assert.Empty(t, name)
}

func TestGetCurrentRefReturnsBranchName(t *testing.T) {
	t.Parallel()

	GetCurrentBranchNameE = func(t terratest_testing.TestingT) (string, error) {
		return "master", nil
	}

	name := GetCurrentGitRef(t)
	assert.Equal(t, "master", name)
}

func TestGetCurrentRefReturns(t *testing.T) {
	t.Parallel()

	GetCurrentBranchNameE = func(t terratest_testing.TestingT) (string, error) {
		return "", nil
	}
	GetCurrentGitRefE = func(t terratest_testing.TestingT) (string, error) {
		return "v0.0.1", nil
	}

	name := GetCurrentGitRef(t)
	assert.Equal(t, "v0.0.1", name)
}
