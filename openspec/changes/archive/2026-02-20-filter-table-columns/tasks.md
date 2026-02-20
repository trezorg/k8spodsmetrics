## 1. Column Type and Validation

- [x] 1.1 Create `internal/columns/columns.go` with `Column` type and constants (total, allocatable, used, request, limit, available, free)
- [x] 1.2 Add `ValidForNodes()` and `ValidForPods()` validation functions
- [x] 1.3 Add `StringList()` helper for usage messages
- [x] 1.4 Add unit tests for column validation

## 2. Config File Support

- [x] 2.1 Add `Columns []string` field to `config.Common` struct
- [x] 2.2 Update `MergeCommon()` to merge columns from file config
- [x] 2.3 Update config package documentation with columns example

## 3. CLI Flag Integration

- [x] 3.1 Add `Columns []string` to `commonConfig` struct in `internal/adapters/stdin/`
- [x] 3.2 Add `--columns` StringSliceFlag to `commonFlags()` with validation action
- [x] 3.3 Wire columns from config file merge into commonConfig

## 4. Table Formatters

- [x] 4.1 Update `ToTable()` in node table formatter to accept columns parameter
- [x] 4.2 Update `ToTable()` in pod metrics table formatter to accept columns parameter
- [x] 4.3 Implement column filtering logic in node table `Print()` function
- [x] 4.4 Implement column filtering logic in pod metrics table `Print()` function
- [x] 4.5 Add unit tests for filtered table output

## 5. Integration

- [x] 5.1 Pass columns from CLI/config to table formatters in pods command
- [x] 5.2 Pass columns from CLI/config to table formatters in summary command
- [x] 5.3 Ensure columns parameter is empty/nil for non-table output formats

## 6. Verification

- [x] 6.1 Run `task check` to verify formatting, linting, and tests pass
- [ ] 6.2 Manual test: pods command with `--columns request,used`
- [ ] 6.3 Manual test: summary command with `--columns total,free`
- [ ] 6.4 Manual test: config file with columns setting
- [ ] 6.5 Manual test: invalid column shows helpful error
