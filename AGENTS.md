<!-- FOR AI AGENTS - Human readability is a side effect, not a goal -->
<!-- Managed by agent: keep sections and order; edit content, not structure -->
<!-- Last updated: 2026-03-21 | Last verified: 2026-03-21 -->

# AGENTS.md

**Precedence:** the **closest `AGENTS.md`** to the files you're changing wins. Root holds global defaults only.

## Commands (verified 2026-03-21)
> Source: Makefile, go.mod

<!-- AGENTS-GENERATED:START commands -->
| Task | Command | ~Time |
|------|---------|-------|
| Typecheck | `go build -v ./...` | ~15s |
| Lint | `golangci-lint run --timeout 10m -E goimports --fix` | ~30s |
| Format | `gofmt -w .` | ~5s |
| Test (single) | `go test -v -race ./internal/ui/...` | ~2s |
| Test (all) | `go test -v ./...` | ~10s |
| Build (local) | `go build .` | ~10s |
| Build (cross) | `make build` | ~30s |
| Install | `go install .` | ~10s |
| Coverage | `make coverage` | ~15s |
<!-- AGENTS-GENERATED:END commands -->

> If commands fail, verify against `Makefile` or ask user to update.

## Workflow
1. **Before coding**: Read nearest `AGENTS.md` for the area you're touching
2. **After each change**: Run the smallest relevant check (lint → typecheck → single test)
3. **Before committing**: Run full test suite if changes affect >2 files or touch shared code
4. **Before claiming done**: Run verification and **show output as evidence**

## File Map
<!-- AGENTS-GENERATED:START filemap -->
```
main.go                          → entrypoint, injects version/commit/date
cmd/minectl/                     → CLI commands (cobra): create, delete, list, update, wizard, plugins, rcon
internal/
  logging/                       → zap-based structured logging
  manifest/                      → YAML manifest parsing and JSON schema validation
  provisioner/                   → cloud provider orchestration via minectl-sdk
  rcon/                          → Minecraft RCON interactive prompt
  ui/                            → TUI components (spinner, table, forms via bubbletea/huh)
config/                          → server config templates per Minecraft edition
docs/                            → documentation and images
Dockerfile                       → chainguard static image, copies binary
Makefile                         → build, test, lint, coverage targets
```
<!-- AGENTS-GENERATED:END filemap -->

## Golden Samples (follow these patterns)
<!-- AGENTS-GENERATED:START golden-samples -->
| For | Reference | Key patterns |
|-----|-----------|--------------|
| CLI command | `cmd/minectl/create.go` | Cobra RunE, flag parsing, provisioner creation |
| Provisioner | `internal/provisioner/provisioner.go` | Interface, cloud provider switch, spinner lifecycle |
| UI component | `internal/ui/spinner.go` | Bubbletea model, headless vs interactive |
| Test pattern | `internal/ui/ui_test.go` | Table-driven tests, headless/interactive variants |
<!-- AGENTS-GENERATED:END golden-samples -->

## Utilities (check before creating new)
<!-- AGENTS-GENERATED:START utilities -->
| Need | Use | Location |
|------|-----|----------|
| Cloud provisioning | `minectl-sdk` | `github.com/dirien/minectl-sdk` (external) |
| CLI framework | `cobra` | `github.com/spf13/cobra` |
| TUI forms | `huh` | `github.com/charmbracelet/huh` |
| TUI spinners/tables | `bubbletea` + `bubbles` | `internal/ui/` |
| Logging | `zap` | `internal/logging/` |
| YAML parsing | `sigs.k8s.io/yaml` | `internal/manifest/` |
| JSON schema validation | `gojsonschema` | `internal/manifest/` |
| Error wrapping | `pkg/errors` | Used throughout |
<!-- AGENTS-GENERATED:END utilities -->

