package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/trezorg/k8spodsmetrics/internal/logger"
	"github.com/trezorg/k8spodsmetrics/internal/metricsresources"
	"github.com/trezorg/k8spodsmetrics/pkg/client"
	"golang.org/x/exp/slog"
	"k8s.io/klog/v2"
)

var version = "0.0.1"

func processK8sMetrics(config metricsresources.Config) error {
	logger.InitLogger(config.LogLevel)
	if config.KubeConfig == "" {
		var err error
		config.KubeConfig, err = client.FindKubeConfig()
		if err != nil {
			return err
		}
	}
	klog.InitFlags(nil)
	defer klog.Flush()
	if err := flag.Set("v", strconv.Itoa(int(config.KLogLevel))); err != nil {
		return err
	}
	ctx, done := context.WithCancel(context.Background())
	defer done()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for sig := range stop {
			logger.Info("Got OS signal", slog.String("signal", sig.String()))
			done()
			return
		}
	}()

	if config.WatchMetrics {
		config.Watch(
			ctx,
			func(rList metricsresources.PodMetricsResourceList) { fmt.Println(rList) },
			func(err error) { logger.Error("", err) },
		)
		return nil
	} else {
		resources, err := config.Request(ctx)
		if err != nil {
			return fmt.Errorf("Cannot get k8s resources: %w", err)
		}
		fmt.Println(resources)
		return nil

	}
}

func main() {
	if err := processArgs(); err != nil {
		log.Fatal(err)
	}
}
