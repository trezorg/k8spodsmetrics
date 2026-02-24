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

**Data Flow:** CLI flags → Config merge (file + CLI) → Service → Kubernetes API → Formatter → stdout

## Key Patterns

- **Config precedence:** CLI flags override YAML config file values
- **Output formatters:** Strategy pattern in `internal/adapters/stdout/{format}/`
- **Watch mode:** Screen wrappers refresh output at `--watch-period` intervals
- **Sorting:** Extensible via `internal/sorting/` with field-specific comparators
