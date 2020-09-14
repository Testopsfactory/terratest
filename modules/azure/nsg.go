package azure

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
)

// GetDefaultNsgRulesClientE returns a rules client
func GetDefaultNsgRulesClientE(subscriptionID string) (network.DefaultSecurityRulesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return network.DefaultSecurityRulesClient{}, err
	}

	nsgClient := network.NewDefaultSecurityRulesClient(subscriptionID)

	// Get an authorizer
	auth, err := NewAuthorizer()
	if err != nil {
		return network.DefaultSecurityRulesClient{}, err
	}

	nsgClient.Authorizer = *auth
	return nsgClient, nil
}

// GetCustomNsgRulesClientE returns a rules client
func GetCustomNsgRulesClientE(subscriptionID string) (network.SecurityRulesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return network.SecurityRulesClient{}, err
	}

	nsgClient := network.NewSecurityRulesClient(subscriptionID)

	// Get an authorizer
	auth, err := NewAuthorizer()
	if err != nil {
		return network.SecurityRulesClient{}, err
	}

	nsgClient.Authorizer = *auth
	return nsgClient, nil
}

// GetAllNSGRulesE returns a slice containing all rules from the NSG (including defauls and custom)
func GetAllNSGRulesE(resourceGroupName, nsgName, subscriptionID string) (NsgRuleSummaryList, error) {
	defaultRulesClient, err := GetDefaultNsgRulesClientE(subscriptionID)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	customRulesClient, err := GetCustomNsgRulesClientE(subscriptionID)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	defaultRuleList, err := defaultRulesClient.ListComplete(context.Background(), resourceGroupName, nsgName)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	customRuleList, err := customRulesClient.ListComplete(context.Background(), resourceGroupName, nsgName)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	rules := make([]NsgRuleSummary, 0)
	bindRuleList(&rules, defaultRuleList)
	bindRuleList(&rules, customRuleList)
	ruleList := NsgRuleSummaryList{}
	ruleList.SummarizedRules = rules
	return ruleList, nil
}

// NsgRuleSummaryList holds a colleciton of NsgRuleSummary structs
type NsgRuleSummaryList struct {
	SummarizedRules []NsgRuleSummary
}

// NsgRuleSummary is a string-based summary of an NSG rule with methods attached.
type NsgRuleSummary struct {
	Name                     string
	Description              string
	Protocol                 string
	SourcePortRange          string
	DestinationPortRange     string
	SourceAddressPrefix      string
	DestinationAddressPrefix string
	Access                   string
	Priority                 int32
	Direction                string
}

func bindRuleList(dest *[]NsgRuleSummary, source network.SecurityRuleListResultIterator) {
	for source.NotDone() {
		v := source.Value()
		*dest = append(*dest, convertToNsgRuleSummary(v.Name, v.SecurityRulePropertiesFormat))
		source.NextWithContext(context.Background())
	}
}

func convertToNsgRuleSummary(name *string, rule *network.SecurityRulePropertiesFormat) NsgRuleSummary {
	summary := NsgRuleSummary{}

	if rule.Description != nil {
		summary.Description = *rule.Description
	} else {
		summary.Description = ""
	}
	summary.Name = *name
	summary.Protocol = string(rule.Protocol)
	summary.SourcePortRange = *rule.SourcePortRange
	summary.DestinationPortRange = *rule.DestinationPortRange
	summary.SourceAddressPrefix = *rule.SourceAddressPrefix
	summary.DestinationAddressPrefix = *rule.DestinationAddressPrefix
	summary.Access = string(rule.Access)
	summary.Priority = *rule.Priority
	summary.Direction = string(rule.Direction)
	return summary
}

// FindRuleByName looks for a matching rule name
func (summarizedRules *NsgRuleSummaryList) FindRuleByName(t *testing.T, name string) NsgRuleSummary {
	for _, r := range summarizedRules.SummarizedRules {
		if r.Name == name {
			return r
		}
	}

	t.Error("Supplied rule name not found")
	return NsgRuleSummary{}
}

// AllowsDestinationPort checks to see if the rule allows a specific destination port.
func (summarizedRule *NsgRuleSummary) AllowsDestinationPort(port string) bool {
	if summarizedRule.DestinationPortRange == "*" {
		return true
	}

	low, high := parsePortRangeString(summarizedRule.DestinationPortRange)
	portAsInt, _ := strconv.ParseInt(port, 10, 16)

	if (port == "*") && (low == 0) && (high == 65535) {
		return true
	}

	return ((low <= uint16(portAsInt)) && (uint16(portAsInt) <= high))
}

// AllowsSourcePort checks to see if the rule allows a specific source port.
func (summarizedRule *NsgRuleSummary) AllowsSourcePort(port string) bool {
	if summarizedRule.SourcePortRange == "*" {
		return true
	}

	low, high := parsePortRangeString(summarizedRule.SourcePortRange)
	portAsInt, _ := strconv.ParseInt(port, 10, 16)

	if (port == "*") && (low == 0) && (high == 65535) {
		return true
	}

	return ((low <= uint16(portAsInt)) && (uint16(portAsInt) <= high))
}

// parsePortRangeString decodes a range string ("2-100") or a single digit ("22") and returns
// a tuple in [low, hi] form. Note that if a single digit is supplied, both members of the
// return tuple will be the same value (e.g., "22" returns (22, 22))
func parsePortRangeString(rangeString string) (low uint16, hight uint16) {
	// Is this a range?
	if strings.Index(rangeString, "-") == -1 {
		val, _ := strconv.ParseInt(rangeString, 10, 16)
		return uint16(val), uint16(val)
	}

	parts := strings.Split(rangeString, "-")
	lowVal, _ := strconv.ParseInt(parts[0], 10, 16)
	highVal, _ := strconv.ParseInt(parts[1], 10, 16)
	return uint16(lowVal), uint16(highVal)
}
