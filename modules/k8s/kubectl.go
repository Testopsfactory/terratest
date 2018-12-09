package k8s

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/shell"
)

// RunKubectl will call kubectl using the provided options and args, failing the test on error.
func RunKubectl(t *testing.T, options *KubectlOptions, args ...string) {
	err := RunKubectlE(t, options, args...)
	if err != nil {
		t.Fatal(err)
	}
}

// RunKubectlE will call kubectl using the provided options and args.
func RunKubectlE(t *testing.T, options *KubectlOptions, args ...string) error {
	_, err := RunKubectlAndGetOutputE(t, options, args...)
	return err
}

// RunKubectlAndGetOutputE will call kubectl using the provided options and args, returning the output of stdout and
// stderr.
func RunKubectlAndGetOutputE(t *testing.T, options *KubectlOptions, args ...string) (string, error) {
	cmdArgs := []string{}
	if options.ContextName != "" {
		cmdArgs = append(cmdArgs, "--context", options.ContextName)
	}
	if options.ConfigPath != "" {
		cmdArgs = append(cmdArgs, "--kubeconfig", options.ConfigPath)
	}
	cmdArgs = append(cmdArgs, args...)
	command := shell.Command{
		Command: "kubectl",
		Args:    cmdArgs,
		Env:     options.Env,
	}
	return shell.RunCommandAndGetOutputE(t, command)
}
