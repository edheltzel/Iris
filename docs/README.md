# Documentation

Use this page to choose the shortest path to the information you need. The repository README stays focused on what Spynel is; operational and implementation detail lives here.

## Scope and living SoT

Handoffs start here, then read the authoritative file.

- **Authoritative** — [iris-vision.md](iris-vision.md) is the full FM-553 living Iris scope/spec (amend history, phases, locks). Not a stub. `box` grok-ship/reports is no longer parallel law.
- **Companion** — [exceed-inventory.md](exceed-inventory.md) is the full FM-554 inventory for provenance. `iris-vision.md` already folded that inventory (FM-555). It is not a second spec.
- **Public positioning** — [vision.md](vision.md) remains the current product-positioning doc until a later identity/rename ship. Do not treat it as the conversion plan.

## Start and operate

- [Getting started and development](getting-started.md) — install, initialize a workspace, run from source, and verify a development checkout.
- [Configuration](configuration.md) — workspace, harness, interface, channel, speech, startup, orchestration, and extension settings.
- [Communication integrations](integrations.md) — TUI, Telegram, WhatsApp, voice, histories, and transport behavior.
- [TUI editing and terminal checks](tui-editing.md) — selection modes, clipboard fallback, word editing, and a local acceptance checklist.
- [Troubleshooting](troubleshooting.md) — common installation, harness, channel, startup, update, and speech problems.

## Coordinate and automate

- [Tasks and goals](tasks-and-goals.md) — durable Markdown workflows, review, recovery, waiting, and notifications.
- [Plain CLI and automation](cli.md) — messages, follow-ups, conversations, status, framework commands, and output contracts.
- [Agent-readable documentation](agent-docs.md) — the offline `spynel docs` interface and its versioned JSON schema.
- [Persistent per-agent instructions](persistent-instructions.md) — workspace-level preferences, role mapping, and precedence.
- [Extensions and hooks](extensions.md) — trusted executable extensions and delivery guarantees.

## Understand and maintain

- [Iris conversion vision](iris-vision.md) — full FM-553 living scope/spec/plan SoT (phases, locks, amend history).
- [Exceed inventory (companion)](exceed-inventory.md) — full FM-554 inventory; provenance only, not a second spec.
- [Product vision](vision.md) — positioning, product boundaries, and the three pillars.
- [Architecture](architecture.md) — application, process, history, harness, channel, and orchestration design.
- [Coding harness compatibility](harness-compatibility.md) — evidence-backed lifecycle coverage and known gaps.
- [Provider-canary threat model](provider-canary-threat-model.md) — controls for authenticated provider testing.
- [Releasing](releasing.md) — native packaging, npm publication, credentials, and release verification.
