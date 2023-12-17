package main

import (
	"log"

	"github.com/trezorg/k8spodsmetrics/internal/adapters/stdin"
)

var version = "0.0.1"

func main() {
	if err := stdin.Start(version); err != nil {
		log.Fatal(err)
	}
}
