package terraform

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/collections"
)

type ValidationOptions struct {
	RootDir     string
	ExcludeDirs []string
}

type ValidationUndefinedRootDirErr struct{}

func (e ValidationUndefinedRootDirErr) Error() string {
	return "RootDir must be defined in ValidationOptions passed to ValidateAllTerraformModules"
}

// readModuleAndExampleSubDirs returns a slice strings representing the filepaths for all valid Terraform modules
// in both the "modules" directory and "examples" directories in the project root, if they exist.
func readModuleAndExampleSubDirs(opts ValidationOptions) ([]string, error) {
	if opts.RootDir == "" {
		return nil, ValidationUndefinedRootDirErr{}
	}
	var validationCandidates []string
	// We want to run InitAndValidate on all valid subdirectories of both the modules and examples dirs
	modulesDir := filepath.Join(opts.RootDir, "modules")
	examplesDir := filepath.Join(opts.RootDir, "examples")

	moduleSubDirs, readModulesErr := os.ReadDir(modulesDir)
	if readModulesErr != nil {
		return nil, readModulesErr
	}
	exampleSubDirs, readExamplesErr := os.ReadDir(examplesDir)
	if readExamplesErr != nil {
		return nil, readExamplesErr
	}

	for _, m := range filterTerraformModulesFromDirs(modulesDir, moduleSubDirs) {
		validationCandidates = append(validationCandidates, m)
	}

	for _, e := range filterTerraformModulesFromDirs(examplesDir, exampleSubDirs) {
		validationCandidates = append(validationCandidates, e)
	}

	// Filter out any filepaths that were explicitly included in opts.ExcludeDirs
	return collections.ListSubtract(validationCandidates, opts.ExcludeDirs), nil
}

// filterTerraformModulesFromDirs accepts a slice of fs.DirEntry representing subDirectories
// (under "modules" or "examples"), for instance, returning only those that contain a main.tf file in their root. This
// is useful for filtering out any sub directories that might ship alongside Terraform modules, but actually be
// Terraform modules themselves
func filterTerraformModulesFromDirs(rootDir string, subDirs []fs.DirEntry) []string {
	var validTerraformModules []string
	for _, m := range subDirs {
		maybeMainTf := filepath.Join(rootDir, m.Name(), "main.tf")
		if _, err := os.Stat(maybeMainTf); err == nil {
			validTerraformModules = append(validTerraformModules, filepath.Join(rootDir, m.Name()))
		}
	}
	return validTerraformModules
}
