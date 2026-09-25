# Quota companion DOX

## Purpose

- Own the external workflow watcher and the task-claim hook that record quota evidence. Neither one changes the semantic heartbeat.

## Local Contracts

- `watch.sh` is not an Iris hook and Iris does not arm it. An unchanged checksum of `.spynel/tasks` and `.spynel/goals` Markdown files prints nothing. A change prints one wake line plus a `quota-axi --no-credential-refresh` snapshot, or a missing-CLI line.
- `hook.sh` runs only on `task.claimed`. It always exits zero. It writes `.spynel/extensions-state/quota/latest.txt` and adds one `Quota evidence:` progress line to the claimed file. A second claim does not add another line.
- Do not rank models, pick a harness, send a reminder, or cancel a claim. `quota-axi` output is evidence, not a route.
- Do not absorb or skip the semantic heartbeat from this pack.

## Child DOX Index

No child DOX files.
