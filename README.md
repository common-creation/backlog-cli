# Backlog CLI

A command-line interface for Backlog, built with Go.

## Features

- Issue management (list, get, create)
- Project listing
- Simple configuration

## Installation

```bash
go install github.com/common-creation/backlog-cli/cmd/backlog@latest
```

## Usage

### Configuration

Before using the CLI, you need to configure it with your Backlog space URL and API key:

```bash
backlog config init --space "https://yourspace.backlog.com" --api-key "your-api-key"
```

### List Projects

```bash
backlog project list
```

### List Issues

```bash
backlog issue list
backlog issue list --project PROJECT_KEY
backlog issue list --status STATUS_ID
backlog issue list --count 50
```

### Get Issue Details

```bash
backlog issue get --key PROJECT-123
```

### Create Issue

```bash
backlog issue create --project PROJECT_KEY --summary "Issue summary" --description "Issue description" --issue-type 1 --priority 3
```

## Development

### Requirements

- Go 1.24 or later

### Building from Source

```bash
git clone https://github.com/common-creation/backlog-cli.git
cd backlog-cli
go build -o backlog ./cmd/backlog
```

## License

MIT
