package k8s

import (
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// IsMinikubeE returns true if the underlying kubernetes cluster is Minikube. This is determined by getting the
// associated nodes and checking if:
// - there is only one node
// - the node has at least one label namespaced with "minikube.k8s.io"
func IsMinikubeE(t testing.TestingT, options *KubectlOptions) (bool, error) {
	nodes, err := GetNodesE(t, options)
	if err != nil {
		return false, err
	}

	// ASSUMPTION: Minikube always only has one node.
	if len(nodes) != 1 {
		return false, nil
	}

	// ASSUMPTION: All minikube setups will have a node with labels that are namespaced with minikube.k8s.io
	node := nodes[0]
	labels := node.GetLabels()
	for key, _ := range labels {
		if strings.HasPrefix(key, "minikube.k8s.io") {
			return true, nil
		}
	}

	// At this point we know that the cluster has 1 node without the expected minikube label, so we assume it is not
	// minikube.
	return false, nil
}
