package metricsresources

import (
	"fmt"

	screen "github.com/aditya43/clear-shell-screen-golang"
	"github.com/trezorg/k8spodsmetrics/internal/logger"
)

type SuccessFunction func(PodMetricsResourceList)
type ErrorFunction func(error)

func ScreenSuccessWriter() SuccessFunction {
	return func(rList PodMetricsResourceList) {
		screen.Clear()
		screen.MoveTopLeft()
		fmt.Println(rList)
		screen.MoveTopLeft()
	}
}

func ScreenErrorWriter() ErrorFunction {
	return func(err error) {
		screen.Clear()
		screen.MoveTopLeft()
		logger.Error("", err)
	}
}
