k8spodsmetrics
====================================

Pods metrics
[![Actions Status]( https://github.com/trezorg/k8spodsmetrics/actions/workflows/go.yml/badge.svg)](https://github.com/trezorg/k8spodsmetrics/actions)

About
------------------------------------

The utility displays the consumption of memory and CPU, requests and limits in the context of kubernetes containers and nodes.
Differs from

    kubernetes top pods
    kubernetes top nodes

in that it adds requests and limits for each container.
Also shows total requests and limits by nodes.

Download
------------------------------------

    curl -sfL https://raw.githubusercontent.com/trezorg/k8spodsmetrics/main/install.sh | bash -s -- -d your_directory
    curl -sfL https://raw.githubusercontent.com/trezorg/k8spodsmetrics/main/install.sh | bash -s

Install
------------------------------------

    go get -u github.com/trezorg/k8spodsmetrics/cmd/k8spodsmetrics

Using
------------------------------------

    k8spodsmetrics --help

Logging
------------------------------------

    --loglevel controls verbosity (DEBUG|INFO|WARN|WARNING|ERROR)
    Logging uses the Go slog JSON handler
