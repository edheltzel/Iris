# No-mistakes gate pack

Refuses an Iris ship when the task mode is `no-mistakes` and the
[no-mistakes](https://github.com/kunchenguid/no-mistakes) CLI is not installed.
Passing this hook does not mean the change was validated.

## Install

Copy the pack into the workspace extensions directory
(`.spynel/extensions` by default):

```bash
cp -R extensions/no-mistakes /path/to/workspace/.spynel/extensions/no-mistakes
```

`iris extension install` cannot install this pack from the Iris repository URL.

## Mode

Write one task mode as the first line of `.spynel/no-mistakes-mode`.
A missing file means `no-mistakes`.

| Mode | Behavior |
| --- | --- |
| `no-mistakes` | CLI missing: exit nonzero. CLI present: journal the handoff and exit zero. |
| `direct-PR` | Skip. Do not look for the CLI. |
| `local-only` | Skip. Do not look for the CLI. |

`no-mistakes-prod-only` is a Firstmate registry policy, not a task mode. The hook refuses it. Resolve the task to `no-mistakes` or `direct-PR` first.

`failed`, `cancelled`, and `waiting` outcomes skip.

## What it does not run

The hook never executes `no-mistakes`, including `status`, `doctor`, and `axi run`.
`status` only shows that `no-mistakes init` ran. `axi run` publishes and blocks
longer than a hook may run.

When the CLI is present, the journal says `/no-mistakes` is still the
post-commit agent step. That step is what validates the change.
