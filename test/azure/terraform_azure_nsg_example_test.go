// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureNsgExample(t *testing.T) {
	t.Parallel()

	//
	// Setup our variables to be unique per test-run:
	//

	// "resource_group_name"
	expectedResourceGroupName := fmt.Sprintf("terratest-nsg-example-%s", random.UniqueId())

	// "vnet_name"
	// "subnet_name"
	expectedVnetName := fmt.Sprintf("vnet_name_%s", random.UniqueId())
	expectedSubnetName := fmt.Sprintf("subnet_name_%s", random.UniqueId())

	// "vm_nic_name"
	// "vm_nic_ip_config_name"
	expectedNICName := fmt.Sprintf("vm_nic_name_%s", random.UniqueId())
	expectedIPConfigName := fmt.Sprintf("vm_nic_ip_config_name_%s", random.UniqueId())

	// "nsg_name"
	// "nsg_rule_name"
	expectedNSGName := fmt.Sprintf("nsg_name_%s", random.UniqueId())
	expectedNSGRuleName := fmt.Sprintf("nsg_rule_name_%s", random.UniqueId())

	// "vm_name"
	// "hostname"
	// "os_disk_name"
	expectedVMName := fmt.Sprintf("vm_name_%s", random.UniqueId())
	expectedHostName := fmt.Sprintf("hostname_%s", random.UniqueId())
	expectedOSDiskName := fmt.Sprintf("os_disk_name_%s", random.UniqueId())

	// "password"
	// "username"

	// Construct options for TF apply
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-nsg-example",
		Vars: map[string]interface{}{
			"resource_group_name":   expectedResourceGroupName,
			"vnet_name":             expectedVnetName,
			"subnet_name":           expectedSubnetName,
			"vm_nic_name":           expectedNICName,
			"vm_nic_ip_config_name": expectedIPConfigName,
			"nsg_name":              expectedNSGName,
			"nsg_rule_name":         expectedNSGRuleName,
			"vm_name":               expectedVMName,
			"hostname":              expectedHostName,
			"os_disk_name":          expectedOSDiskName,
		},
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
	assert.True(t, sshRule.AllowsDestinationPort(t, "22"))

	// But should not allow 80 inbound
	assert.False(t, sshRule.AllowsDestinationPort(t, "80"))

	// SSh is allowed from any port
	assert.True(t, sshRule.AllowsSourcePort(t, "*"))
}
