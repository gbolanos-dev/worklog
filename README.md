# worklog

A CLI tool that turns daily work log entries into polished standups via Claude.

## Install

```bash
go install github.com/gbolanos-dev/worklog/cmd/worklog@latest
```

## Setup

Create `~/.worklog/config.yaml`:

```yaml
anthropic:
  api_key: "sk-ant-..."
```

## Usage

```bash
worklog add "fixed timezone tests in DateUtilitiesTest"
worklog add "rebased SE-1173, resolved conflicts in PythonScriptEditor"
worklog list
worklog standup
```

## Build from source

```bash
make build
./bin/worklog --help
```
