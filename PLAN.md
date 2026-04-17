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

- [x] **Commit 10 — feat: --version flag via ldflags**

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

- [x] **Commit 11 — feat: update check with GitHub releases cache**

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

- [x] **Commit 12 — feat: wire update notification into CLI**

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

- [x] **Commit 13 — feat: YouTrack fetch**

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

- [x] **Commit 14 — feat: GitHub PR fetch**

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

- [x] **Commit 15 — feat: --ticket and --pr flags on standup**

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

- [x] **Commit 16 — feat: worklog chat interactive mode**

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

- [x] **Commit 17 — feat: summary --week command**

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

- [x] **Commit 18 — feat: tag support**

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

---

## Phase 2A: Correctness and Queryable Logs

- [ ] **Commit 19 -- feat(store): add FilterByTag helper**

  Add `FilterByTag(entries []Entry, tag string) []Entry` to `store/store.go`. A pure function that filters an already-loaded slice instead of loading from disk. Do NOT remove `GetEntriesByTag` yet.

  **Go concepts:** Pure functions on slices, no I/O side effects.

  **Verify:** `go build ./...`

  **Commit message:** `feat(store): add FilterByTag helper for in-memory filtering`

---

- [ ] **Commit 20 -- fix: tag filter no longer replaces date-scoped entries**

  In `list.go`, `standup.go`, and `summary.go`, replace:
  ```go
  entries, err = store.GetEntriesByTag(tag)
  ```
  with:
  ```go
  entries = store.FilterByTag(entries, tag)
  ```

  Then remove `GetEntriesByTag` from `store/store.go` (now unused).

  **Go concepts:** Dead code removal. Composing filter functions.

  **Verify:**
  ```bash
  go build ./...
  go run main.go add --tag test "tagged entry"
  go run main.go list --tag test
  # Should only show today's entries with that tag
  ```

  **Commit message:** `fix(cmd): tag filter no longer replaces date-scoped entries`

---

- [ ] **Commit 21 -- fix: replace unused --week flag with --days on summary**

  Remove `var week bool` and `BoolVar(&week, "week", ...)` from `summary.go`. The command always shows the past 7 days regardless of the flag.

  Replace it with `--days N` (default 7). Use the value to compute the `since` date: `time.Now().AddDate(0, 0, -days)`. This gives the user actual control over the summary range instead of removing flexibility.

  Also update `README.md`: replace `worklog summary --week` with `worklog summary` and `worklog summary --days 14`. Update `worklog summary --week --format promo` to `worklog summary --format promo`.

  **Go concepts:** Flag hygiene -- replace dead flags with useful ones. `IntVar` for numeric flags. Keep docs and CLI surface in sync within a single commit.

  **Verify:**
  ```bash
  go run main.go summary --help    # --week should be gone, --days should appear
  go run main.go summary           # defaults to 7 days
  go run main.go summary --days 14 # two-week summary
  grep -- "--week" README.md       # should produce no output
  ```

  **Commit message:** `fix(cmd): replace unused --week flag with --days on summary`

---

- [ ] **Commit 22 -- feat: --date flag on list and standup**

  Add `--date` / `-d` flag to both `list` and `standup`. If provided, use that date instead of today. Validate format with `time.Parse("2006-01-02", date)`.

  For `list`, this shows entries for the specified date. For `standup`, this generates a standup from that date's entries — useful when you forgot to run it yesterday.

  **Go concepts:** `time.Parse` with Go's reference date, input validation at command boundary, applying the same flag pattern to multiple commands.

  **Verify:**
  ```bash
  go run main.go list --date 2026-04-15
  go run main.go list --date bad-date      # should error
  go run main.go standup --date 2026-04-15  # standup for a specific day
  go run main.go standup --date bad-date    # should error
  ```

  **Commit message:** `feat(cmd): add --date flag to list and standup commands`

---

