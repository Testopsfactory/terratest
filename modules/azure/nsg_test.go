// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPortRangeParsing(t *testing.T) {
	var cases = []struct {
		portRange    string
		expectedLo   int
		expectedHi   int
		expectsError bool
	}{
		{"22", 22, 22, false},
		{"22-80", 22, 80, false},
		{"*", 0, 65535, false},
		{"*-*", 0, 0, true},
		{"22-", 0, 0, true},
		{"-80", 0, 0, true},
		{"-", 0, 0, true},
		{"80-22", 22, 80, false},
	}

	for _, tt := range cases {
		t.Run(tt.portRange, func(t *testing.T) {
			lo, hi, err := parsePortRangeString(tt.portRange)
			if !tt.expectsError {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.expectedLo, int(lo))
			assert.Equal(t, tt.expectedHi, int(hi))
		})
	}
}

func TestAllowSourcePort(t *testing.T) {
	var cases = []struct {
		CaseName        string
		SourcePortRange string
		TestPort        string
		Result          bool
	}{
		{"22 allows 22", "22", "22", true},
		{"22 doesn't allow 80", "22", "80", false},
		{"Any allows any", "*", "*", true},
	}

	for _, tt := range cases {
		t.Run(tt.CaseName, func(t *testing.T) {
			summary := NsgRuleSummary{}
			summary.SourcePortRange = tt.SourcePortRange
			result := summary.AllowsSourcePort(t, tt.TestPort)
			assert.Equal(t, result, tt.Result)
		})
	}
}

func TestAllowDestinationPort(t *testing.T) {
	var cases = []struct {
		CaseName        string
		SourcePortRange string
		TestPort        string
		Result          bool
	}{
		{"22 allows 22", "22", "22", true},
		{"22 doesn't allow 80", "22", "80", false},
		{"Any allows any", "*", "*", true},
	}

	for _, tt := range cases {
		t.Run(tt.CaseName, func(t *testing.T) {
			summary := NsgRuleSummary{}
			summary.DestinationPortRange = tt.SourcePortRange
			result := summary.AllowsDestinationPort(t, tt.TestPort)
			assert.Equal(t, result, tt.Result)
		})
	}
}

func TestFindSummarizedRule(t *testing.T) {
	var cases = []struct {
		SearchString string
		Result       bool
	}{
		{"rule_1", true},
		{"rule_2", true},
		{"rule_3", true},
		{"rule_4", true},
		{"rule_5", true},
		{"rule_6", false},
		{"", false},
		{"foo", false},
	}

	ruleList := NsgRuleSummaryList{}
	rules := make([]NsgRuleSummary, 0)

	// Create some base rules
	for i := 1; i <= 5; i++ {
		rule := NsgRuleSummary{}
		rule.Name = fmt.Sprintf("rule_%d", i)
		rules = append(rules, rule)
	}
	ruleList.SummarizedRules = rules

	for _, tt := range cases {
		t.Run(tt.SearchString, func(t *testing.T) {
			match := ruleList.FindRuleByName(tt.SearchString)
			if tt.Result {
				assert.Equal(t, tt.SearchString, match.Name)
			} else {
				assert.Equal(t, match, NsgRuleSummary{})
			}
		})
	}
}
