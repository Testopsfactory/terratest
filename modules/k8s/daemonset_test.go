package k8s

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestGetDaemonSetEReturnsErrorForNonExistantDaemonSet(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "")
	_, err := GetDaemonSetE(t, options, "sample-ds")
	require.Error(t, err)
}

func TestGetDaemonSetEReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "")
	options.Namespace = uniqueID
	configData := fmt.Sprintf(EXAMPLE_DAEMONSET_YAML_TEMPLATE, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	daemonSet := GetDaemonSet(t, options, "fluentd-elasticsearch")
	require.Equal(t, daemonSet.Name, "fluentd-elasticsearch")
	require.Equal(t, daemonSet.Namespace, uniqueID)
}

func TestListDaemonSetsReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "")
	options.Namespace = uniqueID
	configData := fmt.Sprintf(EXAMPLE_DAEMONSET_YAML_TEMPLATE, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	daemonSets := ListDaemonSets(t, options, metav1.ListOptions{})
	require.Equal(t, len(daemonSets), 1)

	daemonSet := daemonSets[0]
	require.Equal(t, daemonSet.Name, "fluentd-elasticsearch")
	require.Equal(t, daemonSet.Namespace, uniqueID)
}

const EXAMPLE_DAEMONSET_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd-elasticsearch
  namespace: %s
  labels:
    k8s-app: fluentd-logging
spec:
  selector:
    matchLabels:
      name: fluentd-elasticsearch
  template:
    metadata:
      labels:
        name: fluentd-elasticsearch
    spec:
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      containers:
      - name: fluentd-elasticsearch
        image: gcr.io/fluentd-elasticsearch/fluentd:v2.5.1
        resources:
          limits:
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 200Mi
        volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
      terminationGracePeriodSeconds: 30
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers

`
