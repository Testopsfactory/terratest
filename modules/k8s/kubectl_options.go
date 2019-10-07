package k8s

import "testing"

// KubectlOptions represents common options necessary to specify for all Kubectl calls
type KubectlOptions struct {
	ContextName string
	ConfigPath  string
	Namespace   string
	Env         map[string]string
}

// NewKubectlOptions will return a pointer to new instance of KubectlOptions with the configured options
func NewKubectlOptions(opts ...string) *KubectlOptions {
	switch {
	case len(opts) = 2:
		return &KubectlOptions{
			ContextName: opts[0],
			ConfigPath:  opts[1],
			Namespace:   opts[2],
		}
	default:
		return &KubectlOptions{
			ContextName: opts[0],
			ConfigPath: opts[1],
		}
	}
}

// GetConfigPath will return a sensible default if the config path is not set on the options.
func (kubectlOptions *KubectlOptions) GetConfigPath(t *testing.T) (string, error) {
	// We predeclare `err` here so that we can update `kubeConfigPath` in the if block below. Otherwise, go complains
	// saying `err` is undefined.
	var err error

	kubeConfigPath := kubectlOptions.ConfigPath
	if kubeConfigPath == "" {
		kubeConfigPath, err = GetKubeConfigPathE(t)
		if err != nil {
			return "", err
		}
	}
	return kubeConfigPath, nil
}
