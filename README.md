# SSH Task Runner (ssht)

A lightweight tool for executing commands on multiple SSH hosts.

[简体中文](./README_zh.md) | English

## Features

- Execute commands on multiple SSH hosts
- Node filtering with `--nodes` parameter
- Configurable logging format (text/json)
- Debug mode with `--debug` flag

## Installation

### From source
```bash
go build -o ssht main.go
```

### Using Docker
```bash
docker build -t ssht .
```

## Usage

1. Configure your hosts in `config.toml`
2. Run commands:

```bash
# Basic usage
./ssht --command "hostname"

# Run on specific nodes
./ssht --command "hostname" --nodes node1,node2

# Debug mode with JSON logging
./ssht --command "hostname" --nodes node1 --debug --log-format json

# Docker usage
docker run -v $(pwd)/config.toml:/app/config.toml ssht --command "hostname"
```

## Configuration

Edit `config.toml` to configure your SSH hosts and groups.

## Logging

Available log formats:
- `text` (default, colored output)
- `json` (structured logging)

Logging options:
- `--log-format`: specify the format (text/json)
- `--log-file`: write logs to specified file (default: stdout)

Examples:
```bash
# Write JSON logs to file
./ssht --command "hostname" --log-format json --log-file logs.json

# Write text logs to file
./ssht --command "hostname" --log-file logs.txt
```
