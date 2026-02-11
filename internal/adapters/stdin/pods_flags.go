package stdin

import (
	"fmt"

	"github.com/trezorg/k8spodsmetrics/internal/resources"
	metricssorting "github.com/trezorg/k8spodsmetrics/internal/sorting/metricsresources"
	"github.com/urfave/cli/v2"
)

func podsFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "namespace",
			Aliases: []string{"n"},
			Value:   "",
			Usage:   "K8S namespace",
		},
		&cli.StringFlag{
			Name:    "label",
			Aliases: []string{"l"},
			Value:   "",
			Usage:   "K8S pod label",
		},
		&cli.StringFlag{
			Name:    "field-selector",
			Aliases: []string{"f"},
			Value:   "",
			Usage:   "K8S field selector",
		},
		&cli.StringSliceFlag{
			Name:    "node",
			Aliases: []string{"nd", "nodes"},
			Usage:   "K8S node names",
		},
		&cli.StringFlag{
			Name:    "sorting",
			Aliases: []string{"s"},
			Value:   "namespace",
			Usage:   fmt.Sprintf("Sorting. [%s]", metricssorting.StringListDefault()),
			Action: func(_ *cli.Context, value string) error {
				return metricssorting.Valid(metricssorting.Sorting(value))
			},
		},
		&cli.BoolFlag{
			Name:    "reverse",
			Aliases: []string{"r"},
			Value:   false,
			Usage:   "Reverse sort",
		},
		&cli.StringSliceFlag{
			Name:    "resource",
			Aliases: []string{"res"},
			Value:   cli.NewStringSlice(string(resources.All)),
			Usage:   fmt.Sprintf("Resources. [%s]", resources.StringListDefault()),
			Action: func(_ *cli.Context, value []string) error {
				outputResources := resources.FromStrings(value...)
				return resources.Valid(outputResources...)
			},
		},
	}
}
