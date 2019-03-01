// +build kubeall kubernetes

package k8s

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetRoleEReturnsErrorForNonExistantRole(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "")
	_, err := GetRoleE(t, options, "non-existing-role")
	require.Error(t, err)
}

func TestGetRoleEReturnsCorrectRoleInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "")
	options.Namespace = uniqueID
	configData := fmt.Sprintf(EXAMPLE_ROLE_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	role := GetRole(t, options, "terratest-role")
	require.Equal(t, role.Name, "terratest-role")
	require.Equal(t, role.Namespace, uniqueID)
}

const EXAMPLE_ROLE_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: '%s'
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: 'terratest-role'
  namespace: '%s'
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
`
