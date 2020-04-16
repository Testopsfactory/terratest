// Package logger contains different methods to log.
package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	gotesting "testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// defaultLogf will be used if the caller uses the function
// Logf, where on the first call to that function, a sane
// default logging function will be set. We keep that in a
// global variable so that we don't need to define it on
// every call to Logf.
var defaultLogf LogFunc

type Loggers []LogFunc
type LogFunc func(t testing.TestingT, format string, args ...interface{})

func With(l ...LogFunc) *Loggers {
	lo := Loggers(l)
	return &lo
}

// Logf logs to all given loggers. If no loggers are given (or it is nil), it will
// use the default logger. This allows for the following usecases:
//   var l *Loggers
//   l.Logf(...)
//   l = With(TestingT)
//   l.Logf(...)
func (l *Loggers) Logf(t testing.TestingT, format string, args ...interface{}) {
	if tt, ok := t.(*gotesting.T); ok {
		tt.Helper()
	}

	// if l is not initialised or no loggers
	// are supplied, use the default logging.
	if l == nil || len(*l) == 0 {
		logDefaultLogf(t, format, args...)
		return
	}

	for _, logf := range *l {
		logf(t, format, args...)
	}
}

// Discard discards all logging.
func Discard(_ testing.TestingT, format string, args ...interface{}) {}

// TestingT can be used to explicitly use Go's testing.T to log.
// It is also used as the default if Go version >= 1.14 (if detected
// correctly). If this is used, but no testing.T is provided, it will
// fallback to Logger.
func TestingT(t testing.TestingT, format string, args ...interface{}) {
	// this should never fail
	tt, ok := t.(*gotesting.T)
	if !ok {
		// fallback
		DoLog(t, 2, os.Stdout, fmt.Sprintf(format, args...))
		return
	}

	tt.Helper()
	tt.Logf(format, args...)
	return
}

// Logger is the conventional logging utility that terratest uses.
// Default up until Go 1.14.
func Logger(t testing.TestingT, format string, args ...interface{}) {
	DoLog(t, 3, os.Stdout, fmt.Sprintf(format, args...))
}

// Logf logs the given format and arguments with the default logging utility. If the Go
// version can be determined and is 1.14 or above, t.Logf will be used (Go 1.14 introduced
// streaming log output). Else, a default logger will be used that adds a timestamp and
// information about what test and file is doing the logging.
// Compared to the Go's builtin testing.T.Logf, this will always print the output instead
// of buffering it and only display it at the very end of the test.
//
// To have more control over logging, use With(...LogFunc) to get a custom logger.
// Builtin alternatives are Discard, TestingT and Logger.
func Logf(t testing.TestingT, format string, args ...interface{}) {
	if tt, ok := t.(*gotesting.T); ok {
		tt.Helper()
	}

	logDefaultLogf(t, format, args...)
}

func logDefaultLogf(t testing.TestingT, format string, args ...interface{}) {
	if defaultLogf == nil {
		// if a gotesting.T is given and the go version is 1.14, use
		// gotesting.T.Logf
		if tt, ok := t.(*gotesting.T); ok && hasStreamingLogf(runtime.Version()) {
			tt.Helper()
			// we should not assign tt.Logf directly as testing.T may change
			// during the execution of the test (consider subtests, for example).
			defaultLogf = TestingT
		} else {
			defaultLogf = Logger
			defaultLogf(t, "streaming logf not supported, falling back to legacy logger")
		}
	}

	defaultLogf(t, format, args...)
}

// hasStreamingLogf returns true if the go runtime version
// is >= Go 1.14, where streaming Logf output has been
// introduced (https://github.com/golang/go/issues/24929)
func hasStreamingLogf(goVersion string) bool {
	noMajor := strings.TrimPrefix(goVersion, "go1.")
	ver, err := strconv.ParseFloat(noMajor, 32)
	if err != nil {
		return false
	}

	return ver >= 14
}

// Log logs the given arguments to stdout, along with a timestamp and information about what test and file is doing the
// logging. This is an alternative to t.Logf that logs to stdout immediately, rather than buffering all log output and
// only displaying it at the very end of the test. See the Logf method for more info.
func Log(t testing.TestingT, args ...interface{}) {
	DoLog(t, 2, os.Stdout, args...)
}

// DoLog logs the given arguments to the given writer, along with a timestamp and information about what test and file is
// doing the logging.
func DoLog(t testing.TestingT, callDepth int, writer io.Writer, args ...interface{}) {
	date := time.Now()
	prefix := fmt.Sprintf("%s %s %s:", t.Name(), date.Format(time.RFC3339), CallerPrefix(callDepth+1))
	allArgs := append([]interface{}{prefix}, args...)
	fmt.Fprintln(writer, allArgs...)
}

// CallerPrefix returns the file and line number information about the methods that called this method, based on the current
// goroutine's stack. The argument callDepth is the number of stack frames to ascend, with 0 identifying the method
// that called CallerPrefix, 1 identifying the method that called that method, and so on.
//
// This code is adapted from testing.go, where it is in a private method called decorate.
func CallerPrefix(callDepth int) string {
	_, file, line, ok := runtime.Caller(callDepth)
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}

	return fmt.Sprintf("%s:%d", file, line)
}
