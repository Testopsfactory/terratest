package packer

import (
	"fmt"
	"testing"
)

func TestExtractAmiIdFromOneLine(t *testing.T) {
	t.Parallel()

	expectedAMIID := "ami-b481b3de"
	text := fmt.Sprintf("1456332887,amazon-ebs,artifact,0,id,us-east-1:%s", expectedAMIID)
	actualAMIID, err := extractAMIID(text)

	if err != nil {
		t.Errorf("Did not expect to get an error when extracting a valid AMI ID: %s", err)
	}

	if actualAMIID != expectedAMIID {
		t.Errorf("Did not get expected AMI ID. Expected: %s. Actual: %s.", expectedAMIID, actualAMIID)
	}
}

func TestExtractAmiIdFromMultipleLines(t *testing.T) {
	t.Parallel()

	expectedAMIID := "ami-b481b3de"
	text := fmt.Sprintf(`
	foo
	bar
	1456332887,amazon-ebs,artifact,0,id,us-east-1:%s
	baz
	blah
	`, expectedAMIID)

	actualAMIID, err := extractAMIID(text)

	if err != nil {
		t.Errorf("Did not expect to get an error when extracting a valid AMI ID: %s", err)
	}

	if actualAMIID != expectedAMIID {
		t.Errorf("Did not get expected AMI ID. Expected: %s. Actual: %s.", expectedAMIID, actualAMIID)
	}
}

func TestExtractAmiIdNoIdPresent(t *testing.T) {
	t.Parallel()

	text := `
	foo
	bar
	baz
	blah
	`

	_, err := extractAMIID(text)

	if err == nil {
		t.Error("Expected to get an error when extracting an AMI ID from text with no AMI in it, but got nil")
	}

}
