k8spodsmetrics
====================================

Pods metrics
[![Actions Status]( https://github.com/trezorg/k8spodsmetrics/actions/workflows/go.yml/badge.svg)](https://github.com/trezorg/k8spodsmetrics/actions)

Download
------------------------------------

    tag_name=$(curl -s https://api.github.com/repos/trezorg/k8spodsmetrics/releases/latest | jq -r '.tag_name')
    curl --fail-with-body -sL \
        "https://github.com/trezorg/k8spodsmetrics/releases/download/${tag_name}/k8spodsmetrics-$(go env GOOS)-$(go env GOARCH)" -o \
        "$(go env GOBIN)/k8spodsmetrics" && \
        chmod +x "$(go env GOBIN)/k8spodsmetrics" && \
        "$(go env GOBIN)/k8spodsmetrics" --help

Install
------------------------------------

    go get -u github.com/trezorg/k8spodsmetrics/cmd/k8spodsmetrics

Using
------------------------------------

    k8spodsmetrics --help
