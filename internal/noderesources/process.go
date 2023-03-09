package noderesources

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/pkg/client"
	"github.com/trezorg/k8spodsmetrics/pkg/nodes"
	"github.com/trezorg/k8spodsmetrics/pkg/pods"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Config struct {
	KubeConfig  string
	KubeContext string
	LogLevel    string
	KLogLevel   uint
}

func (config Config) request(ctx context.Context, client corev1.CoreV1Interface) (NodeResourceList, error) {
	logger.Debug("Getting nodes info...")
	var nodeResources NodeResourceList
	errors := make([]error, 2)
	var podsList pods.PodResourceList
	var nodesList nodes.NodeList
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		nodesList, errors[0] = nodes.Nodes(ctx, client, nodes.NodeFilter{})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		podsList, errors[1] = pods.Pods(ctx, client, pods.PodFilter{})
	}()

	wg.Wait()

	var mErrs *multierror.Error

	for _, err := range errors {
		if err != nil {
			if err := multierror.Append(mErrs, err); err != nil {
				logger.Error("MultiError append error", err)
			}
		}
	}

	if err := mErrs.ErrorOrNil(); err != nil {
		return nodeResources, err
	}

	nodeResources = merge(podsList, nodesList)
	return nodeResources, nil
}

func (config Config) Request(ctx context.Context) (NodeResourceList, error) {
	var err error
	logger.Debug("Preparing client...")
	client, err := client.CoreV1Client(config.KubeConfig, config.KubeContext)
	if err != nil {
		return nil, err
	}
	return config.request(ctx, client)
}
