package dns_helper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var publicDomainNameservers = []string{
	"ns-1499.awsdns-59.org",
	"ns-190.awsdns-23.com",
	"ns-1989.awsdns-56.co.uk",
	"ns-853.awsdns-42.net",
}

var testDNSDatabase = dnsDatabase{
	DNSQuery{"A", "a." + testDomain}: DNSAnswers{
		{"A", "2.2.2.2"},
		{"A", "1.1.1.1"},
	},

	DNSQuery{"AAAA", "aaaa." + testDomain}: DNSAnswers{
		{"AAAA", "2001:db8::aaaa"},
	},

	DNSQuery{"CNAME", "terratest." + testDomain}: DNSAnswers{
		{"CNAME", "gruntwork-io.github.io."},
	},

	DNSQuery{"CNAME", "cname1." + testDomain}: DNSAnswers{
		{"CNAME", "cname2." + testDomain + "."},
	},

	DNSQuery{"A", "cname1." + testDomain}: DNSAnswers{
		{"CNAME", "cname2." + testDomain + "."},
		{"CNAME", "cname3." + testDomain + "."},
		{"CNAME", "cname4." + testDomain + "."},
		{"CNAME", "cnamefinal." + testDomain + "."},
		{"A", "1.1.1.1"},
	},

	DNSQuery{"TXT", "txt." + testDomain}: DNSAnswers{
		{"TXT", `"This is a text."`},
	},

	DNSQuery{"MX", testDomain}: DNSAnswers{
		{"MX", "10 mail." + testDomain + "."},
	},
}

func TestOkDNSFindNameservers(t *testing.T) {
	t.Parallel()
	fqdn := "terratest.gruntwork.io"
	expectedNameservers := publicDomainNameservers
	nameservers, err := DNSFindNameserversE(t, fqdn, nil)
	require.NoError(t, err)
	require.ElementsMatch(t, nameservers, expectedNameservers)
}

func TestErrorDNSFindNameservers(t *testing.T) {
	t.Parallel()
	fqdn := "this.domain.doesnt.exist"
	nameservers, err := DNSFindNameserversE(t, fqdn, nil)
	require.Error(t, err)
	require.Nil(t, nameservers)
}

func TestOkTerratestDNSLookupAuthoritative(t *testing.T) {
	t.Parallel()
	dnsQuery := DNSQuery{"CNAME", "terratest." + testDomain}
	expected := DNSAnswers{{"CNAME", "gruntwork-io.github.io."}}
	res, err := DNSLookupAuthoritativeE(t, dnsQuery, nil)
	require.NoError(t, err)
	require.ElementsMatch(t, res, expected)
}

func TestOkLocalDNSLookupAuthoritative(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServers()
	defer shutDownServers(t, s1, s2)
	for dnsQuery, expected := range testDNSDatabase {
		s1.AddEntryToDNSDatabase(dnsQuery, expected)
		s2.AddEntryToDNSDatabase(dnsQuery, expected)
		res, err := DNSLookupAuthoritativeE(t, dnsQuery, []string{s1.Address(), s2.Address()})
		require.NoError(t, err)
		require.ElementsMatch(t, res, expected)
	}
}

func TestErrorLocalDNSLookupAuthoritative(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServers()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "txt." + testDomain}
	_, err := DNSLookupAuthoritativeE(t, dnsQuery, []string{s1.Address(), s2.Address()})
	if _, ok := err.(*NotFoundError); !ok {
		t.Errorf("unexpected error, got %q", err)
	}
}

func TestOkLocalDNSLookupAuthoritativeAll(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServers()
	defer shutDownServers(t, s1, s2)
	for dnsQuery, expected := range testDNSDatabase {
		s1.AddEntryToDNSDatabase(dnsQuery, expected)
		s2.AddEntryToDNSDatabase(dnsQuery, expected)
		res, err := DNSLookupAuthoritativeE(t, dnsQuery, []string{s1.Address(), s2.Address()})
		require.NoError(t, err)
		require.ElementsMatch(t, res, expected)
	}
}

func TestError1DNSLookupAuthoritativeAll(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServers()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "txt." + testDomain}
	_, err := DNSLookupAuthoritativeAllE(t, dnsQuery, []string{s1.Address(), s2.Address()})
	if _, ok := err.(*NotFoundError); !ok {
		t.Errorf("unexpected error, got %q", err)
	}
}

func TestError2DNSLookupAuthoritativeAll(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServers()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "a." + testDomain}
	s1.AddEntryToDNSDatabase(dnsQuery, DNSAnswers{{"A", "1.1.1.1"}})
	_, err := DNSLookupAuthoritativeAllE(t, dnsQuery, []string{s1.Address(), s2.Address()})
	if _, ok := err.(*NotFoundError); !ok {
		t.Errorf("unexpected error, got %q", err)
	}
}

func TestError3DNSLookupAuthoritativeAll(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServers()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "a." + testDomain}
	s1.AddEntryToDNSDatabase(dnsQuery, DNSAnswers{{"A", "1.1.1.1"}})
	s2.AddEntryToDNSDatabase(dnsQuery, DNSAnswers{{"A", "2.2.2.2"}})
	_, err := DNSLookupAuthoritativeAllE(t, dnsQuery, []string{s1.Address(), s2.Address()})
	if _, ok := err.(*InconsistentAuthoritativeError); !ok {
		t.Errorf("unexpected error, got %q", err)
	}
}

func TestError4DNSLookupAuthoritativeAll(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServers()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "this.domain.doesnt.exist"}
	_, err := DNSLookupAuthoritativeAllE(t, dnsQuery, []string{s1.Address(), s2.Address()})
	if _, ok := err.(*NSNotFoundError); !ok {
		t.Errorf("unexpected error, got %q", err)
	}
}

