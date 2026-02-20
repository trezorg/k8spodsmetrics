# Spec: Shared Default Separator

## ADDED Requirements

### Requirement: Shared Default Separator

The `choices` package **MUST** provide a shared `DefaultSeparator` constant that all choice-based packages use for formatting option lists.

#### Scenario: Constant is exported and accessible

- **WHEN** a package imports `internal/choices`
- **THEN** `choices.DefaultSeparator` **SHALL** equal `"|"`

#### Scenario: Existing behavior is preserved

- **WHEN** `StringListDefault()` is called on any choice type
- **THEN** the output **MUST** use `"|"` as separator
- **AND** the output format **SHALL** be identical to before the change
