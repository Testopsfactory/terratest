package http_helper

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func TestRunDummyServer(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueId()
	text := fmt.Sprintf("dummy-server-%s", uniqueID)

	listener, port := RunDummyServer(t, text)
	defer shutDownServer(t, listener)

	url := fmt.Sprintf("http://localhost:%d", port)
	HttpGetWithValidation(t, url, 200, text)
}

func TestContinuouslyCheck(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueId()
	text := fmt.Sprintf("dummy-server-%s", uniqueID)
	stopChecking := make(chan bool, 1)

	listener, port := RunDummyServer(t, text)

	url := fmt.Sprintf("http://localhost:%d", port)
	wg := ContinuouslyCheckUrl(t, url, stopChecking, 1*time.Second)
	defer func() {
		stopChecking <- true
		wg.Wait()
		shutDownServer(t, listener)
	}()
	time.Sleep(5 * time.Second)
}

func shutDownServer(t *testing.T, listener io.Closer) {
	err := listener.Close()
	assert.NoError(t, err)
}
