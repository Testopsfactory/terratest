// +build kubeall kubernetes

package k8s

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetClusterRoleEReturnsErrorForNonExistantClusterRole(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "")
	_, err := GetClusterRoleE(t, options, "non-existing-role")
	require.Error(t, err)
}

func TestGetClusterRoleEReturnsCorrectClusterRoleInCorrectNamespace(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "")
	defer KubectlDeleteFromString(t, options, EXAMPLE_CLUSTER_ROLE_YAML_TEMPLATE)
	KubectlApplyFromString(t, options, EXAMPLE_CLUSTER_ROLE_YAML_TEMPLATE)

	role := GetClusterRole(t, options, "terratest-cluster-role")
	require.Equal(t, role.Name, "terratest-cluster-role")
}

const EXAMPLE_CLUSTER_ROLE_YAML_TEMPLATE = `---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: 'terratest-cluster-role'
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
`
