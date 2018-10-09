package terraform

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
)

func TestOutputList(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-list", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &Options{
		TerraformDir: testFolder,
	}

	InitAndApply(t, options)
	out := OutputList(t, options, "giant_steps")

	expectedLen := 4
	expectedItem := "John Coltrane"

	assert.Len(t, out, expectedLen, "Output should contain %d items", expectedLen)
	assert.Contains(t, out, expectedItem, "Output should contain %q item", expectedItem)
	assert.Equal(t, out[0], expectedItem, "First item should be %q, got %q", expectedItem, out[0])
}

func TestOutputNotList(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-list", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &Options{
		TerraformDir: testFolder,
	}

	InitAndApply(t, options)
	out := OutputList(t, options, "not_a_list")

	expectedLen := 1
	expectedItem := "This is not a list."

	assert.Len(t, out, expectedLen, "Output should contain %d items", expectedLen)
	assert.Equal(t, out[0], expectedItem, "First item should be %q, got %q", expectedItem, out[0])
}
