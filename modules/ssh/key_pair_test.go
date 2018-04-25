package ssh

import "testing"

// Basic test to ensure we can successfully generate key pairs (no explicit validation for now)
func TestGenerateRSAKeyPair(t *testing.T) {
	t.Parallel()

	GenerateRSAKeyPair(t, 2048)
}


