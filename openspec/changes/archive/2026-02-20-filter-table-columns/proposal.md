## Why

Users viewing node or pod metrics in table format often want to focus on specific columns relevant to their task. Currently, all columns are always displayed, making tables wide and harder to read when only a subset of metrics is needed. Adding column filtering gives users control over table output density and focus.

## What Changes

- Add `--columns` CLI flag to filter displayed columns in table output format
- Support selecting one or multiple columns from the available set
- When `--columns` is not set, display all columns (current behavior preserved)
- Column filtering only applies to table output; other formats (json, yaml, string) remain unchanged

**Node columns (per resource type):**
- `total`, `allocatable`, `used`, `request`, `limit`, `available`, `free`

**Pod columns (per resource type):**
- `request`, `limit`, `used`

## Capabilities

### New Capabilities
- `table-column-filter`: CLI flag and validation for selecting which columns to display in table output

### Modified Capabilities
- None (purely additive feature; existing output behavior unchanged when flag not used)

## Impact

- **CLI**: New `--columns` flag in `internal/adapters/stdin/`
- **Choices**: New choice type for valid column names in `internal/choices/`
- **Table formatters**: Updated in `internal/adapters/stdout/table/{noderesources,metricsresources}/` to filter columns
- **Config**: YAML config file support for `columns` setting
