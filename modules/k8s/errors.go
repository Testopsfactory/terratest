package k8s

type KubernetesError struct {
	message string
}

func (err KubernetesError) Error() string {
	return err.message
}

func NewKubernetesError(message string) KubernetesError {
	return KubernetesError{message}
}