## Heuristics (quick decisions)
<!-- AGENTS-GENERATED:START heuristics -->
| When | Do |
|------|-----|
| Adding cloud provider | Add case in `internal/provisioner/provisioner.go:getProvisioner()`, implement in `minectl-sdk` |
| Adding CLI command | New file in `cmd/minectl/`, register in `minectl.go:init()` |
| Adding server edition | New config dir under `config/`, update manifest parsing |
| Version injection | Via ldflags in Makefile: `-X main.version=...` |
| Headless mode | Use `--headless` flag; switches UI to structured logging |
| Committing | Use Conventional Commits (`feat:`, `fix:`, `docs:`, etc.) |
| Adding dependency | Ask first — we minimize deps |
| Unsure about pattern | Check Golden Samples above |
<!-- AGENTS-GENERATED:END heuristics -->

## Repository Settings
<!-- AGENTS-GENERATED:START repo-settings -->
- **Default branch:** `main`
- **Go version:** 1.24.9
- **Module:** `github.com/dirien/minectl`
- **Docker:** chainguard static base image
- **CI:** GitHub Actions (lock-inactive workflow)
<!-- AGENTS-GENERATED:END repo-settings -->

## Key Decisions
<!-- AGENTS-GENERATED:START key-decisions -->
- **External SDK:** Core cloud logic lives in `github.com/dirien/minectl-sdk`, not in this repo. This repo is the CLI wrapper.
- **Charmbracelet TUI:** Interactive mode uses bubbletea/huh/lipgloss; headless mode falls back to zap logging.
- **Multi-cloud:** 15 cloud providers supported via provider-specific env vars (see `provisioner.go`).
- **Manifest-driven:** Server config is YAML manifests validated against JSON schema.
<!-- AGENTS-GENERATED:END key-decisions -->

## Boundaries

### Always Do
- Run pre-commit checks before committing
- Add tests for new code paths
- Use conventional commit format: `type(scope): subject`
- Follow Go 1.24 conventions and idioms
- Show test output as evidence before claiming work is complete

### Ask First
- Adding new dependencies
- Modifying CI/CD configuration
- Changing public API signatures or CLI flags
- Changes to `minectl-sdk` (external dependency)
- Repo-wide refactoring or rewrites

### Never Do
- Commit secrets, credentials, or cloud provider tokens
- Modify vendor/ or generated files
- Push directly to main branch
- Use `ioutil.*` (deprecated)
- Commit go.sum without go.mod changes

## Contributing (for AI agents)
- **Comprehension**: Understand the problem before submitting code. Read the linked issue, understand *why* the change is needed.
- **Context**: Every PR must explain trade-offs and link to the issue it addresses.
- **Continuity**: Respond to review feedback. Drive-by PRs without follow-up will be closed.

<!-- AGENTS-GENERATED:START module-boundaries -->
## Module Boundaries
- `cmd/minectl/` → CLI layer only; no business logic, delegates to `internal/`
- `internal/provisioner/` → orchestrates `minectl-sdk`; no direct cloud API calls
- `internal/ui/` → presentation only; no business logic
- `internal/manifest/` → parsing/validation only; no side effects
- `internal/logging/` → logging setup; used by `ui`
- `internal/rcon/` → RCON protocol; standalone
<!-- AGENTS-GENERATED:END module-boundaries -->

## Terminology
| Term | Means |
|------|-------|
| minectl | This CLI tool for provisioning Minecraft servers |
| minectl-sdk | External Go SDK containing cloud provider implementations |
| manifest | YAML file describing desired server configuration |
| edition | Minecraft variant: java, bedrock, fabric, forge, papermc, etc. |
| RCON | Remote Console protocol for Minecraft server administration |
| headless | Non-interactive mode for CI/automation (structured log output) |

## Index of scoped AGENTS.md
<!-- AGENTS-GENERATED:START scope-index -->
| Directory | Focus |
|-----------|-------|
| [`cmd/minectl/`](./cmd/minectl/AGENTS.md) | CLI commands, cobra patterns, flag handling |
<!-- AGENTS-GENERATED:END scope-index -->

> **Agents**: When you read or edit files in a listed directory, you **must** load its AGENTS.md first.

## When instructions conflict
The nearest `AGENTS.md` wins. Explicit user prompts override files.
- For Go-specific patterns, defer to language idioms and standard library conventions
