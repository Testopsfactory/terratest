package k8s

import (
	"fmt"
)

type KubernetesError struct {
	message string
}

func (err KubernetesError) Error() string {
	return err.message
}

func NewKubernetesError(message string, args ...interface{}) KubernetesError {
	return KubernetesError{fmt.Sprintf(message, args...)}
}
