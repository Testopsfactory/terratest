package ssh

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSshAgentWithKeyPair(t *testing.T) {
	keyPair := GenerateRSAKeyPair(t, 2048)
	sshAgent := SshAgentWithKeyPair(t, keyPair)

	// ensure that socket directory is set in environment, and it exists
	sockFile := filepath.Join(sshAgent.socketDir, "ssh_auth.sock")

	sockFileEnv, found := os.LookupEnv("SSH_AUTH_SOCK")
	assert.FileExists(t, sockFile)
	assert.True(t, found)
	assert.Equal(t, sockFileEnv, sockFile)

	// assert that there's 1 key in the agent
	keys, err := sshAgent.agent.List()
	assert.NoError(t, err)
	assert.Len(t, keys, 1)

	sshAgent.Stop()

	// is socketDir removed as expected?
	if _, err := os.Stat(sshAgent.socketDir); !os.IsNotExist(err) {
		assert.FailNow(t, "ssh agent failed to remove socketDir on Stop()")
	}
}

func TestSshAgentWithKeyPairs(t *testing.T) {
	keyPair := GenerateRSAKeyPair(t, 2048)
	keyPair2 := GenerateRSAKeyPair(t, 2048)
	sshAgent := SshAgentWithKeyPairs(t, []*KeyPair{keyPair, keyPair2})
	defer sshAgent.Stop()

	keys, _ := sshAgent.agent.List()
	assert.Len(t, keys, 2)
}