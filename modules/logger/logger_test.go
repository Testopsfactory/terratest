package logger

import (
	"testing"
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"fmt"
)

func TestDoLog(t *testing.T) {
	t.Parallel()

	text := "test-do-log"
	var buffer bytes.Buffer

	DoLog(t, 1, &buffer, text)

	assert.Regexp(t, fmt.Sprintf("^%s .+? [[:word:]]+.go:[0-9]+: %s$", t.Name(), text), strings.TrimSpace(buffer.String()))
}
