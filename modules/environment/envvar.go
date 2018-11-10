package environment

import (
	"os"
	"testing"
)

// GetFirstNonEmptyEnvVarOrFatal returns the first non-empty environment variable from ks, or throws a fatal
func GetFirstNonEmptyEnvVarOrFatal(t *testing.T, ks []string) string {
	v := GetFirstNonEmptyEnvVarOrEmptyString(t, ks)
	if v == "" {
		t.Fatalf("All of the following env vars %v are empty. At least one must be non-empty.", ks)
	}

	return v
}

// GetFirstNonEmptyEnvVarOrEmptyString returns the first non-empty environment variable from ks, or returns the empty
// string
func GetFirstNonEmptyEnvVarOrEmptyString(t *testing.T, ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}

	return ""
}
