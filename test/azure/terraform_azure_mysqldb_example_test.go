// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureMySQLDBExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := strings.ToLower(random.UniqueId())
	expectedServerSkuName := "GP_Gen5_2"
	expectedServerStoragemMb := "5120"
	expectedDatabaseCharSet := "utf8"
	expectedDatabaseCollation := "utf8_unicode_ci"

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-mysqldb-example",
		Vars: map[string]interface{}{
			"postfix":                uniquePostfix,
			"mysqlserver_sku_name":   expectedServerSkuName,
			"mysqlserver_storage_mb": expectedServerStoragemMb,
			"mysqldb_charset":        expectedDatabaseCharSet,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	expectedResourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedMYSQLServerName := terraform.Output(t, terraformOptions, "mysql_server_name")

	expectedMYSQLDBName := terraform.Output(t, terraformOptions, "mysql_database_name")

	// website::tag::4:: Get mySQL server details and assert them against the terraform output
	actualMYSQLServerSkuName := azure.GetMYSQLServerSkuName(t, expectedResourceGroupName, expectedMYSQLServerName, "")
	actualServerStoragemMb := azure.GetMYSQLServerStorageMB(t, expectedResourceGroupName, expectedMYSQLServerName, "")

	actualServerState := azure.GetMYSQLServerState(t, expectedResourceGroupName, expectedMYSQLServerName, "")

	assert.Equal(t, expectedServerSkuName, actualMYSQLServerSkuName)
	assert.Equal(t, expectedServerStoragemMb, fmt.Sprint(actualServerStoragemMb))

	assert.Equal(t, mysql.ServerStateReady, actualServerState)

	// website::tag::5:: Get  mySQL server DB details and assert them against the terraform output
	actualDatabaseCharSet := azure.GetMYSQLDBCharset(t, expectedResourceGroupName, expectedMYSQLServerName, expectedMYSQLDBName, "")
	actualDatabaseCollation := azure.GetMYSQLDBCollation(t, expectedResourceGroupName, expectedMYSQLServerName, expectedMYSQLDBName, "")

	assert.Equal(t, expectedDatabaseCharSet, actualDatabaseCharSet)
	assert.Equal(t, expectedDatabaseCollation, actualDatabaseCollation)
}
