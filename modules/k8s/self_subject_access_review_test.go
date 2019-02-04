// +build kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. To avoid overloading the system, we run the
// kubernetes tests separately from the others.

package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	authv1 "k8s.io/api/authorization/v1"
)

// NOTE: See service_account_test.go:TestGetServiceAccountWithAuthTokenGetsTokenThatCanBeUsedForAuth for the deny case,
// as the current authed user is assumed to be a super user and so there is nothing they can't do.

func TestCanIDoReturnsTrueForAllowedAction(t *testing.T) {
	t.Parallel()

	action := authv1.ResourceAttributes{
		Namespace: "kube-system",
		Verb:      "list",
		Resource:  "pod",
	}
	options := NewKubectlOptions("", "")
	assert.True(t, CanIDo(t, options, action))
}
