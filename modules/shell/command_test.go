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

	stdoutText := "Hello, Error"
	stderrText := "Hello, World"
	expectedText := "Hello, Error\nHello, World"
	pythonCode := fmt.Sprintf(
		"from __future__ import print_function; import sys; print('%s', file=sys.stderr); print('%s', file=sys.stdout)",
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
