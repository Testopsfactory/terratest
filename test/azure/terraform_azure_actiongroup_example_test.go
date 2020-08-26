// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"strings"

	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureActionGroupExample(t *testing.T) {
	t.Parallel()
	_random := strings.ToLower(random.UniqueId())

	expectedResourceGroupName := fmt.Sprintf("tmp-rg-%s", _random)
	expectedAppName := fmt.Sprintf("tmp-asp-%s", _random)
	expectedLocation := "westus2"
	expectedShortName := "blah"
	expectedEmailName := "emailTestName"
	expectedEmailAddress := "sample@test.com"
	expectedWebhookName := "webhookTestName"
	expectedWebhookServiceUri := "http://example.com/alert"

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/azure/terraform-azure-actiongroup-example",
		Vars: map[string]interface{}{
			"resourceGroupName": expectedResourceGroupName,
			"appName":           expectedAppName,
			"location":          expectedLocation,
			"shortName":         expectedShortName,
			"enableEmail":       true,
			"emailName":         expectedEmailName,
			"emailAddress":      expectedEmailAddress,
			"enableWebHook":     true,
			"webhookName":       expectedWebhookName,
			"webhookServiceUri": expectedWebhookServiceUri,
		},
	}
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	assert := assert.New(t)
	actionGroupId := terraform.Output(t, terraformOptions, "action_group_id")
	assert.NotNil(actionGroupId)
	assert.Contains(actionGroupId, expectedAppName)

	actionGroup := azure.GetActionGroupResource(t, expectedAppName, expectedResourceGroupName, "")

	assert.NotNil(actionGroup)
	assert.Equal(1, len(*actionGroup.EmailReceivers))
	assert.Equal(0, len(*actionGroup.SmsReceivers))
	assert.Equal(1, len(*actionGroup.WebhookReceivers))
}
