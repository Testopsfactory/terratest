package helm

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/gruntwork-cli/collections"
	"github.com/gruntwork-io/gruntwork-cli/errors"
	"github.com/stretchr/testify/require"
)

// FormatSetValuesAsArgs formats the given values as command line args for helm using set (e.g of the format
// --set key=value)
func FormatSetValuesAsArgs(setValues map[string]string) []string {
	args := []string{}

	// To make it easier to test, go through the keys in sorted order
	keys := collections.Keys(setValues)
	for _, key := range keys {
		value := setValues[key]
		argValue := fmt.Sprintf("%s=%s", key, value)
		args = append(args, "--set", argValue)
	}

	return args
}

// FormatValuesFilesAsArgs formats the given list of values file paths as command line args for helm (e.g of the format
// -f path). This will fail the test if one of the paths do not exist.
func FormatValuesFilesAsArgs(t *testing.T, valuesFiles []string) []string {
	args, err := FormatValuesFilesAsArgsE(t, valuesFiles)
	require.NoError(t, err)
	return args
}

// FormatValuesFilesAsArgsE formats the given list of values file paths as command line args for helm (e.g of the format
// -f path)
func FormatValuesFilesAsArgsE(t *testing.T, valuesFiles []string) ([]string, error) {
	args := []string{}

	for _, valuesFilePath := range valuesFiles {
		// TODO: return a typed error
		// Pass through filepath.Abs to make sure this file exists
		absValuesFilePath, err := filepath.Abs(valuesFilePath)
		if err != nil {
			return args, errors.WithStackTrace(err)
		}
		args = append(args, "-f", absValuesFilePath)
	}

	return args, nil
}

// FormatSetFilesAsArgs formats the given list of keys and file paths as command line args for helm to set from file
// (e.g of the format --set-file key=path). This will fail the test if one of the paths do not exist.
func FormatSetFilesAsArgs(t *testing.T, setFiles map[string]string) []string {
	args, err := FormatSetFilesAsArgsE(t, setFiles)
	require.NoError(t, err)
	return args
}

// FormatSetFilesAsArgsE formats the given list of keys and file paths as command line args for helm to set from file
// (e.g of the format --set-file key=path)
func FormatSetFilesAsArgsE(t *testing.T, setFiles map[string]string) ([]string, error) {
	args := []string{}

	// To make it easier to test, go through the keys in sorted order
	keys := collections.Keys(setFiles)
	for _, key := range keys {
		setFilePath := setFiles[key]
		// TODO: return a typed error
		// Pass through filepath.Abs to make sure this file exists
		absSetFilePath, err := filepath.Abs(setFilePath)
		if err != nil {
			return args, errors.WithStackTrace(err)
		}
		argValue := fmt.Sprintf("%s=%s", key, absSetFilePath)
		args = append(args, "--set-file", argValue)
	}

	return args, nil
}
