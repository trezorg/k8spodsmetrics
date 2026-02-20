# Proposal: Consolidate Default Separator

## Why

Five packages duplicate the same `defaultSeparator = "|"` constant and the associated `StringList`/`StringListDefault` wrapper pattern. The core generic functions already exist in `internal/choices/`, so consolidating the default separator eliminates unnecessary repetition and ensures consistency.

## What Changes

- Add `DefaultSeparator` constant to `internal/choices/`
- Remove local `defaultSeparator` from 5 packages
- Update `StringListDefault()` calls to use `choices.DefaultSeparator`

## Capabilities

### Modified Capabilities
- `choices`: Exports `DefaultSeparator` constant for shared use

## Impact

- `internal/choices/choices.go`: Add exported constant
- `internal/sorting/metricsresources/metricsresources.go`: Use shared constant
- `internal/sorting/noderesources/noderesources.go`: Use shared constant
- `internal/alert/formats.go`: Use shared constant
- `internal/resources/resource.go`: Use shared constant
- `internal/output/formats.go`: Use shared constant
