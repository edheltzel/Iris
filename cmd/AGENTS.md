# Command DOX

## Purpose

- Own thin executable composition for the `iris` binary.

## Local Contracts

- `iris` with no subcommand launches the TUI and all enabled background services.
- `iris serve` runs enabled server channels and orchestration without requiring a terminal UI.
- Plain `iris send`, `followup`, `notify`, `status`, `conversations`, and `command` entry points remain thin adapters over `internal/cli` and the elected application service; they do not duplicate application behavior here. `notify` queues bounded input through the durable outbox and never starts a harness.
- `iris docs` remains a thin offline adapter over the curated `internal/agentdocs` catalog and must never initialize a workspace, contact a primary/local API, or start a harness.
- Command handlers delegate to `internal/cli`; do not put application logic in `main.go`.
- Preserve the npm launcher's dedicated update exit code without printing it as an application error; npm package replacement and restart remain launcher responsibilities.

## Child DOX Index

Direct child DOX files:

| Child | Scope |
| --- | --- |
| [iris/AGENTS.md](iris/AGENTS.md) | Executable process composition. |
