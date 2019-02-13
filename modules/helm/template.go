package helm

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/gruntwork-cli/errors"
	"github.com/gruntwork-io/gruntwork-cli/files"
	"github.com/stretchr/testify/require"
)

// RenderTemplate runs `helm template` to render the template given the provided options and returns stdout/stderr from
// the template command. If you pass in templateFiles, this will only render those templates. This function will fail
// the test if there is an error rendering the template.
func RenderTemplate(t *testing.T, options *Options, chartDir string, templateFiles []string) string {
	out, err := RenderTemplateE(t, options, chartDir, templateFiles)
	require.NoError(t, err)
	return out
}

// RenderTemplateE runs `helm template` to render the template given the provided options and returns stdout/stderr from
// the template command. If you pass in templateFiles, this will only render those templates.
func RenderTemplateE(t *testing.T, options *Options, chartDir string, templateFiles []string) (string, error) {
	// First, verify the charts dir exists
	absChartDir, err := filepath.Abs(chartDir)
	if err != nil {
		return "", errors.WithStackTrace(err)
	}

	// Now construct the args
	// We first construct the template args
	args := []string{}
	for _, templateFile := range templateFiles {
		// validate this is a valid template file
		absTemplateFile := filepath.Join(absChartDir, templateFile)
		if !files.FileExists(absTemplateFile) {
			return "", errors.WithStackTrace(TemplateFileNotFoundError{Path: templateFile, ChartDir: absChartDir})
		}

		// Note: we only get the abs template file path to check it actually exists, but the `helm template` command
		// expects the relative path from the chart.
		args = append(args, "-x", templateFile)
	}
	// ... and add the chart at the end as the command expects
	args = append(args, chartDir)

	// Finally, call out to helm template command
	return RunHelmCommandAndGetOutputE(t, options, "template", args...)
}
