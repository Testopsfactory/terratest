package helm

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/testing"

	"fmt"
	"os"

	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
)

// RenderTemplate runs `helm template` to render the template given the provided options and returns stdout/stderr from
// the template command. If you pass in templateFiles, this will only render those templates. This function will fail
// the test if there is an error rendering the template.
func RenderTemplate(t testing.TestingT, options *Options, chartDir string, releaseName string, templateFiles []string, extraHelmArgs ...string) string {
	out, err := RenderTemplateE(t, options, chartDir, releaseName, templateFiles, extraHelmArgs...)
	require.NoError(t, err)
	return out
}

// RenderTemplateE runs `helm template` to render the template given the provided options and returns stdout/stderr from
// the template command. If you pass in templateFiles, this will only render those templates.
func RenderTemplateE(t testing.TestingT, options *Options, chartDir string, releaseName string, templateFiles []string, extraHelmArgs ...string) (string, error) {
	// First, verify the charts dir exists
	absChartDir, err := filepath.Abs(chartDir)
	if err != nil {
		return "", errors.WithStackTrace(err)
	}
	if !files.FileExists(chartDir) {
		return "", errors.WithStackTrace(ChartNotFoundError{chartDir})
	}

	// check chart dependencies
	if options.BuildDependencies {
		if _, err := RunHelmCommandAndGetOutputE(t, options, "dependency", "build", chartDir); err != nil {
			return "", errors.WithStackTrace(err)
		}
	}

	// Now construct the args
	// We first construct the template args
	args := []string{}
	if options.KubectlOptions != nil && options.KubectlOptions.Namespace != "" {
		args = append(args, "--namespace", options.KubectlOptions.Namespace)
	}
	args, err = getValuesArgsE(t, options, args...)
	if err != nil {
		return "", err
	}
	for _, templateFile := range templateFiles {
		// validate this is a valid template file
		absTemplateFile := filepath.Join(absChartDir, templateFile)
		if !strings.HasPrefix(templateFile, "charts") && !files.FileExists(absTemplateFile) {
			return "", errors.WithStackTrace(TemplateFileNotFoundError{Path: templateFile, ChartDir: absChartDir})
		}

		// Note: we only get the abs template file path to check it actually exists, but the `helm template` command
		// expects the relative path from the chart.
		args = append(args, "--show-only", templateFile)
	}
	// deal extraHelmArgs
	args = append(args, extraHelmArgs...)

	// ... and add the name and chart at the end as the command expects
	args = append(args, releaseName, chartDir)

	// Finally, call out to helm template command
	return RunHelmCommandAndGetStdOutE(t, options, "template", args...)
}

// RenderTemplate runs `helm template` to render a *remote* chart  given the provided options and returns stdout/stderr from
// the template command. If you pass in templateFiles, this will only render those templates. This function will fail
// the test if there is an error rendering the template.
func RenderRemoteTemplate(t testing.TestingT, options *Options, chartURL string, releaseName string, templateFiles []string, extraHelmArgs ...string) string {
	out, err := RenderRemoteTemplateE(t, options, chartURL, releaseName, templateFiles, extraHelmArgs...)
	require.NoError(t, err)
	return out
}

// RenderTemplate runs `helm template` to render a *remote* helm chart  given the provided options and returns stdout/stderr from
// the template command. If you pass in templateFiles, this will only render those templates.
func RenderRemoteTemplateE(t testing.TestingT, options *Options, chartURL string, releaseName string, templateFiles []string, extraHelmArgs ...string) (string, error) {
	// Now construct the args
	// We first construct the template args
	args := []string{}
	if options.KubectlOptions != nil && options.KubectlOptions.Namespace != "" {
		args = append(args, "--namespace", options.KubectlOptions.Namespace)
	}
	args, err := getValuesArgsE(t, options, args...)
	if err != nil {
		return "", err
	}
	for _, templateFile := range templateFiles {
		// As the helm command fails if a non valid template is given as input
		// we do not check if the template file exists or not as we do for local charts
		// as it would add unecessary networking calls
		args = append(args, "--show-only", templateFile)
	}
	// deal extraHelmArgs
	args = append(args, extraHelmArgs...)

	// ... and add the helm chart name, the remote repo and chart URL at the end
	args = append(args, releaseName, "--repo", chartURL)

	// Finally, call out to helm template command
	return RunHelmCommandAndGetStdOutE(t, options, "template", args...)
}

// UnmarshalK8SYaml is the same as UnmarshalK8SYamlE, but will fail the test if there is an error.
func UnmarshalK8SYaml(t testing.TestingT, yamlData string, destinationObj interface{}) {
	require.NoError(t, UnmarshalK8SYamlE(t, yamlData, destinationObj))
}

// UnmarshalK8SYamlE can be used to take template outputs and unmarshal them into the corresponding client-go struct. For
// example, suppose you render the template into a Deployment object. You can unmarshal the yaml as follows:
//
// var deployment appsv1.Deployment
// UnmarshalK8SYamlE(t, renderedOutput, &deployment)
//
// At the end of this, the deployment variable will be populated.
func UnmarshalK8SYamlE(t testing.TestingT, yamlData string, destinationObj interface{}) error {
	// NOTE: the client-go library can only decode json, so we will first convert the yaml to json before unmarshaling
	jsonData, err := yaml.YAMLToJSON([]byte(yamlData))
	if err != nil {
		return errors.WithStackTrace(err)
	}
	err = json.Unmarshal(jsonData, destinationObj)
	if err != nil {
		return errors.WithStackTrace(err)
	}
	return nil
}

// Create/update the manifest snapshot of a chart (e.g bitnami/nginx)
func UpdateSnapshot(yamlData string, releaseName string) {

	snapshotDir := "__snapshot__"
	// Create a directory if not exists
	if !files.FileExists(snapshotDir) {
		os.Mkdir(snapshotDir, 0755)
	}

	filename := snapshotDir + "/" + releaseName + ".yaml"
	// Open a file in write mode
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write the string representation of the "deployment" variable to the file
	_, err = file.WriteString(yamlData)
	if err != nil {
		fmt.Println("Error writing to file: ", filename, err)
		return
	}

	fmt.Println("Content written to file successfully.", filename)
}

// Create/update the manifest snapshot of a chart (e.g bitnami/nginx)
func DiffAgainstSnapshot(yamlData string, releaseName string) int {

	snapshotDir := "__snapshot__"

	filename := snapshotDir + "/" + releaseName + ".yaml"
	from, err := ytbx.LoadFile(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 1
	}

	filename2 := releaseName + ".yaml"
	// Open a file in write mode
	file, err := os.Create(filename2)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return 1
	}

	// Write the string representation of the "deployment" variable to the file
	_, err = file.WriteString(yamlData)
	if err != nil {
		fmt.Println("Error writing to file: ", filename2, err)
		return 1
	}

	defer file.Close()

	to, err := ytbx.LoadFile(filename2)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 1
	}

	compOpt := dyff.KubernetesEntityDetection(false)

	Report, err := dyff.CompareInputFiles(from, to, compOpt)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 1
	}

	reportWriter := &dyff.HumanReport{
		Report:            Report,
		DoNotInspectCerts: false,
		NoTableStyle:      false,
		OmitHeader:        false,
		UseGoPatchPaths:   false,
	}

	number_of_diffs := len(reportWriter.Diffs)

	writer := os.Stdout
	reportWriter.WriteReport(writer)

	// if different, print diff
	// if same, print "no diff"

	defer file.Close()
	return number_of_diffs
}