- [ ] **Commit 23 -- feat: list --since and --until flags**

  Add `--since` and `--until` flags. Use `GetEntriesSince` then post-filter with a new `FilterUntil(entries []Entry, until string) ([]Entry, error)` in store.

  **Validation rules (enforce at the command boundary, before any store calls):**
  1. `--date` is mutually exclusive with BOTH `--since` and `--until` (use `cmd.MarkFlagsMutuallyExclusive` twice -- cobra groups are pairwise).
  2. `--until` requires `--since` to also be set. Using `--until` alone is an error.
  3. If both `--since` and `--until` are set, `--until` must be on or after `--since`. Reject with a clear error message.
  4. All date inputs must parse with `time.Parse("2006-01-02", ...)`. Reject malformed input up front.

  Default behavior (no flags): same as today -- show today's entries only.

  **Go concepts:** `time.Time.After()` / `time.Time.Before()`, `cmd.MarkFlagsMutuallyExclusive`, composing filter pipelines, `PreRunE` for validation that must happen before `RunE`.

  **Verify:**
  ```bash
  go run main.go list --since 2026-04-01
  go run main.go list --since 2026-04-01 --until 2026-04-10
  go run main.go list --date 2026-04-15 --since 2026-04-01   # should error (mutex)
  go run main.go list --date 2026-04-15 --until 2026-04-10   # should error (mutex)
  go run main.go list --until 2026-04-10                     # should error (until requires since)
  go run main.go list --since 2026-04-10 --until 2026-04-01  # should error (until before since)
  go run main.go list --since bad-date                       # should error (parse)
  ```

  **Commit message:** `feat(cmd): add --since and --until flags to list command`

---

## Phase 2B: Entry Lifecycle and CLI UX

- [ ] **Commit 24 -- feat: show entry IDs in list output**

  Edit and delete work by hex ID, so the ID must be visible to the user. Right now `list` shows only `N. entry text`, forcing users to `cat ~/.worklog/log.json` -- we ship this *before* the commands that depend on it.

  In `cmd/list.go`, update the default (non-JSON, once --json exists) output so each row includes a short, copyable form of the ID. Use the first 8 characters of `Entry.ID` to keep it scannable, and accept either the short or full ID in later `delete` / `edit` commands (handled in their commits).

  Suggested format (the exact layout is up to you, but ID must be visible and trivially copyable):
  ```
  1. [a1b2c3d4] fixed timezone tests in DateUtilitiesTest [backend]
  ```

  **Go concepts:** `string` slicing (`id[:8]`), UX-driven output design, shipping discoverability before the feature that depends on it.

  **Verify:**
  ```bash
  go run main.go add "id visibility test"
  go run main.go list
  # Output row should include an 8-char ID prefix that matches the entry in log.json
  ```

  **Commit message:** `feat(cmd): show short entry IDs in list output`

---

- [ ] **Commit 25 -- feat: delete command**

  Create `cmd/delete.go`. Takes one arg: entry ID (full hex or the short 8-char prefix from Commit 24). Add `DeleteEntry(id string) error` to store (load, find by prefix match, remove from slice, write back). If the prefix matches more than one entry, return an error asking the user to disambiguate with a longer prefix. Register in `cli.go`.

  **Confirmation by default:** Before deleting, print the matched entry text and prompt `Delete this entry? [y/N]` using `bufio.NewReader(os.Stdin)`. Only proceed on `y` or `Y`. Add a `--force` / `-f` flag to skip the prompt (useful for scripting).

  **Go concepts:** Deleting from a slice: `append(entries[:i], entries[i+1:]...)`. The `...` variadic unpacking operator. `strings.HasPrefix` for prefix matching. Defensive UX with confirmation prompts. `bufio.NewReader` for single-line stdin reads.

  **Verify:**
  ```bash
  go run main.go add "entry to delete"
  go run main.go list             # copy the 8-char ID from output
  go run main.go delete <id>      # should show entry and ask for confirmation
  go run main.go delete -f <id>   # should delete without asking
  go run main.go list             # entry should be gone
  ```

  **Commit message:** `feat(cmd): delete command removes entry by ID`

---

- [ ] **Commit 26 -- feat: edit command**

  Create `cmd/edit.go`. Takes two args: entry ID (full or short prefix) and new text. Add `EditEntry(id, newText string) error` to store. Same prefix-matching rules as `delete`: reject ambiguous matches with a clear error. Register in `cli.go`.

  **Go concepts:** Mutating structs in a slice -- must use `entries[i].Entry = newText`, not the range copy `e.Entry = newText` (it would silently no-op). Reuse of the same prefix-lookup pattern from Commit 25.

  **Verify:**
  ```bash
  go run main.go add "typo entry"
  go run main.go list               # copy short ID
  go run main.go edit <id> "corrected entry"
  go run main.go list               # should show corrected text
  ```

  **Commit message:** `feat(cmd): edit command updates entry text by ID`

