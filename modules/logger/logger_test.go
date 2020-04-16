package logger

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	tftesting "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
)

func TestDoLog(t *testing.T) {
	t.Parallel()

	text := "test-do-log"
	var buffer bytes.Buffer

	DoLog(t, 1, &buffer, text)

	assert.Regexp(t, fmt.Sprintf("^%s .+? [[:word:]]+.go:[0-9]+: %s$", t.Name(), text), strings.TrimSpace(buffer.String()))
}

func TestVersion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		version      string
		hasStreaming bool
	}{
		{"go1.14", true},
		{"go1.14.1", true},
		{"go1.13.5", false},
		{"go1.15", true},
		{"ae0365fab0", false}, // simulate some commit ID
	}

	for _, test := range cases {
		assert.Equal(t, test.hasStreaming, hasStreamingLogf(test.version))
	}
}

func TestCustomLogger(t *testing.T) {
	var logs []string
	customLogf := func(t tftesting.TestingT, format string, args ...interface{}) {
		logs = append(logs, fmt.Sprintf(format, args...))
	}

	Logf(t, "this should be logged with legacylogger or testing.T if go >=1.14")
	var l *Loggers

	l.Logf(t, "this should be logged with legacylogger or testing.T if go >=1.14")
	l = With()
	l.Logf(t, "this should be logged with legacylogger or testing.T if go >=1.14")

	// try all loggers, though this may spam the output a bit.
	l = With(customLogf, Logger, TestingT, Discard)

	l.Logf(t, "log output 1")
	l.Logf(t, "log output 2")

	t.Run("logger-subtest", func(t *testing.T) {
		l.Logf(t, "subtest log")
	})

	assert.Len(t, logs, 3)
	assert.Equal(t, "log output 1", logs[0])
	assert.Equal(t, "log output 2", logs[1])
	assert.Equal(t, "subtest log", logs[2])
}
