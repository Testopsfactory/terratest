// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. To avoid overloading the system, we run the
// kubernetes tests separately from the others.

package k8s

import (
	"fmt"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestTunnelOpensAPortForwardTunnelToPod(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "")
	options.Namespace = uniqueID
	configData := fmt.Sprintf(EXAMPLE_POD_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)
	WaitUntilPodAvailable(t, options, "nginx-pod", 60, 1*time.Second)

	// Open a tunnel to pod from any available port locally
	localPort := GetAvailablePort(t)
	tunnel := NewTunnel(options, ResourceTypePod, "nginx-pod", localPort, 80)
	defer tunnel.Close()
	tunnel.ForwardPort(t)

	// Try to access the nginx service on the local port, retrying until we get a good response for up to 5 minutes
	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		fmt.Sprintf("http://localhost:%d", localPort),
		60,
		5*time.Second,
		func(statusCode int, body string) bool {
			return statusCode == 200
		},
	)
}

func TestTunnelOpensAPortForwardTunnelToService(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "")
	options.Namespace = uniqueID
	configData := fmt.Sprintf(EXAMPLE_POD_WITH_SERVICE_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)
	WaitUntilPodAvailable(t, options, "nginx-pod", 60, 1*time.Second)
	WaitUntilServiceAvailable(t, options, "nginx-service", 60, 1*time.Second)

	// Open a tunnel from any available port locally
	localPort := GetAvailablePort(t)
	tunnel := NewTunnel(options, ResourceTypeService, "nginx-service", localPort, 80)
	defer tunnel.Close()
	tunnel.ForwardPort(t)

	// Try to access the nginx service on the local port, retrying until we get a good response for up to 5 minutes
	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		fmt.Sprintf("http://localhost:%d", localPort),
		60,
		5*time.Second,
		func(statusCode int, body string) bool {
			return statusCode == 200
		},
	)
}

const EXAMPLE_POD_WITH_SERVICE_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  namespace: %s
  labels:
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.15.7
    ports:
    - containerPort: 80
    readinessProbe:
      httpGet:
        path: /
        port: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
  namespace: %s
spec:
  selector:
    app: nginx
  ports:
  - protocol: TCP
    targetPort: 80
    port: 80
`
