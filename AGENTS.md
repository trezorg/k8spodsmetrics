# Repository Guidelines

## Project Structure & Module Organization
- `cmd/k8spodsmetrics/`: CLI entrypoint (`main`).
- `internal/`: non-exported app code (adapters, output formats, sorting, resources, logger, alerts, humanize).
  - Examples: `internal/adapters/stdin` (CLI wiring), `internal/adapters/stdout/{table,json,yaml,string,screen}` (printers), `internal/{metricsresources,noderesources}` (domain services).
- `pkg/`: public Go APIs (`client`, `pods`, `nodes`, `podmetrics`, `nodemetrics`). Preferred for reuse.
- `build/`: compiled binaries. CI/artifacts target.
- Tests live beside code as `*_test.go`.

## Architecture Overview
- CLI layer: `cmd/k8spodsmetrics` + `internal/adapters/stdin` (urfave/cli) parses flags and builds configs.
- Services: `internal/{metricsresources,noderesources}` fetch/aggregate data and enforce sorting/filters.
- Kubernetes access: `pkg/client` constructs CoreV1 and Metrics clients; domain packages in `pkg/` provide helpers.
- Presentation: `internal/adapters/stdout/{table,json,yaml,string,screen}` implement output strategies; `screen` wraps watch mode.
- Utilities: `internal/{sorting,resources,humanize,alert,logger}` provide common transforms and validation.
- Flow: CLI → Config → Service → `pkg/client` → Formatter → stdout; watch mode uses screen wrappers for live updates.

## Build, Test, and Development Commands
- `make build` — builds CLI to `build/k8spodsmetrics-${GOOS}-${GOARCH}`. Env vars default to `linux/amd64` (override with `GOOS=darwin GOARCH=arm64`).
- `make test` — runs `go test -v -race -count 1 ./...`.
- `make lint` — runs `golangci-lint` and `go vet` (config: `.golangci.yaml`).
- Run locally: `go run ./cmd/k8spodsmetrics --help`.
- Cross-compile matrix: `bash build.sh`.
- Install a release binary: `bash install.sh [-d ~/bin] [-v vX.Y.Z]`.

## Coding Style & Naming Conventions
- Go formatting is mandatory: `gofmt`/`goimports` and `golangci-lint` must pass (`make lint`).
- Idiomatic Go naming: exported `CamelCase`, unexported `lowerCamelCase`; package names short, lowercase.
- Keep packages cohesive; prefer `internal/` for non-public code, `pkg/` for reusable APIs.
- Avoid breaking public `pkg/` APIs without a semver bump.

## Testing Guidelines
- Frameworks: `testing` + `github.com/stretchr/testify/require`.
- Use table-driven tests; name as `TestXxx` and colocate with code.
- Cover sorting, resource parsing, and formatters; include edge cases.
- Run `make test` locally; race detector must be clean.

## Commit & Pull Request Guidelines
- Commits: concise, imperative style (e.g., "Refactor metricsresources", "Add pods resource filter").
- PRs: clear description, rationale, and linked issue; include sample CLI output/screenshot for UX changes.
- CI readiness: `make test lint` must pass; update `README.md` when flags or behavior change.

## Security & Configuration Tips
- Never commit kubeconfigs or secrets. Use `KUBECONFIG` or `--kubeconfig`/`--context` flags.
- Requires Kubernetes Metrics API on the target cluster; verify access before running.
