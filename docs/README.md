# Documentation

Use this page to choose the shortest path to the information you need. The repository README stays focused on what Spynel is; operational and implementation detail lives here.

## Scope and living SoT

Handoffs start here, then read the authoritative file.

- **Authoritative** — [VISION.md](VISION.md) is the living source of truth for Iris scope, spec, and plan (FM-553, landed by FM-603).
- **Companion** — [exceed-inventory.md](exceed-inventory.md) is FM-554 provenance. VISION.md already folded that inventory (FM-555). It is not a second spec.
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

- [Iris conversion VISION](VISION.md) — living scope/spec/plan SoT for handoffs.
- [Exceed inventory (companion)](exceed-inventory.md) — FM-554 provenance; not a second spec.
- [Product vision](vision.md) — positioning, product boundaries, and the three pillars.
- [Architecture](architecture.md) — application, process, history, harness, channel, and orchestration design.
- [Coding harness compatibility](harness-compatibility.md) — evidence-backed lifecycle coverage and known gaps.
- [Provider-canary threat model](provider-canary-threat-model.md) — controls for authenticated provider testing.
- [Releasing](releasing.md) — native packaging, npm publication, credentials, and release verification.
