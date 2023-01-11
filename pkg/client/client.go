package client

import (
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

func restConfig(kubeconfigPath string, context string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		klog.Warning("Neither --kubeconfig nor --master was specified.  Using the inClusterConfig.  This might not work.")
		kubeconfig, err := rest.InClusterConfig()
		if err == nil {
			return kubeconfig, nil
		}
		klog.Warning("error creating inClusterConfig, falling back to default config: ", err)
	}
	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}
	if context != "" {
		configOverrides.CurrentContext = context
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
}

func metricsClient(config *rest.Config) (metricsv1beta1.MetricsV1beta1Interface, error) {
	client, err := metrics.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client.MetricsV1beta1(), nil
}

func podsClient(config *rest.Config) (corev1.CoreV1Interface, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client.CoreV1(), nil
}

func Clients(kubeconfigPath string, context string) (metricsv1beta1.MetricsV1beta1Interface, corev1.CoreV1Interface, error) {
	config, err := restConfig(kubeconfigPath, context)
	if err != nil {
		return nil, nil, err
	}
	mc, err := metricsClient(config)
	if err != nil {
		return nil, nil, err
	}
	pc, err := podsClient(config)
	if err != nil {
		return nil, nil, err
	}
	return mc, pc, nil
}
