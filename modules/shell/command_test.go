package shell

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/random"
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
		"from __future__ import print_function; import sys, time; print('%s', file=sys.stderr); sys.stderr.flush(); time.sleep(1); print('%s', file=sys.stdout); sys.stdout.flush()",
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

func TestRunCommandAndGetOutputConcurrency(t *testing.T) {
	t.Parallel()

	uniqueStderr := random.UniqueId()
	uniqueStdout := random.UniqueId()

	bashCode := fmt.Sprintf(`
echo_stderr(){
	sleep .$[ ( $RANDOM %% 10 ) + 1 ]s
	(>&2 echo "%s")
}
echo_stdout(){
	sleep .$[ ( $RANDOM %% 10 ) + 1 ]s
	echo "%s"
}
for i in {1..5}
do
	echo_stderr &
	echo_stdout &
done
wait
`,
		uniqueStderr,
		uniqueStdout,
	)
	cmd := Command{
		Command: "bash",
		Args:    []string{"-c", bashCode},
	}

	out := RunCommandAndGetOutput(t, cmd)
	stdoutReg := regexp.MustCompile(uniqueStdout)
	stderrReg := regexp.MustCompile(uniqueStderr)
	assert.Equal(t, len(stdoutReg.FindAllString(out, -1)), 5)
	assert.Equal(t, len(stderrReg.FindAllString(out, -1)), 5)
}
