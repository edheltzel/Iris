# GitButler hook pack

Shells the [GitButler](https://docs.gitbutler.com/cli-overview) CLI (`but`) on
Iris task lifecycle events, so orchestrated work lands on its own branch without
anyone running GitButler by hand.

## Install

Copy the pack into the workspace's extensions directory
(`.spynel/extensions` by default):

```bash
cp -R extensions/gitbutler /path/to/workspace/.spynel/extensions/gitbutler
```

`iris extension install` clones a repository and requires
`.spynel-extension.yaml` at the clone root, so it cannot install this pack from
the Iris repository URL. It works once the pack is a repository of its own.

`iris doctor` then reports `extensions: ok (gitbutler)`.

## Behavior

| Hook | `but` command | Effect |
| --- | --- | --- |
| `task.claimed` | `but branch new iris/<task-file-stem>` | Creates the task's branch, skipped when it already exists. |
| `task.completed` | `but commit --branch iris/<task-file-stem> -m "iris(<route>): <stem> <outcome>"` | Commits the workspace onto that branch, once per `event_id`. |

The branch name comes from the task Markdown file's basename, so both events
resolve the same branch.

## Prerequisites

- `but` on `PATH`.
- The workspace is a GitButler project (`but setup` once, then `but status`
  succeeds).

Neither is checked at install time. When either is missing the hook reports the
reason on stderr and exits zero.

## Limits

- **Advisory, never a gate.** Every path exits zero: a nonzero hook exit would
  fail the orchestration operation that delivered the event.
- **`but commit` commits the whole workspace.** Uncommitted changes unrelated to
  the task are included. `but undo` reverses a commit the hook made.
- **Publishing is yours.** The pack never runs `but push`, `but land`, or
  `but pr`.
- **Task events only.** `harness.before` and `harness.after` fire on every
  harness turn, which is the wrong granularity for commits; they are not wired.
- **Redelivery is deduplicated by `event_id`** in
  `.spynel/extensions-state/gitbutler/completed-events`. A commit that fails
  leaves the event retryable.
- **Narrow payload parsing.** Single-line JSON string fields, no escaped quotes,
  no `jq` or `python3` prerequisite.
