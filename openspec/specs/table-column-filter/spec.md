# table-column-filter Specification

## Purpose

Enable users to select which columns appear in table output for node and pod metrics, reducing table width and focusing on relevant data.

## Requirements

### Requirement: Column Selection via CLI Flag

The system **SHALL** provide a `--columns` CLI flag that accepts one or more column names to display in table output.

#### Scenario: Single column selection

- **WHEN** user runs command with `--columns used`
- **THEN** table output **SHALL** display only the "used" column for each resource type

#### Scenario: Multiple column selection

- **WHEN** user runs command with `--columns request,limit,used`
- **THEN** table output **SHALL** display only the specified columns in their defined order

#### Scenario: No columns specified

- **WHEN** user runs command without `--columns` flag
- **THEN** table output **SHALL** display all available columns (default behavior preserved)

### Requirement: Column Validation for Node Commands

The system **SHALL** validate column names against the valid set for node contexts.

#### Scenario: Valid node column

- **WHEN** user runs node/summary command with `--columns total,allocatable,free`
- **THEN** command **SHALL** succeed and display the specified columns

#### Scenario: Invalid node column

- **WHEN** user runs node/summary command with `--columns invalid`
- **THEN** command **SHALL** fail with error message listing valid columns: `total`, `allocatable`, `used`, `request`, `limit`, `available`, `free`

### Requirement: Column Validation for Pod Commands

The system **SHALL** validate column names against the valid set for pod contexts.

#### Scenario: Valid pod column

- **WHEN** user runs pods command with `--columns request,used`
- **THEN** command **SHALL** succeed and display the specified columns

#### Scenario: Invalid pod column

- **WHEN** user runs pods command with `--columns available`
- **THEN** command **SHALL** fail with error message listing valid columns: `request`, `limit`, `used`

### Requirement: Scope Limited to Table Output

Column filtering **SHALL** only affect table output format.

#### Scenario: Columns flag with JSON output

- **WHEN** user runs command with `--output json --columns used`
- **THEN** JSON output **SHALL** include all fields (column filtering ignored)

#### Scenario: Columns flag with YAML output

- **WHEN** user runs command with `--output yaml --columns used`
- **THEN** YAML output **SHALL** include all fields (column filtering ignored)

### Requirement: Config File Support

The system **SHALL** support column selection via YAML configuration file.

#### Scenario: Columns from config file

- **WHEN** config file contains `columns: [request, limit]`
- **AND** no `--columns` CLI flag is provided
- **THEN** table output **SHALL** display only request and limit columns

#### Scenario: CLI overrides config

- **WHEN** config file contains `columns: [request]`
- **AND** CLI provides `--columns used,limit`
- **THEN** table output **SHALL** display only used and limit columns (CLI takes precedence)
