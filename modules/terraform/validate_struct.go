package terraform

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/mattn/go-zglob"
)

type ValidationOptions struct {
	// The target directory to recursively search for all Terraform directories (those that contain .tf files)
	// If you provide RootDir and do not pass entries in either IncludeDirs or ExcludeDirs, then all Terraform directories
	// From the RootDir, recursively, will be validated
	RootDir string
	// If you only want to include certain sub directories of RootDir, add the absolute paths here. For example, if the
	// RootDir is /home/project and you want to only include /home/project/examples, add /home/project/examples here
	// Note that while the struct requires full paths, you can pass relative paths to the NewValidationOptions function
	// which will build the full paths based on the supplied RootDir
	IncludeDirs []string
	// If you want to explicitly exclude certain sub directories of RootDir, add the absolute paths here. For example, if the
	// RootDir is /home/project and you want to include everything EXCEPT /home/project/modules, add
	// /home/project/modules to this slice. Note that ExcludeDirs is only considered when IncludeDirs is not passed
	// Note that while the struct requires full paths, you can pass relative paths to the NewValidationOptions function
	// which will build the full paths based on the supplied RootDir
	ExcludeDirs []string
}

// NewValidationOptions returns a ValidationOptions struct, with override-able sane defaults. Note that the
// ValidationOptions's fields IncludeDirs and ExcludeDirs must be absolute paths, but this method will accept relative paths
// and build the absolute paths when instantiating the ValidationOptions struct,  making it the preferred means of configuring
// ValidationOptions.
//
// For example, if your RootDir is /home/project/ and you want to exclude "modules" and "test" you need
// only pass the relative paths in your excludeDirs slice like so:
// opts, err := NewValidationOptions("/home/project", []string{}, []string{"modules", "test"})
func NewValidationOptions(rootDir string, includeDirs, excludeDirs []string) (ValidationOptions, error) {
	vo := ValidationOptions{
		RootDir:     "",
		IncludeDirs: []string{},
		ExcludeDirs: []string{},
	}

	if rootDir == "" {
		return vo, ValidationUndefinedRootDirErr{}
	}

	vo.RootDir = rootDir

	if len(includeDirs) > 0 {
		vo.IncludeDirs = buildFullPathsFromRelative(vo.RootDir, includeDirs)
	}

	if len(excludeDirs) > 0 {
		vo.ExcludeDirs = buildFullPathsFromRelative(vo.RootDir, excludeDirs)
	}

	return vo, nil
}

func buildFullPathsFromRelative(rootDir string, relativePaths []string) []string {
	var fullPaths []string
	for _, relativePath := range relativePaths {
		fullPaths = append(fullPaths, filepath.Join(rootDir, relativePath))
	}
	return fullPaths
}

// FindTerraformModulePathsInRootE returns a slice strings representing the filepaths for all valid Terraform modules
// in the given RootDir, subject to the include / exclude filters.
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

	if len(opts.IncludeDirs) > 0 {
		terraformDirs = collections.ListIntersection(terraformDirs, opts.IncludeDirs)
	}

	if len(opts.ExcludeDirs) > 0 {
		terraformDirs = collections.ListSubtract(terraformDirs, opts.ExcludeDirs)
	}

	// Filter out any filepaths that were explicitly included in opts.ExcludeDirs
	return terraformDirs, nil
}

// Custom error types

// ValidationUndefinedRootDirErr is returned when NewValidationOptions is called without a RootDir argument
type ValidationUndefinedRootDirErr struct{}

func (e ValidationUndefinedRootDirErr) Error() string {
	return "RootDir must be defined in ValidationOptions passed to ValidateAllTerraformModules"
}
