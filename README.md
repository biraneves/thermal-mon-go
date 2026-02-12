# thermal-mon-go

A lightweight CPU temperature monitor written in Go for Linux systems using `hwmon`.

It periodically reads a thermal sensor file, evaluates temperature against warning/critical thresholds, and prints color-coded status in the terminal.

## Features

- Reads CPU temperature from Linux `hwmon` sensor files.
- Converts millidegrees (`42000`) to Celsius (`42.80°C`).
- Configurable thresholds and polling interval via CLI flags.
- Colorized status output (`OK`, `WARNING`, `CRITICAL`).
- Fast-fail validation for invalid threshold configuration.
- Small codebase with unit tests for core thermal logic.

## Why this project

`thermal-mon-go` is focused on being:

- simple to run;
- easy to understand;
- easy to extend for CLI monitoring workflows.

It is intentionally minimal and avoids unnecessary dependencies.

## Requirements

- Linux (with accessible thermal sensor files under `/sys/class/hwmon/...`)
- Go `1.25.6` (as defined in `go.mod`)

## Installation

### Option 1: Build from source

```bash
git clone https://github.com/biraneves/thermal-mon-go.git
cd thermal-mon-go
go build -o thermal-mon ./cmd
```

### Option 2: Run directly without building

```bash
go run ./cmd
```

## Usage

```bash
thermal-mon [-w warning] [-c critical] [-i interval] [-z thermal_zone_path]
```

### Flags

- `-w` `float`: Warning threshold in Celsius. Default: `75.0`
- `-c` `float`: Critical threshold in Celsius. Default: `85.0`
- `-i` `duration`: Polling interval. Default: `30s`. Examples: `10s`, `1m`
- `-z` `string`: Thermal sensor file path. Default: `/sys/class/hwmon/hwmon1/temp1_input`

### Examples

Run with defaults:

```bash
go run ./cmd
```

Custom thresholds and interval:

```bash
go run ./cmd -w 70 -c 80 -i 10s
```

Custom sensor path:

```bash
go run ./cmd -z /sys/class/hwmon/hwmon0/temp1_input
```

## Output

Typical output format:

```
[2026-02-12 15:04:05] OK: Current temperature: 42.10°C
[2026-02-12 15:04:35] WARNING: Current temperature: 77.30°C
[2026-02-12 15:05:05] CRITICAL: Current temperature: 88.20°C
```

Status colors are applied in terminal output:

- `OK`: green
- `WARNING`: yellow
- `CRITICAL`: red

## Error handling behavior

- Invalid CLI input (for example, malformed flag values) exits with code `2` and prints a clean usage message.
- Invalid threshold relationship (`warning >= critical`) fails fast at startup and exits with code `2`.
- Read/parsing errors for thermal files are logged to `stderr` during runtime, and monitoring continues.

## Project structure

```
cmd/
  main.go        # Composition root
  helpers.go     # CLI parsing, startup output, monitor loop helpers

internal/
  thermal/
    helpers.go       # Thermal file reading + domain threshold logic
    helpers_test.go  # Unit tests for thermal logic
  colors/
    colors.go        # ANSI color constants
```

## Development

Run tests:

```bash
go test ./...
```

Format and vet:

```bash
go fmt ./...
go vet ./...
```

## Design notes

The project keeps domain decisions (temperature status) separate from terminal presentation (ANSI colors), improving maintainability and making it easier to adapt to other interfaces later.

## Contributing

Contributions are welcome.

Suggested workflow:

1. Fork the repository
1. Create a feature branch from `main`
1. Make focused commits (prefer Conventional Commits)
1. Add/update tests when behavior changes
1. Open a pull request with context and rationale

## Commit convention

This projects benefits from [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/), e.g.:

- `feat(cli): add custom interval flag`
- `refactor(thermal): decouple status from color output`
- `test(thermal): add missing file parsing scenarios`

## License

- [MIT License](/LICENSE)