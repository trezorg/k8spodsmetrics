## Context

The application displays Kubernetes pod and node metrics in table format with fixed columns. Users need the ability to filter which columns appear in the output to focus on relevant metrics and reduce table width.

**Current state:**
- Table output always shows all columns for each resource type
- Existing patterns: `Output`, `Sorting`, `Alert` types with validation via `internal/choices`
- CLI uses `urfave/cli/v2` with `StringSliceFlag` for multi-value flags
- Config file supports `StringOrSlice` for YAML flexibility

## Goals / Non-Goals

**Goals:**
- Add `--columns` CLI flag accepting one or more column names
- Validate column names against valid sets per context (nodes vs pods)
- Preserve current behavior when flag not set (show all columns)
- Support column filtering in YAML config file
- Only affect table output format

**Non-Goals:**
- Column reordering (columns display in fixed order)
- Affecting json/yaml/text output formats
- Per-resource-type column selection (same columns for CPU/Memory/Storage)

## Decisions

### D1: Single columns type with context-aware validation

**Decision:** Create a single `Column` type in `internal/columns/columns.go` with separate valid sets for node and pod contexts.

**Rationale:** Columns differ between nodes (7 columns) and pods (3 columns). Using a single type with context-aware validation avoids duplicating types while maintaining strict validation.

**Alternatives considered:**
- Separate `NodeColumn` and `PodColumn` types: More type safety but duplicates validation logic and flag handling
- Single unified column set: Would allow invalid combinations like pods showing "available"

### D2: Empty columns = show all (default behavior)

**Decision:** When `--columns` is not provided or empty, display all columns (current behavior).

**Rationale:** Backward compatible; users only opt-in to filtering when needed.

### D3: Columns stored as `[]Column` in common config

**Decision:** Add `Columns []string` to `Common` config struct, merging from file like other slice values.

**Rationale:** Consistent with existing `StringOrSlice` pattern for namespaces and resources.

### D4: Pass columns to table formatters as parameter

**Decision:** Modify `ToTable()` functions to accept a `columns.Columns` parameter alongside `resources.Resources`.

**Rationale:** Mirrors existing pattern for resource filtering; keeps column logic encapsulated in formatters.

## Risks / Trade-offs

**Invalid column for context** → Validation error with helpful message listing valid columns for the current command

**Column requested but resource type not selected** → Show column as empty/zero (consistent with current behavior when data unavailable)

**Performance** → No impact; filtering happens at render time with simple set lookup
