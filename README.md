k8spodsmetrics
====================================

Pods metrics
[![Actions Status]( https://github.com/trezorg/k8spodsmetrics/actions/workflows/ci.yml/badge.svg)](https://github.com/trezorg/k8spodsmetrics/actions)

About
------------------------------------

The utility displays the consumption of memory and CPU, requests and limits in the context of kubernetes containers and nodes.
Differs from

    kubectl top pods
    kubectl top nodes

in that it adds requests and limits for each container.
Also shows total requests and limits by nodes.

Download
------------------------------------

    curl -sfL https://raw.githubusercontent.com/trezorg/k8spodsmetrics/main/install.sh | bash -s -- -d your_directory
    curl -sfL https://raw.githubusercontent.com/trezorg/k8spodsmetrics/main/install.sh | bash -s
    curl -sfL https://raw.githubusercontent.com/trezorg/k8spodsmetrics/main/install.sh | bash -s -- --checksums

`--checksums` is optional, but recommended to verify release artifact integrity during installation.

Install
------------------------------------

    go install github.com/trezorg/k8spodsmetrics/cmd/k8spodsmetrics@latest

Using
------------------------------------

    k8spodsmetrics --help

Logging
------------------------------------

    --loglevel controls verbosity (DEBUG|INFO|WARN|WARNING|ERROR)
    Logging uses the Go slog JSON handler
    --timeout sets Kubernetes API timeout in seconds (default: 30)

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
  table-view: compact
  alert: cpu
  columns:
    - request
    - limit
    - used
  watch-period: 10
  watch: true
  timeout: 45

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

**Merge Behavior:** CLI flags take precedence over file config values. Empty/zero values from CLI are replaced with file config values. For boolean flags, file values are used unless the CLI flag is explicitly set, so `--watch=false` and `--reverse=false` override `true` values from the config file. For timeout, the config `common.timeout` value is used unless `--timeout` is explicitly provided.

Table Views
------------------------------------

When `--output table` is selected, you can choose between two table layouts:

- `--table-view expanded` keeps the original wide matrix view and supports `--columns`
- `--table-view compact` renders a denser single-header summary view for `pods` and `summary` and is the default

Examples:

    k8spodsmetrics --output table --table-view compact pods --resources cpu,memory
    k8spodsmetrics --output table --table-view compact summary --resources all

`--columns` implies `--table-view expanded` when no table view is explicitly set. An explicit `--table-view compact --columns ...` combination is still rejected.
