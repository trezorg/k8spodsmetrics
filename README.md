k8spodsmetrics
====================================

Pods metrics
[![Actions Status]( https://github.com/trezorg/k8spodsmetrics/actions/workflows/go.yml/badge.svg)](https://github.com/trezorg/k8spodsmetrics/actions)

Download
------------------------------------

    curl --fail-with-body -sL "https://github.com/trezorg/k8spodsmetrics/releases/download/v0.0.2/k8spodsmetrics-$(go env GOOS)-$(go env GOARCH)" -o \
        "$(go env GOBIN)/k8spodsmetrics" && echo "ok"

Install
------------------------------------

    go get -u github.com/trezorg/k8spodsmetrics/cmd/k8spodsmetrics

Using
------------------------------------

    k8spodsmetrics --help
