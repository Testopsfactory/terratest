package http_helper

import (
	"time"
	"net/http"
	"io/ioutil"
	"strings"
	"fmt"
	"testing"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
)

// Perform an HTTP GET on the given URL and return the HTTP status code and body. If there's any error, fail the test.
func HttpGet(t *testing.T, url string) (int, string) {
	statusCode, body, err := HttpGetE(t, url)
	if err != nil {
		t.Fatal(err)
	}
	return statusCode, body
}

// Perform an HTTP GET on the given URL and return the HTTP status code, body, and any error.
func HttpGetE(t *testing.T, url string) (int, string, error) {
	logger.Logf(t, "Making an HTTP GET call to URL", url)

	client := http.Client{
		// By default, Go does not impose a timeout, so an HTTP connection attempt can hang for a LONG time.
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return -1, "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, strings.TrimSpace(string(body)), nil
}

// Repeatedly perform an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetry(t *testing.T, url string, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) {
	err := HttpGetWithRetryE(t, url, expectedStatus, expectedBody, retries, sleepBetweenRetries)
	if err != nil {
		t.Fatal(err)
	}
}

// Repeatedly perform an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetryE(t *testing.T, url string, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", url), retries, sleepBetweenRetries, func() (string, error) {
		status, body, err := HttpGetE(t, url)

		if err != nil {
			return "", err
		} else if status != expectedStatus || body != expectedBody {
			return "", UnexpectedHttpResponse{Url: url, ExpectedStatus: expectedStatus, ActualStatus: status, ExpectedBody: expectedBody, ActualBody: body}
		} else {
			logger.Logf(t, "Got expected status code %d from URL %s and expected body:\n%s\n", expectedStatus, url, body)
			return body, nil
		}
	})

	return err
}

// Repeatedly perform an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidation(t *testing.T, url string, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) {
	err := HttpGetWithRetryWithCustomValidationE(t, url, retries, sleepBetweenRetries, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// Repeatedly perform an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidationE(t *testing.T, url string, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", url), retries, sleepBetweenRetries, func() (string, error) {
		status, body, err := HttpGetE(t, url)

		if err != nil {
			return "", err
		} else if !validateResponse(status, body) {
			return "", ValidationFunctionFailed{Url: url, Status: status, Body: body}
		} else {
			logger.Logf(t, "Validation function passed for URL %s, with status code %d and body: %s\n", url, status, body)
			return body, nil
		}
	})

	return err
}

type UnexpectedHttpResponse struct {
	Url            string
	ExpectedStatus int
	ActualStatus   int
	ExpectedBody   string
	ActualBody     string
}

func (err UnexpectedHttpResponse) Error() string {
	return fmt.Sprintf("Unexpected HTTP response from URL %s.\nExpected status: %d\nActual status: %d\n\nExpected body:\n%s\n\nActual body\n%s\n", err.Url, err.ExpectedStatus, err.ActualStatus, err.ExpectedBody, err.ActualBody)
}

type ValidationFunctionFailed struct {
	Url    string
	Status int
	Body   string
}

func (err ValidationFunctionFailed) Error() string {
	return fmt.Sprintf("Validation failed for URL %s. Response status: %d. Response body:\n%s", err.Url, err.Status, err.Body)
}
