# worklog

A CLI tool that logs daily work entries and generates standups, summaries, and more via Claude.

## Install

```bash
go install github.com/gbolanos-dev/worklog/cmd/worklog@latest
```

## Setup

Create `~/.worklog/config.yaml`:

```yaml
anthropic:
  api_key: "sk-ant-..."

# Optional: enable YouTrack and GitHub context
youtrack:
  base_url: "https://youtrack.example.com"
  token: "perm:..."

github:
  token: "ghp_..."
  default_repo: "owner/repo"
```

## Usage

```bash
# Log entries throughout the day
worklog add "refactored auth middleware"
worklog add --tag backend "added rate limiting to API"

# View today's entries
worklog list
worklog list --tag backend

# Generate a standup
worklog standup
worklog standup --issue PROJ-100 --pr 42

# Weekly summary
worklog summary --week
worklog summary --week --format promo

# Interactive chat with Claude using your work context
worklog chat
```

## Build from source

```bash
make build
./bin/worklog --help
```