// Retry until any authoritative nameserver gives an answer
func TestOkDNSLookupAuthoritativeWithRetry(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServersRetry()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "a." + testDomain}
	expectedRes := DNSAnswers{{"A", "1.1.1.1"}, {"A", "2.2.2.2"}}
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	res, err := DNSLookupAuthoritativeWithRetryE(t, dnsQuery, []string{s1.Address(), s2.Address()}, 5, time.Second)
	require.NoError(t, err)
	require.ElementsMatch(t, res, expectedRes)
}

// Retry will fail as the record will never exist
func TestErrorDNSLookupAuthoritativeWithRetry(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServersRetry()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "txt." + testDomain}
	_, err := DNSLookupAuthoritativeWithRetryE(t, dnsQuery, []string{s1.Address(), s2.Address()}, 5, time.Second)
	require.Error(t, err)
	if _, ok := err.(*MaxRetriesExceeded); !ok {
		t.Errorf("unexpected error, got %q", err)
	}
}

// Retry until all authoritative nameservers give the same answers
func TestOkDNSLookupAuthoritativeAllWithRetryNotfound(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServersRetry()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "a." + testDomain}
	expectedRes := DNSAnswers{{"A", "1.1.1.1"}, {"A", "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	res, err := DNSLookupAuthoritativeAllWithRetryE(t, dnsQuery, []string{s1.Address(), s2.Address()}, 5, time.Second)
	require.NoError(t, err)
	require.ElementsMatch(t, res, expectedRes)
}

// Retry until all authoritative nameservers give the same answers
func TestOkDNSLookupAuthoritativeAllWithRetryInconsistent(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServersRetry()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "a." + testDomain}
	expectedRes := DNSAnswers{{"A", "1.1.1.1"}, {"A", "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabase(dnsQuery, DNSAnswers{{"A", "2.2.2.2"}})
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	res, err := DNSLookupAuthoritativeAllWithRetryE(t, dnsQuery, []string{s1.Address(), s2.Address()}, 5, time.Second)
	require.NoError(t, err)
	require.ElementsMatch(t, res, expectedRes)
}

// Retry will fail as one authoritative nameserver will always give an extra answer
func TestErrorDNSLookupAuthoritativeAllWithRetry(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServersRetry()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "a." + testDomain}
	s1.AddEntryToDNSDatabase(dnsQuery, DNSAnswers{{"A", "2.2.2.2"}})
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, DNSAnswers{{"A", "1.1.1.1"}, {"A", "2.2.2.2"}})
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, DNSAnswers{{"A", "1.1.1.1"}})
	_, err := DNSLookupAuthoritativeAllWithRetryE(t, dnsQuery, []string{s1.Address(), s2.Address()}, 5, time.Second)
	require.Error(t, err)
	if _, ok := err.(*MaxRetriesExceeded); !ok {
		t.Errorf("unexpected error, got %q", err)
	}
}

// Validate all authoritative nameservers give the expected answers
func TestOkDNSLookupAuthoritativeAllWithValidation(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServers()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "a." + testDomain}
	expectedRes := DNSAnswers{{"A", "1.1.1.1"}, {"A", "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	err := DNSLookupAuthoritativeAllWithValidationE(t, dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes)
	require.NoError(t, err)
}

// Retry until all authoritative nameservers give the expected answers
func TestOkDNSLookupAuthoritativeAllWithValidationRetry(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServersRetry()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "a." + testDomain}
	expectedRes := DNSAnswers{{"A", "1.1.1.1"}, {"A", "2.2.2.2"}}
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	err := DNSLookupAuthoritativeAllWithValidationRetryE(t, dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes, 5, time.Second)
	require.NoError(t, err)
}

// Retry until all authoritative nameservers give the expected answers
func TestOk2DNSLookupAuthoritativeAllWithValidationRetry(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServersRetry()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "a." + testDomain}
	expectedRes := DNSAnswers{{"A", "1.1.1.1"}, {"A", "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, DNSAnswers{{"A", "2.2.2.2"}})
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	err := DNSLookupAuthoritativeAllWithValidationRetryE(t, dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes, 5, time.Second)
	require.NoError(t, err)
}

// Retry until all authoritative nameservers give the expected answers
func TestOk3DNSLookupAuthoritativeAllWithValidationRetry(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServersRetry()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "a." + testDomain}
	expectedRes := DNSAnswers{{"A", "1.1.1.1"}, {"A", "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabase(dnsQuery, DNSAnswers{{"A", "2.2.2.2"}})
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	err := DNSLookupAuthoritativeAllWithValidationRetryE(t, dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes, 5, time.Second)
	require.NoError(t, err)
}

// Retry will fail as one authoritative nameserver will never give the expected answers
func TestErrorDNSLookupAuthoritativeAllWithValidationRetry(t *testing.T) {
	t.Parallel()
	s1, s2 := setupTestDNSServersRetry()
	defer shutDownServers(t, s1, s2)
	dnsQuery := DNSQuery{"A", "a." + testDomain}
	expectedRes := DNSAnswers{{"A", "1.1.1.1"}, {"A", "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabase(dnsQuery, DNSAnswers{{"A", "2.2.2.2"}})
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, DNSAnswers{{"A", "2.2.2.2"}})
	err := DNSLookupAuthoritativeAllWithValidationRetryE(t, dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes, 5, time.Second)
	if _, ok := err.(*MaxRetriesExceeded); !ok {
		t.Errorf("unexpected error, got %q", err)
	}
}

func shutDownServers(t *testing.T, s1, s2 *dnsTestServer) {
	err := s1.Server.Shutdown()
	assert.NoError(t, err)
	err = s2.Server.Shutdown()
	assert.NoError(t, err)
}
