package k8s

// The following code is a fork of the Helm client. The main differences are:
// - Support testing context for better logging
// - Support resources other than pods
// See: https://github.com/helm/helm/blob/master/pkg/kube/tunnel.go

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/kubernetes/pkg/kubectl/generate"

	"github.com/gruntwork-io/terratest/modules/logger"
)

// KubeResourceType is an enum representing known resource types that can support port forwarding
type KubeResourceType int

const (
	Pod KubeResourceType = iota
	Service
)

func (resourceType KubeResourceType) String() string {
	switch resourceType {
	case Pod:
		return "pod"
	case Service:
		return "svc"
	default:
		// This should not happen
		return ""
	}
}

// Tunnel is the main struct that configures and manages port forwading tunnels to Kubernetes resources.
type Tunnel struct {
	KubectlOptions *KubectlOptions
	LocalPort      int
	RemotePort     int
	ResourceType   KubeResourceType
	ResourceName   string
	Out            io.Writer
	StopChan       chan struct{}
	ReadyChan      chan struct{}
}

// NewTunnel will create a new Tunnel struct.
func NewTunnel(kubectlOptions *KubectlOptions, resourceType KubeResourceType, resourceName string, local int, remote int) *Tunnel {
	return &Tunnel{
		KubectlOptions: kubectlOptions,
		LocalPort:      local,
		RemotePort:     remote,
		ResourceType:   resourceType,
		ResourceName:   resourceName,
		Out:            ioutil.Discard,
		StopChan:       make(chan struct{}, 1),
		ReadyChan:      make(chan struct{}, 1),
	}
}

// Close disconnects a tunnel connection by closing the StopChan, thereby stopping the goroutine.
func (tunnel *Tunnel) Close() {
	close(tunnel.StopChan)
}

// getAttachablePodForResource will find a pod that can be port forwarded to given the provided resource type and return
// the name.
func (tunnel *Tunnel) getAttachablePodForResourceE(t *testing.T) (string, error) {
	switch tunnel.ResourceType {
	case Pod:
		return tunnel.ResourceName, nil
	case Service:
		return tunnel.getAttachablePodForServiceE(t)
	default:
		return "", UnknownKubeResourceType{tunnel.ResourceType}
	}
}

// getAttachablePodForServiceE will find an active pod associated with the Service and return the pod name.
func (tunnel *Tunnel) getAttachablePodForServiceE(t *testing.T) (string, error) {
	service, err := GetServiceE(t, tunnel.KubectlOptions, tunnel.ResourceName)
	if err != nil {
		return "", err
	}
	selectorLabelsOfPods := generate.MakeLabels(service.Spec.Selector)
	servicePods, err := ListPodsE(t, tunnel.KubectlOptions, metav1.ListOptions{LabelSelector: selectorLabelsOfPods})
	if err != nil {
		return "", err
	}
	for _, pod := range servicePods {
		if IsPodAvailable(&pod) {
			return pod.Name, nil
		}
	}
	return "", ServiceNotAvailable{service}
}

// ForwardPort opens a tunnel to a kubernetes resource, as specified by the provided tunnel struct. This will fail the
// test if there is an error attempting to open the port.
func (tunnel *Tunnel) ForwardPort(t *testing.T) {
	require.NoError(t, tunnel.ForwardPortE(t))
}

// ForwardPortE opens a tunnel to a kubernetes resource, as specified by the provided tunnel struct.
func (tunnel *Tunnel) ForwardPortE(t *testing.T) error {
	logger.Logf(
		t,
		"Creating a port forwarding tunnel for resource %s/%s routing local port %d to remote port %d",
		tunnel.ResourceType.String(),
		tunnel.ResourceName,
		tunnel.LocalPort,
		tunnel.RemotePort,
	)

	// Prepare a kubernetes client for the client-go library
	clientset, err := GetKubernetesClientFromOptionsE(t, tunnel.KubectlOptions)
	if err != nil {
		logger.Logf(t, "Error creating a new Kubernetes client: %s", err)
		return err
	}
	kubeConfigPath, err := tunnel.KubectlOptions.GetConfigPath(t)
	if err != nil {
		logger.Logf(t, "Error getting kube config path: %s", err)
		return err
	}
	config, err := LoadApiClientConfigE(kubeConfigPath, tunnel.KubectlOptions.ContextName)
	if err != nil {
		logger.Logf(t, "Error loading Kubernetes config: %s", err)
		return err
	}

	// Find the pod to port forward to
	podName, err := tunnel.getAttachablePodForResourceE(t)
	if err != nil {
		logger.Logf(t, "Error finding available pod: %s", err)
		return err
	}
	logger.Logf(t, "Selected pod %s to open port forward to", podName)

	// Build a url to the portforward endpoint
	// example: http://localhost:8080/api/v1/namespaces/helm/pods/tiller-deploy-9itlq/portforward
	postEndpoint := clientset.CoreV1().RESTClient().Post()
	namespace := tunnel.KubectlOptions.Namespace
	portForwardCreateURL := postEndpoint.
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("portforward").
		URL()

	logger.Logf(t, "Using URL %s to create portforward", portForwardCreateURL)

	// Construct the spdy client required by the client-go portforward library
	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		logger.Logf(t, "Error creating http client: %s", err)
		return err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", portForwardCreateURL)

	// Construct a new PortForwarder struct that manages the instructed port forward tunnel
	ports := []string{fmt.Sprintf("%d:%d", tunnel.LocalPort, tunnel.RemotePort)}
	portforwarder, err := portforward.New(dialer, ports, tunnel.StopChan, tunnel.ReadyChan, tunnel.Out, tunnel.Out)
	if err != nil {
		logger.Logf(t, "Error creating port forwarding tunnel: %s", err)
		return err
	}

	// Open the tunnel in a goroutine so that it is available in the background. Report errors to the main goroutine via
	// a new channel.
	errChan := make(chan error)
	go func() {
		errChan <- portforwarder.ForwardPorts()
	}()

	// Wait for an error or the tunnel to be ready
	select {
	case err = <-errChan:
		logger.Logf(t, "Error starting port forwarding tunnel: %s", err)
		return err
	case <-portforwarder.Ready:
		logger.Logf(t, "Successfully created port forwarding tunnel")
		return nil
	}
}

// GetAvailablePort retrieves an available port on the host machine. This delegates the port selection to the golang net
// library by starting a server and then checking the port that the server is using. This will fail the test if it could
// not find an avilable port.
func GetAvailablePort(t *testing.T) int {
	port, err := GetAvailablePortE(t)
	require.NoError(t, err)
	return port
}

// GetAvailablePortE retrieves an available port on the host machine. This delegates the port selection to the golang net
// library by starting a server and then checking the port that the server is using.
func GetAvailablePortE(t *testing.T) (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()

	_, p, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return port, err
}
