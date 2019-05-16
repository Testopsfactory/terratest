package terraform

import (
	"fmt"
	"testing"
)

// InitAndPlan runs terraform init and plan with the given options and return stdout/stderr from the apply command.
func InitAndPlan(t *testing.T, options *Options) int {
	exitCode, err := InitAndPlanE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode
}

// InitAndPlanE runs terraform init and plan with the given options and return stdout/stderr from the apply command.
func InitAndPlanE(t *testing.T, options *Options) (int, error) {
	if _, err := InitE(t, options); err != nil {
		return DefaultErrorExitCode, err
	}

	return PlanExitCodeE(t, options)
}

// PlanExitCode runs terraform apply with the given options and returns the detailed exitcode.
func PlanExitCode(t *testing.T, options *Options) int {
	exitCode, err := PlanExitCodeE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode
}

// PlanExitCodeE runs terraform plan with the given options and returns the detailed exitcode.
func PlanExitCodeE(t *testing.T, options *Options) (int, error) {
	return GetExitCodeForTerraformCommandE(t, options, FormatArgs(options, "plan", "-input=false", "-lock=true", "-detailed-exitcode")...)
}

// PlanAllExitCodeTg runs terragrunt plan-all with the given options and returns the detailed exitcode.
func PlanAllExitCodeTg(t *testing.T, options *Options) int {
	exitCode, err := PlanAllExitCodeTgE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode
}

// PlanAllExitCodeTgE runs terragrunt plan-all with the given options and returns the detailed exitcode.
func PlanAllExitCodeTgE(t *testing.T, options *Options) (int, error) {
	if options.TerraformBinary != "terragrunt" {
		return 1, fmt.Errorf("terragrunt must be set as TerraformBinary to use this method")
	}

	return GetExitCodeForTerraformCommandE(t, options, FormatArgs(options, "plan-all", "--input=false", "--lock=true", "--detailed-exitcode")...)
}
