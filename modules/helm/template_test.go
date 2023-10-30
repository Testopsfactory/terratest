//go:build kubeall || helm
// +build kubeall helm

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests, and further differentiate helm
// tests. This is done because minikube is heavy and can interfere with docker related tests in terratest. Similarly,
// helm can overload the minikube system and thus interfere with the other kubernetes tests. To avoid overloading the
// system, we run the kubernetes tests and helm tests separately from the others.

package helm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

const (
	remote2ChartSource  = "https://charts.bitnami.com/bitnami"
	remote2ChartName    = "nginx"
	remote2ChartVersion = "13.2.23"
)

// Test that we can render locally a remote chart (e.g bitnami/nginx)
func TestRemoteChartRender(t *testing.T) {
	t.Parallel()

	namespaceName := fmt.Sprintf(
		"%s-%s",
		strings.ToLower(t.Name()),
		strings.ToLower(random.UniqueId()),
	)

	releaseName := "keda"

	options := &Options{
		SetValues: map[string]string{
			"metricsServer.replicaCount":           "999",
			"resources.metricServer.limits.memory": "1234Mi",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", namespaceName),
	}

	// Run RenderTemplate to render the template and capture the output. Note that we use the version without `E`, since
	// we want to assert that the template renders without any errors.
	// Additionally, although we know there is only one yaml file in the template, we deliberately path a templateFiles
	// arg to demonstrate how to select individual templates to render.
	output := RenderRemoteTemplate(t, options, "https://kedacore.github.io/charts", releaseName, []string{"templates/metrics-server/deployment.yaml"})

	// Now we use kubernetes/client-go library to render the template output into the Deployment struct. This will
	// ensure the Deployment resource is rendered correctly.
	var deployment appsv1.Deployment
	UnmarshalK8SYaml(t, output, &deployment)

	// Verify the namespace matches the expected supplied namespace.
	require.Equal(t, namespaceName, deployment.Namespace)

	// Finally, we verify the deployment pod template spec is set to the expected container image value
	var expectedMetricsServerReplica int32
	expectedMetricsServerReplica = 999
	deploymentMetricsServerReplica := *deployment.Spec.Replicas
	require.Equal(t, expectedMetricsServerReplica, deploymentMetricsServerReplica)
}
