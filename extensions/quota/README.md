# Quota companion

Copy this directory to `.spynel/extensions/quota`. It does not change the
semantic heartbeat.

## Watcher

`watch.sh` is an external poll. Iris does not start it.

```bash
SPYNEL_WORKSPACE=/path/to/workspace extensions/quota/watch.sh
```

An unchanged checksum of `.spynel/tasks` and `.spynel/goals` Markdown files
prints nothing. A changed file prints `wake: workflow files changed` and a
`quota-axi --no-credential-refresh` snapshot when that CLI is installed.
The script does not rank models or send a reminder.

## Claim hook

`task.claimed` writes `.spynel/extensions-state/quota/latest.txt` and adds one
progress line, `Quota evidence: <path> (not a route)`, to the claimed file.
The hook exits zero even when `quota-axi` is missing. It does not cancel the
claim.
