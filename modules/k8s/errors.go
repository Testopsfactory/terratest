package k8s

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// ServiceNotAvailable is returned when a Kubernetes service is not yet available to accept traffic.
type ServiceNotAvailable struct {
	service *corev1.Service
}

func (err ServiceNotAvailable) Error() string {
	return fmt.Sprintf("Service %s is not available", err.service.Name)
}

func NewServiceNotAvailableError(service *corev1.Service) ServiceNotAvailable {
	return ServiceNotAvailable{service}
}

// UnknownServiceType is returned when a Kubernetes service has a type that is not yet handled by the test functions.
type UnknownServiceType struct {
	service *corev1.Service
}

func (err UnknownServiceType) Error() string {
	return fmt.Sprintf("Service %s has an unknown service type", err.service.Name)
}

func NewUnknownServiceTypeError(service *corev1.Service) UnknownServiceType {
	return UnknownServiceType{service}
}

// UnknownServicePort is returned when the given service port is not an exported port of the service.
type UnknownServicePort struct {
	service *corev1.Service
	port    int32
}

func (err UnknownServicePort) Error() string {
	return fmt.Sprintf("Port %d is not a part of the service %s", err.port, err.service.Name)
}

func NewUnknownServicePortError(service *corev1.Service, port int32) UnknownServicePort {
	return UnknownServicePort{service, port}
}

// NoNodesInKubernetes is returned when the Kubernetes cluster has no nodes registered.
type NoNodesInKubernetes struct{}

func (err NoNodesInKubernetes) Error() string {
	return "There are no nodes in the Kubernetes cluster"
}

func NewNoNodesInKubernetesError() NoNodesInKubernetes {
	return NoNodesInKubernetes{}
}

// NodeHasNoHostname is returned when a Kubernetes node has no discernible hostname
type NodeHasNoHostname struct {
	node *corev1.Node
}

func (err NodeHasNoHostname) Error() string {
	return fmt.Sprintf("Node %s has no hostname", err.node.Name)
}

func NewNodeHasNoHostnameError(node *corev1.Node) NodeHasNoHostname {
	return NodeHasNoHostname{node}
}

// MalformedNodeID is returned when a Kubernetes node has a malformed node id scheme
type MalformedNodeID struct {
	node *corev1.Node
}

func (err MalformedNodeID) Error() string {
	return fmt.Sprintf("Node %s has malformed ID %s", err.node.Name, err.node.Spec.ProviderID)
}

func NewMalformedNodeIDError(node *corev1.Node) MalformedNodeID {
	return MalformedNodeID{node}
}
