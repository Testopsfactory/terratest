package shell

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommandAndGetOutput(t *testing.T) {
	t.Parallel()

	text := "Hello, World"
	cmd := Command{
		Command: "echo",
		Args:    []string{text},
	}

	out := RunCommandAndGetOutput(t, cmd)
	assert.Equal(t, text, strings.TrimSpace(out))
}

func TestRunCommandAndGetOutputOrder(t *testing.T) {
	t.Parallel()

	stderrText := "Hello, Error"
	stdoutText := "Hello, World"
	expectedText := "Hello, Error\nHello, World"
	pythonCode := fmt.Sprintf(
		"from __future__ import print_function; import sys; print('%s', file=sys.stderr); sys.stderr.flush(); print('%s', file=sys.stdout); sys.stdout.flush()",
		stderrText,
		stdoutText,
	)
	cmd := Command{
		Command: "python",
		Args:    []string{"-c", pythonCode},
	}

	out := RunCommandAndGetOutput(t, cmd)
	assert.Equal(t, strings.TrimSpace(out), expectedText)
}
