// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureNsgExample(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-nsg-example",
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	nsgName := terraform.Output(t, terraformOptions, "nsg_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// A default NSG has 6 rules, and we have one custom rule for a total of 7
	rules, err := azure.GetAllNSGRulesE(resourceGroupName, nsgName, "")
	assert.NoError(t, err)
	assert.Equal(t, 7, len(rules.SummarizedRules))

	// We should have a rule named "allowSSH"
	sshRule := rules.FindRuleByName("allowSSH")

	// That rule should allow port 22 inbound
	assert.True(t, sshRule.AllowsDestinationPort("22"))

	// But should not allow 80 inbound
	assert.False(t, sshRule.AllowsDestinationPort("80"))

	// SSh is allowed from any port
	assert.True(t, sshRule.AllowsSourcePort("*"))
}
