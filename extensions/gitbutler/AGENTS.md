# GitButler hook pack DOX

## Purpose

- Own the hook pack that shells the GitButler CLI (`but`) on Iris task lifecycle events.

## Local Contracts

- Integrate GitButler only by shelling the `but` CLI. Never read or write GitButler's own store, and never add a Go dependency on it.
- `task.claimed` creates the task's branch when it is absent; `task.completed` commits the workspace onto that branch once per `event_id`. Branch names are `iris/<task-file-stem>` so both events resolve the same branch from the payload's `file` field.
- Never push, land, open a review, or rewrite history from a hook. Publishing stays a human action.
- Every path exits zero and reports refusals on stderr: absent `SPYNEL_WORKSPACE`, uninstalled `but`, an unreachable workspace, a repository outside GitButler mode, a payload without `file`, and a failing `but` subcommand.
- Record a completion receipt only after `but commit` succeeds, so a failed commit stays retryable and a redelivered event does not commit twice.
- Payload parsing is narrow by contract: single-line JSON string fields without escaped quotes. Do not grow it into a general JSON parser or add a `jq`/`python3` prerequisite.

## Child DOX Index

No child DOX files.
