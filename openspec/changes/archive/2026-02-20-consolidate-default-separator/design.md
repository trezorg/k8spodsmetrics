# Design: Consolidate Default Separator

## Context

Five packages define an identical `defaultSeparator string = "|"` constant and wrap the generic `choices.StringList` function. The `internal/choices` package already provides the generic `StringList[T ~string]` function used by all of them.

## Goals / Non-Goals

**Goals:**
- Eliminate duplicate constant definitions
- Maintain 100% behavioral compatibility
- Keep existing tests passing without modification

**Non-Goals:**
- Refactoring the `StringList` or `StringListDefault` wrapper functions themselves
- Changing the separator value from `"|"` to something else

## Decisions

### Decision 1: Export constant from choices package

Add `DefaultSeparator` as an exported constant in `internal/choices/choices.go` rather than creating a new package or leaving it unexported.

**Rationale:** The `choices` package is already the canonical location for shared choice-related utilities. Adding the constant here follows the existing pattern.

### Decision 2: Keep wrapper functions in each package

Each package will continue to have its own `StringList()` and `StringListDefault()` wrapper functions—they just reference the shared constant instead of a local one.

**Rationale:** This preserves the API of each package and avoids breaking any callers. The wrappers are thin and appropriate for each type's domain.
