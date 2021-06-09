package terraform

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestInitAndValidateWithNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir: testFolder,
	}

	out := InitAndValidate(t, options)
	require.Contains(t, out, "The configuration is valid")
}

func TestInitAndValidateWithError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-plan-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir: testFolder,
	}

	out, err := InitAndValidateE(t, options)
	require.Error(t, err)
	require.Contains(t, out, "Reference to undeclared input variable")
}

func TestReadModuleAndExampleSubDirsRejectsEmptyOpts(t *testing.T) {
	opts := ValidationOptions{
		RootDir:     "",
		ExcludeDirs: []string{},
	}

	_, err := readModuleAndExampleSubDirs(opts)
	require.Error(t, err)
}

func TestReadModuleAndExampleSubDirsExamples(t *testing.T) {
	cwd, cwdErr := os.Getwd()
	require.NoError(t, cwdErr)

	opts := ValidationOptions{
		RootDir:     filepath.Join(cwd, "../../"),
		ExcludeDirs: []string{},
	}

	subDirs, err := readModuleAndExampleSubDirs(opts)
	require.NoError(t, err)
	// There are many valid Terraform modules in the root/examples directory of the Terratest project, so we should get back many results
	require.Greater(t, len(subDirs), 0)
}

// Verify ExcludeDirs is working properly, by explicitly passing a list of two modules and two examples to exclude
// and ensuring at the end that they do not appear in the returned slice of sub directories to validate
func TestReadModuleAndExampleSubDirsWithResultsExclusion(t *testing.T) {

	cwd, cwdErr := os.Getwd()
	require.NoError(t, cwdErr)

	projectRootDir := filepath.Join(cwd, "../..")

	exclusions := []string{
		filepath.Join(projectRootDir, "retry"),
		filepath.Join(projectRootDir, "ssh"),
		filepath.Join(projectRootDir, "terraform-packer-example"),
		filepath.Join(projectRootDir, "terraform-hello-world-example"),
	}

	opts := ValidationOptions{
		RootDir:     projectRootDir,
		ExcludeDirs: exclusions,
	}

	subDirs, err := readModuleAndExampleSubDirs(opts)
	require.NoError(t, err)
	require.Greater(t, len(subDirs), 0)
	for _, exclusion := range exclusions {
		require.False(t, collections.ListContains(subDirs, exclusion))
	}
}
