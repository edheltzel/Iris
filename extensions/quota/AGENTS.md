# Quota companion DOX

## Purpose

- Own the optional heartbeat companion that absorbs a tick with no waiting task and attaches a quota snapshot when a waiting task exists.

## Local Contracts

- The script decides. Iris only honors a first-line `absorb`. Any other result, including a missing script, keeps the current provider heartbeat.
- Do not encode reminder thresholds, inactivity timers, or quota ranking here. A waiting task always dispatches so the heartbeat agent can decide.
- `quota-axi` is evidence on stdout. Never pick a harness or model from it.
- A failing or missing `quota-axi` still dispatches. Token savings never outrank a waiting task.
- Journal absorb and dispatch under `.spynel/extensions-state/quota/journal`.

## Child DOX Index

No child DOX files.
