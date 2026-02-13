# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Quick Reference

**Build & Test:**
```bash
task build          # Build binary (defaults: GOOS=linux GOARCH=amd64)
task test           # Run tests with race detector
task lint           # Run golangci-lint + go vet
task check          # Full check: format, lint, test (run before commits)
```

**Run locally:**
```bash
go run ./cmd/k8spodsmetrics --help
go run ./cmd/k8spodsmetrics pods --namespace default
go run ./cmd/k8spodsmetrics summary
```

## Architecture

```
cmd/k8spodsmetrics/       # CLI entrypoint (main)
internal/
├── adapters/
│   ├── stdin/           # CLI flags and config loading (urfave/cli)
│   └── stdout/          # Output formatters (table/json/yaml/string/screen)
├── config/              # YAML config file parsing
├── metricsresources/    # Pod metrics service (fetch, aggregate, sort)
├── noderesources/       # Node summary service
├── sorting/             # Sort strategies for metrics/nodes
├── resources/           # Resource type parsing (cpu/memory)
├── humanize/            # Byte/unit formatting
├── alert/               # Alert threshold handling
└── choices/             # Enum validation (output format, sort fields)
pkg/                      # Public APIs (client, pods, nodes, podmetrics, nodemetrics)
```

**Data Flow:** CLI flags → Config merge (file + CLI) → Service → Kubernetes API → Formatter → stdout

## Key Patterns

- **Config precedence:** CLI flags override YAML config file values
- **Output formatters:** Strategy pattern in `internal/adapters/stdout/{format}/`
- **Watch mode:** Screen wrappers refresh output at `--watch-period` intervals
- **Sorting:** Extensible via `internal/sorting/` with field-specific comparators

## Dependencies

- CLI: `github.com/urfave/cli/v2`
- Kubernetes: `k8s.io/client-go`, `k8s.io/metrics`
- Output: `github.com/jedib0t/go-pretty/v6` (tables)
- Testing: `github.com/stretchr/testify`

See `AGENTS.md` for detailed guidelines on coding style, testing, and PR requirements.
