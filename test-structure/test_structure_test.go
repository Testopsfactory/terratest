package test_structure

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	terralog "github.com/gruntwork-io/terratest/log"
	"github.com/gruntwork-io/terratest/files"
	"github.com/gruntwork-io/terratest"
)

type testData struct {
	Foo string
	Bar bool
	Baz map[string]interface{}
}

func TestSaveAndLoadTestData(t *testing.T) {
	t.Parallel()

	logger := terralog.NewLogger("TestSaveAndLoadTestData")

	tmpFile, err := ioutil.TempFile("", "save-and-load-test-data")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	expectedData := testData{
		Foo: "foo",
		Bar: true,
		Baz: map[string]interface{}{"abc": "def", "ghi": 1.0, "klm": false},
	}
	SaveTestData(t, tmpFile.Name(), expectedData, logger)

	actualData := testData{}
	LoadTestData(t, tmpFile.Name(), &actualData, logger)
	assert.Equal(t, expectedData, actualData)

	CleanupTestData(t, tmpFile.Name(), logger)
	assert.False(t, files.FileExists(tmpFile.Name()))
}

func TestSaveAndLoadTerratestOptions(t *testing.T) {
	t.Parallel()

	logger := terralog.NewLogger("TestSaveAndLoadTerratestOptions")

	tmpFolder, err := ioutil.TempDir("", "save-and-load-terratest-options")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	expectedData := &terratest.TerratestOptions{
		UniqueId: "foo",
		TestName: "bar",
		TemplatePath: "/abc/def/ghi",
		Vars: map[string]interface{}{},
	}
	SaveTerratestOptions(t, tmpFolder, expectedData, logger)

	actualData := LoadTerratestOptions(t, tmpFolder, logger)
	assert.Equal(t, expectedData, actualData)

	CleanupTerratestOptions(t, tmpFolder, logger)
	assert.False(t, files.FileExists(formatTerratestOptionsPath(tmpFolder)))
}
