package k8s

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/gruntwork-io/terratest/modules/aws"
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
	clientset, err := GetKubernetesClientE(t)
	require.NoError(t, err)
	endpoint, err := GetServiceEndpointFromClientE(clientset, service, servicePort)
	require.NoError(t, err)
	return endpoint
}

// GetServiceEndpointFromClientE will return the service access point using the following logic:
// - For ClusterIP service type, return the URL that maps to ClusterIP and Service Port
// - For NodePort service type, identify the public IP of the node (if it exists, otherwise return the bound hostname),
//   and the assigned node port for the provided service port, and return the URL that maps to node ip and node port.
// - For LoadBalancer service type, return the publicly accessible hostname of the load balancer.
// - All other service types are not supported.
func GetServiceEndpointFromClientE(clientset *kubernetes.Clientset, service *corev1.Service, servicePort int) (string, error) {
	switch service.Spec.Type {
	case corev1.ServiceTypeClusterIP:
		// ClusterIP service type will map directly to service port
		return fmt.Sprintf("%s:%d", service.Spec.ClusterIP, servicePort), nil
	case corev1.ServiceTypeNodePort:
		return findEndpointForNodePortService(clientset, service, int32(servicePort))
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

// Extracts a endpoint that can be reached outside the kubernetes cluster. NodePort type needs to find the right
// allocated node port mapped to the service port, as well as find out the externally reachable ip (if available).
func findEndpointForNodePortService(
	clientset *kubernetes.Clientset,
	service *corev1.Service,
	servicePort int32,
) (string, error) {
	nodePort, err := findNodePort(service, int32(servicePort))
	if err != nil {
		return "", err
	}
	node, err := pickRandomNode(clientset)
	if err != nil {
		return "", err
	}
	nodeHostname, err := findNodeHostname(node)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", nodeHostname, nodePort), nil
}

// Given the desired servicePort, return the allocated nodeport
func findNodePort(service *corev1.Service, servicePort int32) (int32, error) {
	for _, port := range service.Spec.Ports {
		if port.Port == servicePort {
			return port.NodePort, nil
		}
	}
	return -1, NewKubernetesError("Port %d is not a part of the service %s", servicePort, service.Name)
}

// pickRandomNode will pick a random node in the kubernetes cluster
func pickRandomNode(clientset *kubernetes.Clientset) (corev1.Node, error) {
	nodes, err := GetNodesByClientE(clientset, metav1.ListOptions{})
	if err != nil {
		return corev1.Node{}, err
	}
	if len(nodes) == 0 {
		return corev1.Node{}, NewKubernetesError("There are no nodes in the cluster")
	}
	// TODO: randomly pick one
	return nodes[0], nil
}

// Given a node, return the ip address, preferring the external IP
func findNodeHostname(node corev1.Node) (string, error) {
	nodeIDUri, err := url.Parse(node.Spec.ProviderID)
	if err != nil {
		return "", err
	}
	switch nodeIDUri.Scheme {
	case "aws":
		return findAwsNodeHostname(node, nodeIDUri)
	default:
		return findDefaultNodeHostname(node)
	}
}

// findAwsNodeHostname will return the public ip of the node, assuming the node is an AWS EC2 instance.
// If the instance does not have a public IP, will return the internal hostname as recorded on the Kubernetes node
// object.
func findAwsNodeHostname(node corev1.Node, awsIDUri *url.URL) (string, error) {
	// Path is /AVAILABILITY_ZONE/INSTANCE_ID
	parts := strings.Split(awsIDUri.Path, "/")
	instanceID := parts[2]
	availabilityZone := parts[1]
	// Availability Zone name is known to be region code + 1 letter
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html
	region := availabilityZone[:len(availabilityZone)-1]

	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return "", err
	}
	ec2Client := ec2.New(sess)

	ipMap, err := aws.GetPublicIpsOfEc2InstancesFromClientE(ec2Client, []string{instanceID}, region)
	if err != nil {
		return "", err
	}

	publicIp, containsIp := ipMap[instanceID]
	if !containsIp {
		// return default hostname
		return findDefaultNodeHostname(node)
	}
	return publicIp, nil
}

// findDefaultNodeHostname returns the hostname recorded on the Kubernetes node object.
func findDefaultNodeHostname(node corev1.Node) (string, error) {
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeHostName {
			return address.Address, nil
		}
	}
	return "", NewKubernetesError("Node %s has no hostname", node.Name)
}
