package stdin

import (
	"fmt"

	"github.com/trezorg/k8spodsmetrics/internal/resources"
	nodesorting "github.com/trezorg/k8spodsmetrics/internal/sorting/noderesources"
	"github.com/urfave/cli/v2"
)

func summaryFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "label",
			Aliases: []string{"l"},
			Value:   "",
			Usage:   "K8S node label",
		},
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Value:   "",
			Usage:   "K8S node name",
		},
		&cli.StringFlag{
			Name:    "sorting",
			Aliases: []string{"s"},
			Value:   "name",
			Usage:   fmt.Sprintf("Sorting. [%s]", nodesorting.StringListDefault()),
			Action: func(_ *cli.Context, value string) error {
				return nodesorting.Valid(nodesorting.Sorting(value))
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
