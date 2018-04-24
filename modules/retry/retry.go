package retry

import (
	"github.com/gruntwork-io/terratest/logger"
	"time"
	"fmt"
	"github.com/gruntwork-io/terratest/parallel"
	"testing"
)

// Run the specified action and wait up to the specified timeout for it to complete. Return the output of the action if
// it completes on time or fail the test otherwise.
func DoWithTimeout(t *testing.T, actionDescription string, timeout time.Duration, action func() (string, error)) string {
	out, err := DoWithTimeoutE(t, actionDescription, timeout, action)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Run the specified action and wait up to the specified timeout for it to complete. Return the output of the action if
// it completes on time or an error otherwise.
func DoWithTimeoutE(t *testing.T, actionDescription string, timeout time.Duration, action func() (string, error)) (string, error) {
	resultChannel := make(chan parallel.TestResult, 1)

	go func() {
		out, err := action()
		resultChannel <- parallel.TestResult{Description: actionDescription, Value: out, Err: err}
	}()

	select {
	case result := <-resultChannel:
		return result.Value, result.Err
	case <-time.After(timeout):
		return "", TimeoutExceeded{Description: actionDescription, Timeout: timeout}
	}
}

// Run the specified action. If it returns a value, return that value. If it returns a FatalError, return that error
// immediately. If it returns any other type of error, sleep for sleepBetweenRetries and try again, up to a maximum of
// maxRetries retries. If maxRetries is exceeded, fail the test.
func DoWithRetry(t *testing.T, actionDescription string, maxRetries int, sleepBetweenRetries time.Duration, action func() (string, error)) string {
	out, err := DoWithRetryE(t, actionDescription, maxRetries, sleepBetweenRetries, action)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Run the specified action. If it returns a value, return that value. If it returns a FatalError, return that error
// immediately. If it returns any other type of error, sleep for sleepBetweenRetries and try again, up to a maximum of
// maxRetries retries. If maxRetries is exceeded, return a MaxRetriesExceeded error.
func DoWithRetryE(t *testing.T, actionDescription string, maxRetries int, sleepBetweenRetries time.Duration, action func() (string, error)) (string, error) {
	for i := 0; i < maxRetries; i++ {
		logger.Log(t, actionDescription)

		output, err := action()
		if err == nil {
			return output, nil
		}

		if _, isFatalErr := err.(FatalError); isFatalErr {
			return "", err
		}

		logger.Logf(t, "%s returned an error: %s. Sleeping for %s and will try again.", actionDescription, err.Error(), sleepBetweenRetries)
		time.Sleep(sleepBetweenRetries)
	}

	return "", MaxRetriesExceeded{Description: actionDescription, MaxRetries: maxRetries}
}

// Custom error types

type TimeoutExceeded struct {
	Description string
	Timeout     time.Duration
}

func (err TimeoutExceeded) Error() string {
	return fmt.Sprintf("'%s' did not complete before timeout of %s", err.Description, err.Timeout)
}

type MaxRetriesExceeded struct {
	Description string
	MaxRetries  int
}

func (err MaxRetriesExceeded) Error() string {
	return fmt.Sprintf("'%s' unsuccessful after %d retries", err.Description, err.MaxRetries)
}

// Marker interface for errors that should not be retried
type FatalError error
