package terraform

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// Output calls terraform output for the given variable and return its value.
func Output(t *testing.T, options *Options, key string) string {
	out, err := OutputE(t, options, key)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// OutputE calls terraform output for the given variable and return its value.
func OutputE(t *testing.T, options *Options, key string) (string, error) {
	output, err := RunTerraformCommandE(t, options, "output", "-no-color", key)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

// OutputRequired calls terraform output for the given variable and return its value. If the value is empty, fail the test.
func OutputRequired(t *testing.T, options *Options, key string) string {
	out, err := OutputRequiredE(t, options, key)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// OutputRequiredE calls terraform output for the given variable and return its value. If the value is empty, return an error.
func OutputRequiredE(t *testing.T, options *Options, key string) (string, error) {
	out, err := OutputE(t, options, key)

	if err != nil {
		return "", err
	}
	if out == "" {
		return "", EmptyOutput(key)
	}

	return out, nil
}

// OutputList calls terraform output for the given variable and return its value as a list.
// If the output is a single value, then it returns a list with just one item.
func OutputList(t *testing.T, options *Options, key string) []string {
	out, err := OutputListE(t, options, key)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// OutputListE calls terraform output for the given variable and return its value as a list.
// If the output is a single value, then it returns a list with just one item.
func OutputListE(t *testing.T, options *Options, key string) ([]string, error) {
	out, err := RunTerraformCommandE(t, options, "output", "-no-color", "-json", key)
	if err != nil {
		return nil, err
	}

	outputMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(out), &outputMap); err != nil {
		t.Fatalf("Failed to parse JSON for value %s: %v", out, err)
	}

	value := outputMap["value"]
	var list []string
	switch t := value.(type) {
	case []interface{}:
		list = make([]string, len(t))
		for i, item := range t {
			list[i] = item.(string)
		}
	case interface{}:
		list = make([]string, 1)
		list[0] = t.(string)
	}

	return list, nil
}

// EmptyOutput is an error that occurs when an output is empty.
type EmptyOutput string

func (outputName EmptyOutput) Error() string {
	return fmt.Sprintf("Required output %s was empty", string(outputName))
}
