# Tasks: Consolidate Default Separator

## 1. Add shared constant to choices package

- [x] 1.1 Add `DefaultSeparator` constant to `internal/choices/choices.go`

## 2. Update packages to use shared constant

- [x] 2.1 Update `internal/sorting/metricsresources/metricsresources.go`
- [x] 2.2 Update `internal/sorting/noderesources/noderesources.go`
- [x] 2.3 Update `internal/alert/formats.go`
- [x] 2.4 Update `internal/resources/resource.go`
- [x] 2.5 Update `internal/output/formats.go`

## 3. Verify

- [x] 3.1 Run tests to confirm no regressions
- [x] 3.2 Run linter to confirm code quality
