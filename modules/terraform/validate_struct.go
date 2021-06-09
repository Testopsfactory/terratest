package terraform

import (
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/collections"
)

type ValidationOptions struct {
	RootDir string
	//Sub-directories relative to the RootDir to walk recursively for all Terraform modules. For example, if your
	// RootDir is /home/project and you want to validate all modules in home/project/terraform then pass "terraform"
	// in TargetSubDirs
	TargetSubDirs []string
	//ExcludeDirs must be relative to the RootDir. For example if your RootDir is home/project and you want to exclude
	// the path /home/project/test, pass "test" in ExcludeDirs
	ExcludeDirs []string
}

// NewValidationOptions returns a ValidationOptions struct, with override-able sane defaults
func NewValidationOptions(rootDir string, targetSubDirs, excludeDirs []string) (ValidationOptions, error) {
	vo := ValidationOptions{
		ExcludeDirs: []string{},
	}

	if rootDir == "" {
		return vo, ValidationUndefinedRootDirErr{}
	}
	vo.RootDir = rootDir

	// If no target sub directories are passed, default to recursively searching "modules" and "examples"
	if len(targetSubDirs) == 0 {
		vo.TargetSubDirs = []string{
			"modules",
			"examples",
		}
	} else {
		vo.TargetSubDirs = targetSubDirs
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

type ValidationUndefinedRootDirErr struct{}

func (e ValidationUndefinedRootDirErr) Error() string {
	return "RootDir must be defined in ValidationOptions passed to ValidateAllTerraformModules"
}

// readModuleAndExampleSubDirs returns a slice strings representing the filepaths for all valid Terraform modules
// in both the "modules" directory and "examples" directories in the project root, if they exist.
func ReadModuleAndExampleSubDirs(opts ValidationOptions) ([]string, error) {
	var terraformModuleCandidates []string
	// We want to run InitAndValidate on all valid subdirectories of both the modules and examples dirs

	for _, dir := range opts.TargetSubDirs {
		target := filepath.Join(opts.RootDir, dir)
		err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if yes, err := IsTerraformModuleDirectory(path); err == nil && yes {
				terraformModuleCandidates = append(terraformModuleCandidates, path)
			}
			return nil
		})
		if err != nil {
			return terraformModuleCandidates, err
		}
	}

	// Filter out any filepaths that were explicitly included in opts.ExcludeDirs
	return collections.ListSubtract(terraformModuleCandidates, opts.ExcludeDirs), nil
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
