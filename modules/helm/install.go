package helm

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/gruntwork-cli/errors"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/files"
)

func Install(t *testing.T, options *Options, chartDir string, releaseName string) {
	require.NoError(t, InstallE(t, options, chartDir, releaseName))
}

func InstallE(t *testing.T, options *Options, chartDir string, releaseName string) error {
	// First, verify the charts dir exists
	absChartDir, err := filepath.Abs(chartDir)
	if err != nil {
		return errors.WithStackTrace(err)
	}
	if !files.FileExists(chartDir) {
		return errors.WithStackTrace(ChartNotFoundError{chartDir})
	}

	// Now call out to helm install to install the charts with the provided options
	args := []string{}
	if options.KubectlOptions != nil && options.KubectlOptions.Namespace != "" {
		args = append(args, "--namespace", options.KubectlOptions.Namespace)
	}
	args, err = getValuesArgsE(t, options, args...)
	if err != nil {
		return err
	}
	args = append(args, "-n", releaseName, absChartDir)
	_, err = RunHelmCommandAndGetOutputE(t, options, "install", args...)
	return err
}
