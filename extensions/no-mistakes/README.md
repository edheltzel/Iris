# No-mistakes gate pack

Refuses Iris task completion when the delivery mode requires
[no-mistakes](https://github.com/kunchenguid/no-mistakes) and the CLI is missing
or this workspace has no gate.

## Install

Copy the pack into the workspace extensions directory
(`.spynel/extensions` by default):

```bash
cp -R extensions/no-mistakes /path/to/workspace/.spynel/extensions/no-mistakes
```

`iris extension install` cannot install this pack from the Iris repository URL.
The manifest has to sit at the clone root.

## Mode

Write one Firstmate mode as the first line of `.spynel/no-mistakes-mode`.
A missing file means `no-mistakes`.

| Mode | Hook behavior |
| --- | --- |
| `no-mistakes` | Require the CLI and a repo gate. |
| `no-mistakes-prod-only` | Same. Iris cannot tell product work from internal work, so it requires the gate. |
| `direct-PR` | Journal a skip. Do not run the CLI. |
| `local-only` | Journal a skip. Do not run the CLI. |

## What it runs

`task.completed` for outcome `done` runs `no-mistakes status` in the workspace.
No `gate:` line, a failing status, or a missing CLI exits nonzero, and Iris
does not settle that task. Other outcomes skip.

`harness.after` runs the same check and journals it, then exits zero. Failing
that hook would replace the chat response, and the dispatch has already happened.

The pack never runs `no-mistakes axi run`. That command pushes and opens a PR,
and it blocks longer than a hook may run.

## Journal

Skip, fail, and pass lines go to stderr and to
`.spynel/extensions-state/no-mistakes/journal`.
