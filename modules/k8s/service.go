package k8s

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
)

// GetService returns a Kubernetes service resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetService(t *testing.T, namespace string, serviceName string) *corev1.Service {
	clientset, err := GetKubernetesClientE(t)
	require.NoError(t, err)
	service, err := GetServiceFromClientE(clientset, namespace, serviceName)
	require.NoError(t, err)
	return service
}

// GetServiceFromClientE returns a Kubernetes service resource in the provided namespace with the given name.
func GetServiceFromClientE(clientset *kubernetes.Clientset, namespace string, serviceName string) (*corev1.Service, error) {
	return clientset.CoreV1().Services(namespace).Get(serviceName, metav1.GetOptions{})
}

// WaitUntilServiceAvailable waits until the service endpoint is ready to accept traffic.
func WaitUntilServiceAvailable(t *testing.T, namespace string, serviceName string, retries int, sleepBetweenRetries time.Duration) {
	clientset, err := GetKubernetesClientE(t)
	require.NoError(t, err)
	statusMsg := fmt.Sprintf("Wait for service %s to be provisioned.", serviceName)
	message, err := retry.DoWithRetryE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			service, err := GetServiceFromClientE(clientset, namespace, serviceName)
			if err != nil {
				return "", err
			}
			if !IsServiceAvailable(service) {
				return "", NewKubernetesError("Service is not available")
			}
			return "Service is now available", nil
		},
	)
	if err != nil {
		logger.Logf(t, "Error waiting for service %s to be available: %s", serviceName, err)
		t.Fatal(err)
	}
	logger.Logf(t, message)
}

// IsServiceAvailable returns true if the service endpoint is ready to accept traffic.
func IsServiceAvailable(service *corev1.Service) bool {
	// Only the LoadBalancer type has a delay. All other service types are available if the resource exists.
	switch service.Spec.Type {
	case corev1.ServiceTypeLoadBalancer:
		ingress := service.Status.LoadBalancer.Ingress
		// The load balancer is ready if it has at least one ingress point
		return len(ingress) > 0
	default:
		return true
	}
}

// GetServiceEndpoint will return the service access point. If the service endpoint is not ready, will fail the test
// immediately.
func GetServiceEndpoint(t *testing.T, service *corev1.Service, servicePort int) string {
	endpoint, err := GetServiceEndpointE(service, servicePort)
	require.NoError(t, err)
	return endpoint
}

// GetServiceEndpointE will return the service access point
func GetServiceEndpointE(service *corev1.Service, servicePort int) (string, error) {
	switch service.Spec.Type {
	case corev1.ServiceTypeClusterIP:
		// ClusterIP service type will map directly to service port
		return fmt.Sprintf("%s:%d", service.Spec.ClusterIP, servicePort), nil
	case corev1.ServiceTypeNodePort:
		// NodePort type needs to find the right port mapped to the service port
		nodePort, err := findNodePort(service, int32(servicePort))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s:%d", service.Spec.ExternalIPs[0], nodePort), nil
	case corev1.ServiceTypeLoadBalancer:
		ingress := service.Status.LoadBalancer.Ingress
		if len(ingress) == 0 {
			return "", NewKubernetesError("Load Balancer is not ready")
		}
		// Load Balancer service type will map directly to service port
		return fmt.Sprintf("%s:%d", ingress[0].Hostname, servicePort), nil
	default:
		return "", NewKubernetesError("Unknown service type")
	}
}

func findNodePort(service *corev1.Service, servicePort int32) (int32, error) {
	for _, port := range service.Spec.Ports {
		if port.Port == servicePort {
			return port.NodePort, nil
		}
	}
	return -1, NewKubernetesError(fmt.Sprintf("Port %d is not a part of the service %s", servicePort, service.Name))
}