---

- [ ] **Commit 27 -- feat: list --json output mode**

  Add `--json` bool flag to `list`. When true, output `json.MarshalIndent(entries, "", "  ")` instead of numbered text.

  **Go concepts:** `json.MarshalIndent`, struct tags controlling output, same data with different presentation.

  **Verify:**
  ```bash
  go run main.go list --json
  go run main.go list --json --tag work
  ```

  **Commit message:** `feat(cmd): add --json output mode to list command`

---

- [ ] **Commit 28 -- feat: init command for interactive config setup**

  Create `cmd/init.go`. Prompts for API keys/URLs via `bufio.NewReader(os.Stdin)`. Builds `config.Config`, marshals to YAML, writes to `~/.worklog/config.yaml`. Warns if file exists.

  **Go concepts:** `bufio.NewReader` vs `bufio.Scanner`, `yaml.Marshal` (inverse of Unmarshal), `os.Stat` for file existence.

  **Verify:**
  ```bash
  cp ~/.worklog/config.yaml ~/.worklog/config.yaml.bak
  go run main.go init
  cat ~/.worklog/config.yaml
  mv ~/.worklog/config.yaml.bak ~/.worklog/config.yaml
  ```

  **Commit message:** `feat(cmd): init command for interactive config setup`

---

- [ ] **Commit 29 -- feat: config validation**

  Add `Validate(sections ...string) error` method on `*Config`. Checks required fields are non-empty. Returns descriptive error with hint to run `worklog init`. Call it in standup, summary, chat before API calls.

  **Go concepts:** Methods on structs, variadic parameters, `fmt.Errorf`, early return validation pattern.

  **Verify:**
  ```bash
  # Temporarily clear api_key in config, then:
  go run main.go standup
  # Should print: config: anthropic.api_key is required (run "worklog init" to set up)
  ```

  **Commit message:** `feat(config): validate required fields before commands run`

---

- [ ] **Commit 30 -- feat: configurable model and max tokens**

  Add `anthropic.model` (default `claude-haiku-4-5-20251001`) and `anthropic.max_tokens` (default `1024`) to the config struct in `config/config.go`. Update `claude/client.go` to accept these values instead of hardcoding them — pass them through `NewClient` or add a `ClientOptions` struct.

  This needs to land before retro/reporting where longer outputs will hit the 1024 token ceiling.

  **Go concepts:** Struct defaults with zero-value handling — if `model` is empty string after YAML load, fall back to the default. Options pattern for configuring a client. Separating tuning config from required config.

  **Verify:**
  ```bash
  go build ./...
  # With no model/max_tokens in config.yaml, should still work (defaults apply)
  go run main.go standup
  # Add to config.yaml:
  #   anthropic:
  #     api_key: "sk-ant-..."
  #     model: "claude-sonnet-4-5-20241022"
  #     max_tokens: 2048
  go run main.go standup  # should use the configured model
  ```

  **Commit message:** `feat(config): configurable model and max tokens`

---

## Phase 2C: CLI Polish

- [ ] **Commit 31 -- feat: add lipgloss for styled output**

  Add `charmbracelet/lipgloss` for colored, styled terminal output across all commands.

  - `go get github.com/charmbracelet/lipgloss`
  - Create `internal/ui/styles.go` with package `ui` that defines reusable styles:
    - `Success` -- green text (for "Logged:", "Deleted:", "Updated:" confirmations)
    - `Warning` -- yellow text (for "No entries found", update notifications)
    - `Header` -- bold cyan (for section headers in stats, retro)
    - `Muted` -- gray (for IDs, dates, secondary info)
    - `Error` -- red bold (for error messages)
  - Update `cmd/add.go`, `cmd/list.go`, and `cmd/delete.go` (once it exists) to use styled output
  - Keep it restrained: color enhances, text should still be clear without ANSI

  **Go concepts:** Third-party library integration. Package-level variables for reusable config. The `lipgloss.NewStyle()` builder pattern (method chaining). Separating presentation from logic by centralizing styles in one package.

  **Verify:**
  ```bash
  go build ./...
  go run main.go add "styled entry"
  # "Logged:" should appear in green
  go run main.go list
  # Entry numbers or dates should have subtle color
  ```

  **Commit message:** `feat(ui): add lipgloss for styled terminal output`

