package terraform

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

func TestGetVariablesFromVarFiles(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())
	randomFileName2 := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	testHcl2 := []byte(`
		aws_region     = "us-west-2"`)
	
	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	WriteFile(t, randomFileName2, testHcl2)
	defer os.Remove(randomFileName2)

	varMap, err := GetVariablesFromVarFiles(&Options{
		VarFiles: []string{randomFileName, randomFileName2},
	})

	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}

	require.Equal(t, 2, len(varMap))
	require.Equal(t, "us-east-2", varMap[0]["aws_region"])
	require.Equal(t, "111111111111", varMap[0]["aws_account_id"])
	require.Equal(t, map[string]interface{}{ "foo": "bar", }, varMap[0]["tags"].([]map[string]interface{})[0])
	require.Equal(t, []interface{}{ "item1" }, varMap[0]["list"].([]interface{}))
	require.Equal(t, "us-west-2", varMap[1]["aws_region"])
}

func TestGetVariablesFromVarFilesNoFileError(t *testing.T) {
	_, err := GetVariablesFromVarFiles(&Options{
		VarFiles: []string{"thisdoesntexist"},
	})

	require.Equal(t, "open thisdoesntexist: no such file or directory", err.Error())
}

func TestGetVariablesFromVarFilesBadFile(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())
	testHcl := []byte(`
		thiswillnotwork`)
	
	err := ioutil.WriteFile(randomFileName, testHcl, 0644)

	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}

	defer os.Remove(randomFileName)

	_, err = GetVariablesFromVarFiles(&Options{
		VarFiles: []string{randomFileName},
	})

	if err == nil {
		t.FailNow()
	}

	// HCL library could change their error string, so we are only testing the error string contains what we add to it
	require.Regexp(t, fmt.Sprintf("^%s - ", randomFileName), err.Error())

}

// Helper function to write a file to the filesystem
// Will immediately fail the test if it could not write the file
func WriteFile(t *testing.T, fileName string, bytes []byte) {
	err := ioutil.WriteFile(fileName, bytes, 0644)

	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
}