// +build kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. To avoid overloading the system, we run the
// kubernetes tests separately from the others.

package k8s

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests that:
// - kubectl is properly configured to talk to a kubernetes cluster
// - GetNodes will return a list of nodes registered with kubernetes
func TestGetNodes(t *testing.T) {
	t.Parallel()

	// Assumes local kubernetes (minikube or docker-for-desktop kube), where there is only one node
	nodes := GetNodes(t)
	require.Equal(t, len(nodes), 1)

	node := nodes[0]
	// Make sure node name is not blank, indicating an uninitialized Node object
	assert.NotEqual(t, node.Name, "")
}

// Tests that:
// - kubectl is properly configured to talk to a kubernetes cluster
// - GetReadyNodes will return a list of ready nodes registered with kubernetes
func TestGetReadyNodes(t *testing.T) {
	t.Parallel()

	// Assumes local kubernetes (minikube or docker-for-desktop kube), where there is only one node
	nodes := GetReadyNodes(t)
	require.Equal(t, len(nodes), 1)

	node := nodes[0]
	// Make sure node name is not blank, indicating an uninitialized Node object
	assert.NotEqual(t, node.Name, "")
}

// Tests that:
// - kubectl is properly configured to talk to a kubernetes cluster
// - WaitUntilAllNodesReady checks if all nodes in the cluster are ready
func TestWaitUntilAllNodesReady(t *testing.T) {
	t.Parallel()

	WaitUntilAllNodesReady(t, 12, 5*time.Second)

	nodes := GetNodes(t)
	nodeNames := map[string]bool{}
	for _, node := range nodes {
		nodeNames[node.Name] = true
	}

	readyNodes := GetReadyNodes(t)
	readyNodeNames := map[string]bool{}
	for _, node := range readyNodes {
		readyNodeNames[node.Name] = true
	}

	assert.Equal(t, nodeNames, readyNodeNames)
}
