package terraform

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/mattn/go-zglob"
)

type ValidationOptions struct {
	RootDir string
	//Sub-directories relative to the RootDir to walk recursively for all Terraform modules. For example, if your
	// RootDir is /home/project and you want to validate all modules in home/project/terraform then pass "terraform"
	// in IncludeDirs
	IncludeDirs []string
	//ExcludeDirs must be relative to the RootDir. For example if your RootDir is home/project and you want to exclude
	// the path /home/project/test, pass "test" in ExcludeDirs
	ExcludeDirs []string
}

// NewValidationOptions returns a ValidationOptions struct, with override-able sane defaults
func NewValidationOptions(rootDir string, includeDirs, excludeDirs []string) (ValidationOptions, error) {
	vo := ValidationOptions{
		ExcludeDirs: []string{},
	}

	if rootDir == "" {
		return vo, ValidationUndefinedRootDirErr{}
	}
	vo.RootDir = rootDir

	// If no target sub directories are passed, default to recursively searching "modules" and "examples"
	if len(includeDirs) == 0 {
		vo.IncludeDirs = []string{
			"modules",
			"examples",
		}
	} else {
		vo.IncludeDirs = includeDirs
	}

	if len(excludeDirs) > 0 {
		var fullExclusionPaths []string
		for _, excludedPath := range excludeDirs {
			fullExclusionPaths = append(fullExclusionPaths, filepath.Join(rootDir, excludedPath))
		}
		vo.ExcludeDirs = fullExclusionPaths
	}

	return vo, nil
}

// readModuleAndExampleSubDirs returns a slice strings representing the filepaths for all valid Terraform modules
// in both the "modules" directory and "examples" directories in the project root, if they exist.
func FindTerraformModulePathsInRootE(opts ValidationOptions) ([]string, error) {
	// Find all Terraform files from the configured RootDir
	pattern := fmt.Sprintf("%s/**/*.tf", opts.RootDir)
	matches, err := zglob.Glob(pattern)
	if err != nil {
		return matches, err
	}
	// Keep a unique set of the base dirs that contain Terraform files
	terraformDirSet := make(map[string]bool)
	for _, match := range matches {
		// The glob match returns all full paths to every .tf file, whereas we're only interested in their root
		// directories for the purposes of running Terraform validate
		rootDir := path.Dir(match)
		terraformDirSet[rootDir] = true
	}

	// Return the unique slice of Terraform directories found starting at opts.RootDir
	terraformDirs := make([]string, 0, len(terraformDirSet))
	for dir := range terraformDirSet {
		terraformDirs = append(terraformDirs, dir)
	}

	// Filter out any filepaths that were explicitly included in opts.ExcludeDirs
	return collections.ListSubtract(terraformDirs, opts.ExcludeDirs), nil
}

// IsTerraformModuleDirectory accepts a slice of string paths representing sub directories (under "modules" or "examples"),
// for instance, returning true for any.tf files. This is useful for filtering out any sub directories that might ship
// alongside Terraform modules, but not actually be Terraform modules themselves
func IsTerraformModuleDirectory(path string) (bool, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	for _, f := range files {
		if filepath.Ext(filepath.Join(path, f.Name())) == ".tf" {
			return true, nil
		}
	}
	return false, nil
}

// Custom error types

// ValidationUndefinedRootDirErr is returned when NewValidationOptions is called without a RootDir argument
type ValidationUndefinedRootDirErr struct{}

func (e ValidationUndefinedRootDirErr) Error() string {
	return "RootDir must be defined in ValidationOptions passed to ValidateAllTerraformModules"
}
