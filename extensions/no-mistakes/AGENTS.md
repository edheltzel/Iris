# No-mistakes gate pack DOX

## Purpose

- Own the hook pack that refuses a completed task when Firstmate's delivery mode requires the no-mistakes CLI and that CLI or its repo gate is absent.

## Local Contracts

- Modes are the Firstmate closed set: `no-mistakes`, `direct-PR`, `local-only`, and `no-mistakes-prod-only`. Read the first line of `.spynel/no-mistakes-mode`. A missing or blank file means `no-mistakes`.
- `direct-PR` and `local-only` journal a skip and exit zero. They do not invoke the CLI.
- `no-mistakes-prod-only` has no product-versus-internal classifier here, so it requires the gate.
- An unknown mode fails `task.completed`.
- Only outcome `done`, or a missing outcome, is a ship. `failed`, `cancelled`, and `waiting` skip.
- Refuse by exiting nonzero from `task.completed` only. `harness.after` journals the same failure and exits zero, because failing that hook replaces the chat response and cannot refuse a dispatch that already happened.
- Run `no-mistakes status` from the workspace. Missing CLI, nonzero status, or no `gate:` line is a failure. Never run `no-mistakes axi run` from this hook: it publishes and exceeds the hook timeout.
- Journal skip, fail, and pass to stderr and `.spynel/extensions-state/no-mistakes/journal`.

## Child DOX Index

No child DOX files.
