package client

import (
	"os"
	"strconv"

	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog/v2"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned/typed/metrics/v1beta1"
)

const (
	defaultClientQPS   float32 = 10
	defaultClientBurst int     = 20

	clientQPSEnvVar   = "K8SPODSMETRICS_CLIENT_QPS"
	clientBurstEnvVar = "K8SPODSMETRICS_CLIENT_BURST"
)

func rateLimitFromEnv() (float32, int) {
	qps := defaultClientQPS
	burst := defaultClientBurst

	if rawQPS := os.Getenv(clientQPSEnvVar); rawQPS != "" {
		parsedQPS, err := strconv.ParseFloat(rawQPS, 32)
		if err != nil || parsedQPS <= 0 {
			klog.Warningf("invalid %s=%q, using default %.2f", clientQPSEnvVar, rawQPS, defaultClientQPS)
		} else {
			qps = float32(parsedQPS)
		}
	}

	if rawBurst := os.Getenv(clientBurstEnvVar); rawBurst != "" {
		parsedBurst, err := strconv.Atoi(rawBurst)
		if err != nil || parsedBurst <= 0 {
			klog.Warningf("invalid %s=%q, using default %d", clientBurstEnvVar, rawBurst, defaultClientBurst)
		} else {
			burst = parsedBurst
		}
	}

	return qps, burst
}

func applyRateLimit(config *rest.Config) {
	qps, burst := rateLimitFromEnv()
	config.QPS = qps
	config.Burst = burst
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(qps, burst)
}

func restConfig(kubeconfigPath string, context string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		klog.Warning("--kubeconfig was not specified. Using the inClusterConfig.  This might not work.")
		kubeconfig, err := rest.InClusterConfig()
		if err == nil {
			applyRateLimit(kubeconfig)
			return kubeconfig, nil
		}
		klog.Warning("error creating inClusterConfig, falling back to default config: ", err)
	}
	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}
	if context != "" {
		configOverrides.CurrentContext = context
	}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, err
	}
	applyRateLimit(config)
	return config, nil
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

func CoreV1Client(kubeconfigPath string, context string) (corev1.CoreV1Interface, error) {
	config, err := restConfig(kubeconfigPath, context)
	if err != nil {
		return nil, err
	}
	pc, err := podsClient(config)
	if err != nil {
		return nil, err
	}
	return pc, nil
}

func ForMetrics(kubeconfigPath string, context string) (metricsv1beta1.MetricsV1beta1Interface, error) {
	config, err := restConfig(kubeconfigPath, context)
	if err != nil {
		return nil, err
	}
	mc, err := metricsClient(config)
	if err != nil {
		return nil, err
	}
	return mc, nil
}

func FindKubeConfig() (string, error) {
	env := os.Getenv("KUBECONFIG")
	if env != "" {
		return env, nil
	}
	path, err := homedir.Expand("~/.kube/config")
	if err != nil {
		return "", err
	}
	return path, nil
}
