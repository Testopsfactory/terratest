package k8s

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
)

// GetPod returns a Kubernetes pod resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetPod(t *testing.T, options *KubectlOptions, podName string) *corev1.Pod {
	pod, err := GetPodE(t, options, podName)
	require.NoError(t, err)
	return pod
}

// GetPodE returns a Kubernetes pod resource in the provided namespace with the given name.
func GetPodE(t *testing.T, options *KubectlOptions, podName string) (*corev1.Pod, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().Pods(options.Namespace).Get(podName, metav1.GetOptions{})
}

// WaitUntilPodAvailable waits until the pod is running.
func WaitUntilPodAvailable(t *testing.T, options *KubectlOptions, podName string, retries int, sleepBetweenRetries time.Duration) {
	statusMsg := fmt.Sprintf("Wait for pod %s to be provisioned.", podName)
	message := retry.DoWithRetry(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pod, err := GetPodE(t, options, podName)
			if err != nil {
				return "", err
			}
			if !IsPodAvailable(pod) {
				return "", NewPodNotAvailableError(pod)
			}
			return "Pod is now available", nil
		},
	)
	logger.Logf(t, message)
}

// IsPodAvailable returns true if the pod is running.
func IsPodAvailable(pod *corev1.Pod) bool {
	return pod.Status.Phase == corev1.PodRunning
}
