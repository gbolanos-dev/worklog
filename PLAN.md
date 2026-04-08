# worklog - Build Plan

> This file is the source of truth for the worklog project build plan.
> Claude Code should follow each commit step in order.
> After each commit is pushed, mark the item as complete by changing `[ ]` to `[x]`.

---

## How to Use This File

1. Read the current unchecked commit step
2. Implement only what that step describes — nothing more
3. Run any listed verification commands
4. Commit and push with the exact commit message shown
5. Mark the step `[x]` and stop — wait for next instruction

---

## MVP Commits

- [x] **Commit 1 — Init project scaffolding**

  Create the following structure from scratch:

  ```
  worklog/
  ├── main.go
  ├── plan.md              # this file, copied into the repo
  ├── .gitignore
  ├── go.mod               # go mod init github.com/<your-username>/worklog
  ├── cmd/
  │   ├── add.go
  │   ├── list.go
  │   └── standup.go
  ├── store/
  │   └── store.go
  └── config/
      └── config.go
  ```

  **File contents:**

  `main.go`
  ```go
  package main

  func main() {}
  ```

  `cmd/add.go`, `cmd/list.go`, `cmd/standup.go`
  ```go
  package cmd
  ```

  `store/store.go`
  ```go
  package store
  ```

  `config/config.go`
  ```go
  package config
  ```

  `.gitignore`
  ```
  worklog
  bin/
  *.yaml
  config.yaml
  ```

  **Verify:**
  ```bash
  go build ./...
  ```

  **Commit message:** `init: project structure and build plan`

  **Push:** `git push -u origin main`

---

- [x] **Commit 2 — Config loader**

  Implement `config/config.go` to load `~/.worklog/config.yaml` into a Go struct.

  - Define a `Config` struct with an `Anthropic.APIKey` field
  - Use `gopkg.in/yaml.v3` to unmarshal the file
  - Export a `Load() (*Config, error)` function
  - Create `~/.worklog/` directory if it does not exist
  - Add `gopkg.in/yaml.v3` via `go get gopkg.in/yaml.v3`

  **Example config file** (not committed, user creates manually):
  ```yaml
  anthropic:
    api_key: "sk-ant-..."
  ```

  **Verify:** temporarily add a `fmt.Println` in `main.go` that loads and prints
  the config, confirm it works, then remove it before committing.

  **Commit message:** `feat(config): load api key from ~/.worklog/config.yaml`

---

- [x] **Commit 3 — Store (read/write log entries)**

  Implement `store/store.go` to manage `~/.worklog/log.json`.

  - Define an `Entry` struct with `ID`, `Date`, and `Entry` string fields
  - Export `AddEntry(text string) error` — appends a new entry to the JSON file
  - Export `GetEntriesForDate(date string) ([]Entry, error)` — filters by date
  - Create `~/.worklog/` directory if it does not exist
  - Use `encoding/json` only, no external deps
  - Generate ID using `fmt.Sprintf("%x", time.Now().UnixNano())`

  **Commit message:** `feat(store): read/write log entries to ~/.worklog/log.json`

---

- [x] **Commit 4 — `add` command**

  Wire up `cmd/add.go` using the cobra CLI library.

  - Add `github.com/spf13/cobra` via `go get github.com/spf13/cobra`
  - Define an `AddCmd` cobra command that takes a single string argument
  - Call `store.AddEntry()` with the argument
  - Print confirmation: `Logged: <entry text>`
  - Set up root command in `main.go` and register `AddCmd` on it

  **Verify:**
  ```bash
  go run main.go add "fixed timezone tests"
  cat ~/.worklog/log.json
  ```

  **Commit message:** `feat(cmd): add command logs entry to store`

---

- [x] **Commit 5 — `list` command**

  Implement `cmd/list.go` to display today's entries.

  - Define a `ListCmd` cobra command
  - Call `store.GetEntriesForDate(today)` and print each entry
  - Format output as a numbered list
  - Register on root command

  **Verify:**
  ```bash
  go run main.go list
  ```

  **Commit message:** `feat(cmd): list command prints today's log entries`

---

- [x] **Commit 6 — Claude API client**

  Implement a minimal Claude API wrapper.

  - Create `claude/client.go` with package `claude`
  - Define a `Client` struct that holds the API key
  - Export `NewClient(apiKey string) *Client`
  - Export `Complete(prompt string) (string, error)` — POSTs to
    `https://api.anthropic.com/v1/messages` using `net/http` and `encoding/json`
  - Use model `claude-haiku-4-5-20251001`, max_tokens 1024
  - Set required headers: `x-api-key`, `anthropic-version: 2023-06-01`, `content-type`
  - No external HTTP libraries — stdlib only

  **Commit message:** `feat(claude): minimal API client wrapping /v1/messages`

---