---

- [ ] **Commit 32 -- feat: table-formatted list output**

  Replace the plain numbered list in `list` with a formatted table using lipgloss table.

  - Use `github.com/charmbracelet/lipgloss/table` (included with lipgloss)
  - Table columns: `#`, `Date`, `Entry`, `Tags`
  - Apply header styling using the styles from `internal/ui/styles.go`
  - Only use table format for the default (non-JSON) output
  - Also apply table formatting to `stats` output (once it exists)

  **Go concepts:** Working with lipgloss table API: `table.New()`, `table.Row()`, `table.Headers()`. Rendering styled tables to stdout. Conditional formatting (table vs JSON vs plain).

  **Verify:**
  ```bash
  go build ./...
  go run main.go add --tag work "table test"
  go run main.go list
  # Should display a nicely formatted table with columns
  go run main.go list --json
  # JSON output should remain unchanged
  ```

  **Commit message:** `feat(ui): table-formatted list output with lipgloss`

---

- [ ] **Commit 33 -- feat: spinner for API calls**

  Add a spinner/loading indicator while waiting for Claude API responses in standup, summary, retro, and chat.

  - `go get github.com/schollz/progressbar/v3`
  - `go get golang.org/x/term` (for TTY detection)
  - Create a helper in `internal/ui/spinner.go`:
    - `StartSpinner(label string) *progressbar.ProgressBar` -- creates and starts an indeterminate spinner
    - `StopSpinner(bar *progressbar.ProgressBar)` -- clears the spinner line
  - **TTY detection:** `StartSpinner` must check `term.IsTerminal(int(os.Stderr.Fd()))` first. If stderr is not a terminal (piped, redirected, CI, non-interactive), return `nil` and render nothing. `StopSpinner` must be a no-op when passed `nil`. This keeps scripted use clean -- no control characters in logs.
  - Wrap `client.Complete()` calls: start spinner before, stop after
  - For chat, show spinner during each turn's API call
  - Spinner writes to stderr so it doesn't pollute stdout when the command's stdout is piped

  **Go concepts:** Writing to `os.Stderr` vs `os.Stdout`. `golang.org/x/term.IsTerminal` for interactivity detection. Nil-safe helper APIs. Cleaning up terminal state (clearing the spinner line). The pattern of wrapping I/O calls with UI feedback without changing the underlying logic.

  **Verify:**
  ```bash
  go build ./...
  go run main.go standup
  # Should see a spinner while waiting for Claude's response

  go run main.go standup 2>/tmp/stderr.log
  # No spinner visible, and /tmp/stderr.log should contain no ANSI escape codes

  go run main.go standup | cat
  # Stdout is piped but stderr is still the terminal -- spinner should still render
  ```

  **Commit message:** `feat(ui): add spinner for API calls`

---

- [ ] **Commit 34 -- feat: use huh for interactive init command**

  Replace raw `bufio` prompts in the `init` command with `charmbracelet/huh` forms.

  - `go get github.com/charmbracelet/huh/v2`
  - Rewrite `cmd/init.go` to use `huh.NewForm()` with:
    - `huh.NewInput()` for text fields (API key, base URL, token, repo)
    - `huh.NewConfirm()` for overwrite confirmation
  - Group fields into logical sections: "Anthropic", "YouTrack (optional)", "GitHub (optional)"
  - Optional sections can be skipped

  **Go concepts:** Declarative UI with builder pattern. `huh.NewForm().WithGroups(...)`. Error handling from form submission. Comparing imperative I/O (`bufio` read loop) vs declarative forms (`huh` schema). This is a common pattern shift in modern Go tooling.

  **Verify:**
  ```bash
  go build ./...
  cp ~/.worklog/config.yaml ~/.worklog/config.yaml.bak
  go run main.go init
  # Should show a polished interactive form
  cat ~/.worklog/config.yaml
  mv ~/.worklog/config.yaml.bak ~/.worklog/config.yaml
  ```

  **Commit message:** `feat(cmd): use huh for interactive init prompts`

