k8spodsmetrics
====================================

Pods metrics
[![Actions Status]( https://github.com/trezorg/k8spodsmetrics/actions/workflows/ci.yml/badge.svg)](https://github.com/trezorg/k8spodsmetrics/actions)

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

Configuration File
------------------------------------

You can specify default values in a YAML configuration file using the `--config` flag:

    k8spodsmetrics --config /path/to/config.yaml pods

Example configuration file:

```yaml
common:
  kubeconfig: /path/to/kubeconfig
  context: my-context
  output: json
  alert: cpu
  kloglevel: 2
  watch-period: 10
  watch: true

pods:
  namespace: default
  label: app=nginx
  field-selector: status.phase=Running
  nodes:
    - node1
    - node2
  sorting: name
  reverse: true
  resources:
    - cpu
    - memory

summary:
  name: node-name
  label: kubernetes.io/role=master
  sorting: used_cpu
  reverse: false
  resources:
    - all
```

**Merge Behavior:** CLI flags take precedence over file config values. Empty/zero values from CLI are replaced with file config values. Note that boolean flags have a limitation: if the file has `watch: true` or `reverse: true`, the CLI default of `false` cannot override it.