- [x] **Commit 7 — `standup` command**

  Implement `cmd/standup.go` to generate a standup from today's entries.

  - Load config, load today's entries, build a prompt, call `claude.Complete()`
  - Prompt instructs Claude to write a standup with Yesterday, Today, and Blockers
  - Print Claude's response to stdout
  - Register on root command

  **Verify:**
  ```bash
  go run main.go add "fixed timezone tests in DateUtilitiesTest"
  go run main.go add "rebased SE-1173, resolved conflicts in PythonScriptEditor"
  go run main.go standup
  ```

  **Commit message:** `feat(cmd): standup command generates update via Claude`

---

## v2 Commits — Install & Update Infrastructure

- [x] **Commit 8 — refactor: extract main logic into internal/cli**

  Prepare for multiple binary entry points without duplicating logic.

  - Create `internal/cli/cli.go` with `Run(args []string) int`
  - Move root command setup and all command registration into `internal/cli`
  - `main.go` becomes a thin wrapper:
    ```go
    package main

    import (
        "os"
        "github.com/<you>/worklog/internal/cli"
    )

    func main() {
        os.Exit(cli.Run(os.Args[1:]))
    }
    ```
  - No behavior change — pure refactor

  **Verify:**
  ```bash
  go run main.go add "test refactor"
  go run main.go list
  go run main.go standup
  ```

  **Commit message:** `refactor: extract main logic into internal/cli`

---

- [x] **Commit 9 — feat: add cmd/worklog entry point and Makefile**

  Add a dedicated entry point so `go install` produces a binary named `worklog`.

  - Create `cmd/worklog/main.go`:
    ```go
    package main

    import (
        "os"
        "github.com/<you>/worklog/internal/cli"
    )

    func main() {
        os.Exit(cli.Run(os.Args[1:]))
    }
    ```
  - Add `Makefile`:
    ```makefile
    VERSION ?= dev
    LDFLAGS  = -ldflags "-X github.com/<you>/worklog/internal/cli.Version=$(VERSION)"

    .PHONY: build clean

    build:
    	go build $(LDFLAGS) -o bin/worklog ./cmd/worklog

    clean:
    	rm -rf bin/
    ```
  - Update `README.md` with install instructions:
    ```bash
    go install github.com/<you>/worklog/cmd/worklog@latest
    ```

  **Verify:**
  ```bash
  make build
  ./bin/worklog list
  ```

  **Commit message:** `feat: add cmd/worklog entry point and Makefile`

---

- [ ] **Commit 10 — feat: --version flag via ldflags**

  Wire a version string that gets stamped at build time from git tags.

  - Add `var Version = "dev"` to `internal/cli/cli.go`
  - Add `--version` flag to root command — prints `worklog <version>` and exits
  - Makefile already injects `VERSION` via ldflags from Commit 9

  **Verify:**
  ```bash
  make build
  ./bin/worklog --version
  # → worklog dev

  make VERSION=v0.1.0 build
  ./bin/worklog --version
  # → worklog v0.1.0
  ```

  **Commit message:** `feat: add --version flag with ldflags`

---

- [ ] **Commit 11 — feat: update check with GitHub releases cache**

  Background update detection with zero latency impact on commands.

  - Create `internal/update/version.go`:
    - `func newer(current, latest string) bool` — simple semver comparison
  - Create `internal/update/update.go`:
    - `type Result struct { Available bool; Latest string }`
    - `func Check(currentVersion, cacheDir string) *Result` — reads
      `~/.worklog/update-check.json` only, returns result if cache is fresh (24h TTL)
    - `func Refresh(cacheDir string)` — GET GitHub releases API with 3s timeout,
      writes cache file, fails silently on any network error
    - `func DoUpdate(module string) error` — runs `go install <module>@latest`
    - Cache path: `~/.worklog/update-check.json`
  - Create `internal/update/update_test.go` — test version comparison and cache read/write
  - Skip all update logic when `Version == "dev"`

  **Verify:**
  ```bash
  go test ./internal/update/...
  ```

  **Commit message:** `feat: add update check with GitHub releases cache`

---

- [ ] **Commit 12 — feat: wire update notification into CLI**

  Connect the update system to the CLI so every command run checks for updates.

  After any command completes in `cli.Run()` (skip for `--version`):

  1. Call `update.Check(Version, cacheDir)` — reads cache, instant
  2. If result available: print to stderr:
     ```
     Update available: vX.Y.Z  →  run: worklog --update
     ```
  3. If cache was stale: call `update.Refresh(cacheDir)` — fetch post-command
  4. Add `--update` flag: calls `DoUpdate()`, prints result, exits

  **Two-run model:** check now (from cache), notify next run, HTTP only after pipeline.

  **Verify:**
  ```bash
  # Tag and push a release on GitHub as v0.2.0
  # Delete ~/.worklog/update-check.json to simulate stale cache, then:
  ./bin/worklog list
  # → (entries printed)
  # → (cache refreshed silently in background)

  # Run again — cache now has v0.2.0:
  ./bin/worklog list
  # → (entries)
  # → Update available: v0.2.0  →  run: worklog --update

  ./bin/worklog --update
  # → Installing github.com/<you>/worklog/cmd/worklog@latest...
  # → Done. worklog v0.2.0 installed.

  worklog --version
  # → worklog v0.2.0
  ```

  **Commit message:** `feat: wire update notification and --update flag into CLI`

