<!-- Managed by agent: keep sections and order; edit content, not structure. Last updated: 2026-03-21 -->

# AGENTS.md — cmd/minectl (CLI Commands)

<!-- AGENTS-GENERATED:START overview -->
## Overview
Cobra-based CLI entry point for minectl. Each file defines one subcommand. All commands delegate to `internal/` packages for business logic.
<!-- AGENTS-GENERATED:END overview -->

<!-- AGENTS-GENERATED:START filemap -->
## Key Files
| File | Purpose |
|------|---------|
| `minectl.go` | Root command, global flags (`--headless`, `--verbose`, `--log-encoding`), version check, subcommand registration in `init()` |
| `create.go` | `minectl create` — provisions a Minecraft server from manifest |
| `delete.go` | `minectl delete` — tears down a server by ID |
| `update.go` | `minectl update` — updates server edition/config |
| `list.go` | `minectl list` — lists servers for a given provider/region |
| `rcon.go` | `minectl rcon` — opens interactive RCON console |
| `plugins.go` | `minectl plugins` — uploads plugins to a running server |
| `wizard.go` | `minectl wizard` — interactive TUI wizard for server creation |
<!-- AGENTS-GENERATED:END filemap -->

<!-- AGENTS-GENERATED:START golden-samples -->
## Golden Samples (follow these patterns)
| For | Reference | Key patterns |
|-----|-----------|--------------|
| Standard CRUD command | `create.go` | `RunE` + `RunFunc` wrapper, flag setup, provisioner factory |
| Root + global flags | `minectl.go` | `PersistentPreRunE` for logging/UI setup, `PersistentPostRunE` for version check |
| Interactive command | `wizard.go` | Uses `internal/ui/form.go` for TUI prompts |
<!-- AGENTS-GENERATED:END golden-samples -->

<!-- AGENTS-GENERATED:START setup -->
## Setup & environment
- **CLI framework:** Cobra (`github.com/spf13/cobra`)
- **Build output:** `bin/minectl` (linux), `bin/minectl-darwin`, `bin/minectl.exe` (windows)
- Version injected via ldflags: `main.version`, `main.commit`, `main.date`
<!-- AGENTS-GENERATED:END setup -->

<!-- AGENTS-GENERATED:START commands -->
## Build & tests
- Build: `go build .` or `make build` (cross-platform)
- Run: `go run .`
- Test: `go test -v ./...` (tests are in `internal/`, not here)
- Lint: `golangci-lint run --timeout 10m -E goimports --fix`
<!-- AGENTS-GENERATED:END commands -->

<!-- AGENTS-GENERATED:START code-style -->
## Code style & conventions
- Use Cobra `RunE` (not `Run`) for commands that can fail
- Wrap `RunE` with `RunFunc()` from `minectl.go` for consistent error handling + exit codes
- Common flag sets: `--filename` (manifest), `--id` (server ID), `--ssh-key` (private key path)
- `createUpdatePluginProvisioner()` is the shared factory for commands needing a provisioner
- Register new commands in `init()` inside `minectl.go`
- Provide `--help` text for all commands and flags
- Exit codes: 0 = success, 1 = error (via `os.Exit` in `RunFunc`)
- Errors: displayed via `minectlUI.ErrorMsg()` (stderr in headless, TUI in interactive)
- `--headless` flag: enables structured logging, disables TUI rendering
<!-- AGENTS-GENERATED:END code-style -->

<!-- AGENTS-GENERATED:START security -->
## Security & safety
- Cloud credentials come from environment variables only — never from flags or config files
- Never log or display cloud tokens
- Validate manifest paths before reading
- SSH key paths: validate existence, do not log contents
<!-- AGENTS-GENERATED:END security -->

<!-- AGENTS-GENERATED:START checklist -->
## PR/commit checklist
- [ ] `--help` text is clear and accurate for new/modified commands
- [ ] New command registered in `minectl.go:init()`
- [ ] Uses `RunFunc` wrapper for error handling
- [ ] Flags have meaningful defaults and descriptions
- [ ] Works in both interactive and `--headless` mode
- [ ] Errors go through `minectlUI.ErrorMsg()`
<!-- AGENTS-GENERATED:END checklist -->

<!-- AGENTS-GENERATED:START examples -->
## Patterns to Follow
> **Prefer looking at real code in this repo over generic examples.**
> See **Golden Samples** section above for files that demonstrate correct patterns.
<!-- AGENTS-GENERATED:END examples -->

<!-- AGENTS-GENERATED:START help -->
## When stuck
- Cobra docs: see `github.com/spf13/cobra` README
- Check existing commands (especially `create.go`) for patterns
- Root AGENTS.md for project-wide conventions
<!-- AGENTS-GENERATED:END help -->
