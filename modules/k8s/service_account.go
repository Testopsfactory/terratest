package k8s

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/gruntwork-cli/errors"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
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

// CreateServiceAccountWithAuthToken will create a new ServiceAccount resource in the configured namespace and generate
// an auth token that can be used to configure Kubectl to login as that service account. This will fail the test if
// there is an error.
func CreateServiceAccountWithAuthToken(t *testing.T, kubectlOptions *KubectlOptions, name string) string {
	token, err := CreateServiceAccountWithAuthTokenE(t, kubectlOptions, name)
	require.NoError(t, err)
	return token
}

// CreateServiceAccountWithAuthTokenE will create a new ServiceAccount resource in the configured namespace and generate
// an auth token that can be used to configure Kubectl to login as that service account.
func CreateServiceAccountWithAuthTokenE(t *testing.T, kubectlOptions *KubectlOptions, name string) (string, error) {
	// Create a new service account that we will use for auth.
	if err := CreateServiceAccountE(t, kubectlOptions, name); err != nil {
		return "", err
	}

	// Wait for the TokenController to provision a ServiceAccount token
	msg, err := retry.DoWithRetryE(
		t,
		"Waiting for ServiceAccount Token to be provisioned",
		30,
		10*time.Second,
		func() (string, error) {
			logger.Logf(t, "Checking if service account has secret")
			serviceAccount := GetServiceAccount(t, kubectlOptions, name)
			if len(serviceAccount.Secrets) == 0 {
				msg := "No secrets on the service account yet"
				logger.Logf(t, msg)
				return "", fmt.Errorf(msg)
			}
			return "Service Account has secret", nil
		},
	)
	if err != nil {
		return "", err
	}
	logger.Logf(t, msg)

	// Then get the service account token
	serviceAccount, err := GetServiceAccountE(t, kubectlOptions, name)
	if err != nil {
		return "", err
	}
	if len(serviceAccount.Secrets) != 1 {
		return "", errors.WithStackTrace(ServiceAccountTokenNotAvailable{name})
	}
	secret := GetSecret(t, kubectlOptions, serviceAccount.Secrets[0].Name)
	return string(secret.Data["token"]), nil
}

// AddConfigContextForServiceAccountE will add a new config context that binds the ServiceAccount auth token to the
// Kubernetes cluster of the current config context.
func AddConfigContextForServiceAccountE(
	t *testing.T,
	kubectlOptions *KubectlOptions,
	contextName string,
	serviceAccountName string,
	token string,
) error {
	// First load the config context
	config := LoadConfigFromPath(kubectlOptions.ConfigPath)
	rawConfig, err := config.RawConfig()
	if err != nil {
		return errors.WithStackTrace(err)
	}

	// Next get the current cluster
	currentContext := rawConfig.Contexts[rawConfig.CurrentContext]
	currentCluster := currentContext.Cluster

	// Now insert the auth info for the service account
	rawConfig.AuthInfos[serviceAccountName] = &api.AuthInfo{Token: token}

	// We now have enough info to add the new context
	UpsertConfigContext(&rawConfig, contextName, currentCluster, serviceAccountName)

	// Finally, overwrite the config
	if err := clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, false); err != nil {
		return errors.WithStackTrace(err)
	}
	return nil
}