---

## v2 Commits — YouTrack & GitHub Context

- [ ] **Commit 13 — feat: YouTrack fetch**

  Pull live ticket context into prompts.

  - Create `fetch/youtrack.go` with package `fetch`
  - Define `YouTrackTicket` struct: `ID`, `Title`, `Description`, `MyComments []string`
  - Export `FetchTicket(baseURL, token, ticketID string) (*YouTrackTicket, error)`
  - Use YouTrack REST API:
    `GET /api/issues/{id}?fields=id,summary,description,comments(text,author(login))`
  - Filter comments to only those authored by the configured username
  - Add `youtrack.base_url`, `youtrack.token`, `youtrack.username` to config struct

  **Verify:**
  ```bash
  go run main.go standup --ticket IGN-15026
  ```

  **Commit message:** `feat(fetch): YouTrack ticket context`

---

- [ ] **Commit 14 — feat: GitHub PR fetch**

  Pull PR context into prompts.

  - Create `fetch/github.go` with package `fetch`
  - Define `GitHubPR` struct: `Number`, `Title`, `Body`, `FilesChanged []string`
  - Export `FetchPR(token, owner, repo string, number int) (*GitHubPR, error)`
  - Use GitHub REST API:
    - `GET /repos/{owner}/{repo}/pulls/{number}`
    - `GET /repos/{owner}/{repo}/pulls/{number}/files`
  - Add `github.token`, `github.default_repo` to config struct

  **Verify:**
  ```bash
  go run main.go standup --pr 4821
  ```

  **Commit message:** `feat(fetch): GitHub PR context`

---

- [ ] **Commit 15 — feat: --ticket and --pr flags on standup**

  Combine manual log entries with live ticket and PR context in the Claude prompt.

  - Add `--ticket` (repeatable) and `--pr` (repeatable) flags to `standup` command
  - Fetch all provided tickets and PRs before building the prompt
  - Create `claude/context.go` with `BuildPrompt(ctx WorkContext, task string) string`
  - `WorkContext` holds `LogEntries []store.Entry`, `Tickets []fetch.YouTrackTicket`,
    `PullRequests []fetch.GitHubPR`

  **Verify:**
  ```bash
  go run main.go standup --ticket SE-1173 --pr 4821
  ```

  **Commit message:** `feat(cmd): --ticket and --pr flags on standup`

---

- [ ] **Commit 16 — feat: worklog chat interactive mode**

  Load all context once then iterate on output interactively.

  - Add `chat` cobra command
  - Accept `--ticket` and `--pr` flags, fetch context upfront
  - Start a read-eval loop: read user input from stdin, append to conversation
    history, call Claude, print response
  - Type `exit` or `quit` to end the session
  - Maintain full conversation history across turns so Claude retains context

  **Verify:**
  ```bash
  go run main.go chat --ticket SE-1173 --pr 4821
  # > summarize what I did
  # > make it shorter
  # > write it as a YouTrack comment
  # > exit
  ```

  **Commit message:** `feat(cmd): worklog chat interactive refinement mode`

---

- [ ] **Commit 17 — feat: summary --week command**

  Weekly digest across all log entries.

  - Add `summary` cobra command with `--week` flag
  - Load all entries from the past 7 days via store
  - Build prompt asking Claude for a structured weekly summary
  - Support `--format=promo` flag for promotion-doc-friendly output

  **Verify:**
  ```bash
  go run main.go summary --week
  go run main.go summary --week --format=promo
  ```

  **Commit message:** `feat(cmd): summary --week generates weekly digest via Claude`

---

- [ ] **Commit 18 — feat: tag support**

  Tag entries at log time and filter by tag on any command.

  - Add `--tag` flag (repeatable) to `add` command
  - Add `Tags []string` field to `Entry` struct in store
  - Export `GetEntriesByTag(tag string) ([]Entry, error)` from store
  - Add `--tag` filter flag to `list`, `standup`, and `summary` commands

  **Verify:**
  ```bash
  go run main.go add --tag SE-1173 "wired hyperlink listener"
  go run main.go list --tag SE-1173
  go run main.go standup --tag SE-1173
  ```

  **Commit message:** `feat: tag support on add and filter commands`
