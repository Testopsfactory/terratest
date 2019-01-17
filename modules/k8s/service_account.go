package k8s

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetServiceAccount returns a Kubernetes service account resource in the provided namespace with the given name. The
// namespace used is the one provided in the KubectlOptions. This will fail the test if there is an error.
func GetServiceAccount(t *testing.T, options *KubectlOptions, serviceAccountName string) *corev1.ServiceAccount {
	serviceAccount, err := GetServiceAccountE(t, options, serviceAccountName)
	require.NoError(t, err)
	return serviceAccount
}

// GetServiceAccount returns a Kubernetes service account resource in the provided namespace with the given name. The
// namespace used is the one provided in the KubectlOptions.
func GetServiceAccountE(t *testing.T, options *KubectlOptions, serviceAccountName string) (*corev1.ServiceAccount, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().ServiceAccounts(options.Namespace).Get(serviceAccountName, metav1.GetOptions{})
}

// CreateServiceAccount will create a new service account resource in the provided namespace with the given name. The
// namespace used is the one provided in the KubectlOptions. This will fail the test if there is an error.
func CreateServiceAccount(t *testing.T, options *KubectlOptions, serviceAccountName string) {
	require.NoError(t, CreateServiceAccountE(t, options, serviceAccountName))
}

// CreateServiceAccountE will create a new service account resource in the provided namespace with the given name. The
// namespace used is the one provided in the KubectlOptions.
func CreateServiceAccountE(t *testing.T, options *KubectlOptions, serviceAccountName string) error {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return err
	}

	serviceAccount := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: options.Namespace,
		},
	}
	_, err = clientset.CoreV1().ServiceAccounts(options.Namespace).Create(&serviceAccount)
	return err
}
