package http_helper

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
)

// Continuously check the given URL every 1 second until the stopChecking channel receives a signal to stop.
func ContinuouslyCheckUrl(t *testing.T, url string, stopChecking <-chan bool, sleepBetweenChecks time.Duration) {
	go func() {
		for {
			select {
			case <-stopChecking:
				logger.Logf(t, "Got signal to stop downtime checks for URL %s.\n", url)
				return
			case <-time.After(sleepBetweenChecks):
				statusCode, body, err := HttpGetE(t, url)
				logger.Logf(t, "Got response %d and err %v from URL at %s", statusCode, err, url)
				if err != nil {
					t.Fatalf("Failed to make HTTP request to the URL at %s: %s\n", url, err.Error())
				} else if statusCode != 200 {
					t.Fatalf("Got a non-200 response (%d) from the URL at %s, which means there was downtime! Response body: %s", statusCode, url, body)
				}
			}
		}
	}()
}
