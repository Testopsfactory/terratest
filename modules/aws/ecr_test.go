package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcrRepo(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	repo1, err := CreateECRRepoE(t, region, "terratest")
	defer DeleteECRRepo(t, region, repo1)

	assert.Nil(t, err)
	assert.Equal(t, "terratest", *repo1.RepositoryName)

	repo2, err := GetECRRepoE(t, region, "terratest")

	assert.Nil(t, err)
	assert.Equal(t, "terratest", *repo2.RepositoryName)
}
