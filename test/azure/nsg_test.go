// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
)

func TestRangeParsingSinglePort(t *testing.T) {
	lo, hi := azure.ParsePortRangeString("22")
	assert.Equal(t, uint16(22), lo)
	assert.Equal(t, uint16(22), hi)
}

func TestRangeParsingPortRange(t *testing.T) {
	lo, hi := azure.ParsePortRangeString("22-80")
	assert.Equal(t, uint16(22), lo)
	assert.Equal(t, uint16(80), hi)
}

func TestRangeParsingAsterisk(t *testing.T) {
	lo, hi := azure.ParsePortRangeString("*")
	assert.Equal(t, uint16(0), lo)
	assert.Equal(t, uint16(65535), hi)
}

func TestRuleSummaryAllowSourcePort(t *testing.T) {
	summary := azure.NsgRuleSummary{}
	summary.SourcePortRange = "22"

	result := summary.AllowsSourcePort("22")
	assert.True(t, result)
}

func TestRuleSummaryAllowSourcePortAsterisk(t *testing.T) {
	summary := azure.NsgRuleSummary{}
	summary.SourcePortRange = "*"

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := summary.AllowsSourcePort(string(uint16(rand.Int())))
	assert.True(t, result)
}

func TestRuleSummaryAllowDestinationPort(t *testing.T) {
	summary := azure.NsgRuleSummary{}
	summary.DestinationPortRange = "80"

	result := summary.AllowsDestinationPort("80")
	assert.True(t, result)
}

func TestRuleSummaryAllowDestinationPortAsterisk(t *testing.T) {
	summary := azure.NsgRuleSummary{}
	summary.DestinationPortRange = "*"

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := summary.AllowsDestinationPort(string(uint16(rand.Int())))
	assert.True(t, result)
}

func TestFindSummarizedRule(t *testing.T) {
	ruleList := azure.NsgRuleSummaryList{}
	rules := make([]azure.NsgRuleSummary, 0)

	// Create some rules
	for i := 1; i <= 10; i++ {
		rule := azure.NsgRuleSummary{}
		rule.Name = fmt.Sprintf("rule_%d", i)
		rules = append(rules, rule)
	}
	ruleList.SummarizedRules = rules

	// Look for a rule that exists
	match1 := ruleList.FindRuleByName("rule_5")
	assert.Equal(t, "rule_5", match1.Name)

	// Look for a rule that doesn't exist
	match2 := ruleList.FindRuleByName("foo")
	assert.Equal(t, azure.NsgRuleSummary{}, match2)
}