---

## Phase 2D: AI Workflows and Reporting

- [ ] **Commit 35 -- feat: retro command**

  Create `prompts/retro.md` (wins/blockers/improvements template). Embed it in `prompts/prompts.go`. Create `cmd/retro.go` with `--days` (default 14) and `--tag` flags.

  **Go concepts:** Reinforces `//go:embed`, the full AI workflow pattern, `strconv.Itoa`.

  **Verify:**
  ```bash
  go run main.go retro
  go run main.go retro --days 7 --tag backend
  ```

  **Commit message:** `feat(cmd): retro command generates AI retrospective`

---

- [ ] **Commit 36 -- feat: export command**

  Create `cmd/export.go`. Dumps entries as markdown grouped by date. Add `GetAllEntries() ([]Entry, error)` to store. Supports `--since`, `--until`, `--tag` flags. Outputs to stdout for piping.

  **Go concepts:** `strings.Builder` for efficient concatenation, `sort.Slice` for ordering, `map[string][]Entry` for grouping.

  **Verify:**
  ```bash
  go run main.go export
  go run main.go export --since 2026-04-01 --tag work
  go run main.go export > /tmp/log.md && cat /tmp/log.md
  ```

  **Commit message:** `feat(cmd): export command dumps entries as markdown`

---

- [ ] **Commit 37 -- feat: stats command**

  Create `cmd/stats.go`. Shows total entries, entries per day (with bar chart), tag breakdown sorted by count, most active day. Supports `--since` and `--until` flags.

  **Go concepts:** `map[string]int` for counting, sorting maps by value, `strings.Repeat` for bars, `fmt.Sprintf("%-15s", ...)` for aligned columns.

  **Verify:**
  ```bash
  go run main.go stats
  go run main.go stats --since 2026-04-01
  ```

  **Commit message:** `feat(cmd): stats command shows entry analytics`

---

## Phase 3: Time-Scale Reporting (Draft)

> This phase extends `summary` into a multi-scale reporting tool rather than adding a parallel `report` command. Start narrow: prove monthly works, then expand to yearly.

### 3A: Monthly and Yearly Summaries

- [ ] **Commit 38 -- feat: summary --days 30 with monthly prompt**

  Create `prompts/monthly.md` — a prompt template optimized for monthly-scale output: themes, project progress, key accomplishments, patterns. Embed it in `prompts/prompts.go`.

  In `summary.go`, select the prompt template based on `--days`:
  - 1–14 → `prompts.Summary` (weekly)
  - 15–60 → `prompts.Monthly`
  - 61+ → `prompts.Yearly` (once it exists, fall back to Monthly until then)

  **Verify:**
  ```bash
  go run main.go summary --days 30
  go run main.go summary --days 30 --tag backend
  ```

  **Commit message:** `feat(cmd): monthly summary prompt for longer time ranges`

---

- [ ] **Commit 39 -- feat: yearly summary prompt**

  Create `prompts/yearly.md` — annual-scale template: major projects, growth areas, impact highlights. Designed to feed into self-reviews and promo docs. Embed it.

  Update the threshold in `summary.go`: `--days 61+` selects `prompts.Yearly`.

  **Verify:**
  ```bash
  go run main.go summary --days 365
  go run main.go summary --days 90 --format promo
  ```

  **Commit message:** `feat(cmd): yearly summary prompt for annual reviews`

---

### 3B: Report Saving (Explicit)

- [ ] **Commit 40 -- feat: summary --save flag**

  Add `--save` flag to `summary`. When passed, write the generated output to `~/.worklog/reports/` as a timestamped markdown file with metadata header (date range, tags, model used). Without `--save`, behavior is unchanged (stdout only).

  Add `worklog report list` subcommand to list saved reports. Add `worklog report show <filename>` to display one.

  **Verify:**
  ```bash
  go run main.go summary --days 30 --save
  go run main.go report list
  go run main.go report show <filename>
  ```

  **Commit message:** `feat(cmd): save summary output to reports directory`
