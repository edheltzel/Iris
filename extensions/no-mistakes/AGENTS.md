# No-mistakes gate pack DOX

## Purpose

- Own the hook pack that refuses a ship when the resolved Firstmate task mode is `no-mistakes` and the `no-mistakes` CLI is not installed.

## Local Contracts

- Task modes are only `no-mistakes`, `direct-PR`, and `local-only`. Read the first line of `.spynel/no-mistakes-mode`. A missing or blank file means `no-mistakes`.
- `no-mistakes-prod-only` is a Firstmate registry policy, not a task mode. Refuse it until intake resolves it to `no-mistakes` or `direct-PR`.
- `direct-PR` and `local-only` journal a skip and exit zero. They do not look for the CLI.
- Only outcome `done`, or a missing outcome, is a ship. `failed`, `cancelled`, and `waiting` skip.
- When the mode is `no-mistakes` and the CLI is missing, both `task.completed` and `harness.after` exit nonzero.
- When the CLI is present, journal that `/no-mistakes` is the post-commit agent step and exit zero. Do not run `no-mistakes`, `no-mistakes status`, `no-mistakes doctor`, or `no-mistakes axi run`. A status `gate:` line is not validation.
- Journal skip, fail, and pass to stderr and `.spynel/extensions-state/no-mistakes/journal`.

## Child DOX Index

No child DOX files.
