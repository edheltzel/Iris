# Extension packs DOX

## Purpose

- Own the first-party executable hook packs Iris ships for copying into a workspace's extensions directory.

## Local Contracts

- A pack here is ordinary extension content: a validated `.spynel-extension.yaml` manifest plus executables, discovered and trusted exactly like a third-party repository. Packs receive no privileged path into the application. `iris extension install` requires the manifest at a clone root, so a pack under this directory is copied in rather than installed by URL.
- Packs use only the hook contract `internal/extensions` already emits and the `SPYNEL_HOOK`, `SPYNEL_EXTENSION`, and `SPYNEL_WORKSPACE` environment. Never add a pack that requires a new Go capability without shipping that capability first.
- Hooks that integrate an external CLI exit zero when the tool, project mode, or payload field is absent. A nonzero exit fails the orchestration operation that delivered the event, so only a pack enforcing a deliberate gate may refuse.
- Packs consuming at-least-once events deduplicate visible effects persistently by `event_id` under the workspace's `.spynel/` state, never in memory.
- Each pack owns its executable tests beside it so `go test ./...` covers pack behavior with a stubbed external CLI.

## Child DOX Index

Direct child DOX files:

| Child | Scope |
| --- | --- |
| [gitbutler/AGENTS.md](gitbutler/AGENTS.md) | GitButler `but` task-lifecycle hook pack. |
| [no-mistakes/AGENTS.md](no-mistakes/AGENTS.md) | No-mistakes delivery-mode gate pack. |
| [quota/AGENTS.md](quota/AGENTS.md) | External workflow watcher and task-claim quota evidence. |
