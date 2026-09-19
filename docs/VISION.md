This in-repo file is the living source of truth for Iris scope, spec, and plan.

# FM-553 — Iris vision

**SoT:** `docs/VISION.md`  
**Companion provenance:** [`docs/exceed-inventory.md`](exceed-inventory.md) (FM-554)  
**Index:** [`docs/README.md`](README.md)

## Amend history

| ID | Change |
| --- | --- |
| FM-553 | Origin vision / scope / spec SoT for the Iris conversion. |
| FM-555 | Fold the FM-554 exceed inventory into this document so FM-553 remains the single competing spec. |
| FM-603 | Land this living SoT in-repo as `docs/VISION.md`. Handoffs use the repo copy, not a ticket export. |

## Authority

- This file is the authoritative **scope / spec / plan** for Iris conversion work and handoffs.
- [`docs/README.md`](README.md) indexes what is authoritative versus companion.
- [`docs/exceed-inventory.md`](exceed-inventory.md) is FM-554 provenance only. It is not a second SoT. FM-555 already folded that inventory here.
- [`docs/vision.md`](vision.md) is public product positioning for the current Spynel-named application. It is not the conversion plan. Do not rewrite it as part of a docs-only SoT landing.
- Root [`AGENTS.md`](../AGENTS.md) and the DOX chain remain the binding engineering contract for current code. This file does not weaken DOX.

## Handoff contract

1. Read [`docs/README.md`](README.md), then this file, before planning or implementing Iris conversion work.
2. Treat later tickets, chat, and scratch notes as proposals until they are amended into this file.
3. Smallest logical change. Subtract before add.
4. Do not invent capabilities, integrations, guarantees, or adoption claims.
5. Docs-only landings (including FM-603) do not ship product code, rename, or Phase 1 implementation.

## What Iris is

This repository is **Iris**: a personal fork of Spynel for one-chat orchestration.

Public tagline on the GitHub repo: **One chat, unlimited AI orchestration.**

The inherited product model (still described in [`docs/vision.md`](vision.md) and root `AGENTS.md` until a later identity/rename ship) is:

- A classic, non-AI orchestration program. External coding harnesses supply intelligence.
- **One human → one agent → infinite agents.** “One agent” is the single assistant-facing relationship, not a claim that Iris is itself an AI agent. “Infinite agents” is scalable leverage, not a capacity guarantee.
- Three pillars: communication interface, Markdown task/goal management, agentic loops over external harnesses.
- Slogans remain **Simplicity at scale.** and **Simplicity. Leverage. Quality.** until a later copy ship amends them here and in public docs together.

Iris conversion means this fork becomes the Iris product surface without turning the orchestration core into an agent, and without competing with Codex, Claude Code, Pi, ACP, or future harnesses.

## Conversion scope

In scope for the Iris conversion program (later tickets; not this docs-only landing):

- Keep the orchestration model: deterministic communication, durable Markdown work, harness-neutral execution, fail-closed channels, single primary election.
- Make repository identity, docs, and later user-visible copy match Iris, when a dedicated identity/rename ship is opened.
- Use this file as the plan SoT so handoffs do not re-litigate scope in chat.

Out of scope until this file is amended:

- Recreating harness features inside Iris.
- Treating conceptual scale (“infinite agents”) as a resource, performance, or availability guarantee.
- Inventing integrations, platforms, or guarantees not already in the current product contract.
- Using FM-554 or any other companion as a competing spec.

## Exceed (folded from FM-554 via FM-555)

Iris should exceed stock Spynel as a personal one-chat orchestration surface. The exceed inventory originated in FM-554 and was folded here by FM-555.

Authoritative interpretation of “exceed” lives in this file. The companion [`docs/exceed-inventory.md`](exceed-inventory.md) is provenance only.

Exceed means: smaller, clearer, and more leverage for one human talking to one assistant who coordinates many harness sessions — not a larger surface area, not a second agent product, and not a rewrite of harness intelligence.

Concrete exceed items are added or removed by amending this file. Do not accumulate a parallel list.

## Non-goals for FM-603

FM-603 is a thin docs-only landing:

- No product code.
- No rename ship.
- No Phase 1 implementation.

RedTeam reviews the branch before a pull request is raised.

## How later ships should use this file

- Start from this SoT and the nearest DOX chain.
- If a change would alter scope, identity, exceed items, or sequencing, amend this file in the same change (or a docs-first predecessor).
- If a change is local engineering under the current Spynel-named contracts, follow root `AGENTS.md` / DOX and do not wait on conversion work.
- Public copy stays aligned with [`docs/vision.md`](vision.md) until an identity/rename ship updates both this SoT and public positioning together.
