`docs/iris-vision.md` is the living Iris scope SoT. This file is FM-554 companion provenance, not a second spec.

# FM-554 — Iris exceed-bar inventory

**Purpose:** Table-heavy inventory of capabilities Iris must match or beat (Firstmate + Sentinel), plus what Iris inherits day-one from Spynel.  
**Date:** 2026-09-17 (America/New_York)

---

## Lock reminder

**Dual track** (Sentinel #49 / FM-543): keep hacking Sentinel until captain is happy with Spynel; **destination = Spynel-base fork/adapt** named **Iris** in this factory (FM-554 / FM-553).

Locked destination prefs:

- Portable **hooks** early (gates / axi / GitButler CLI)
- **GitButler** instead of Treehouse
- **Herdr socket-first** (external watcher; not day-one Go panes)
- **Atlas stays sole Grok Bot front door** (no second liaison)
- Steal Firstmate/Sentinel goods into the fork as we go

---

## Master exceed table

| Capability | Firstmate | Sentinel | Spynel day-one (Iris inherits) | Iris exceed target (match/beat) | Notes / evidence |
| --- | --- | --- | --- | --- | --- |
| **no-mistakes / validation gate** | **Shipped.** Project modes `no-mistakes` / `direct-PR` / `local-only`; worker owns `no-mistakes axi`; gate refuse (`fm-gate-refuse-lib.sh`); `.no-mistakes.yaml` | **Not shipped.** Land = honor flag only. Idea **#45** (steal FM pipeline). Companion `no-mistakes` on doctor PATH list | **Partial.** Built-in Markdown review lifecycle (`todo→…→done`), `harness.reviews`, leases — **not** the `no-mistakes` CLI. Hooks can shell gate (FM-551) | **Beat:** ship-mode gate via hook pack early; modes match FM; refuse dispatch when CLI missing | FM: README Features, `docs/architecture.md` “No-mistakes gate”, `.no-mistakes.yaml`. Sentinel: README Land / #45. Spynel: `docs/tasks-and-goals.md` |
| **zero-token / quota watch** | **Shipped.** Zero-token bash watcher `bin/fm-watch.sh` (absorb benign wakes). `quota-axi` essential tool; dispatch profiles resolve arrays from quota output | **Partial.** Absorb = no model wake on idle ticks. `quota-axi` companion + doctor floor (≥0.1.29). No FM-depth zero-token fleet watcher documented as shipped | **Absent / different.** Semantic heartbeat = harness tokens when fired; no `quota-axi` found | **Match/beat:** zero-token absorb path + quota-aware dispatch (hook/companion or Go later) | FM: README “Event-driven, zero-token”; architecture Event-driven supervision; bootstrap `quota-axi`. Sentinel: README Absorb / Companion CLIs. Spynel: FM-543 snapshot |
| **wait contracts** | **Shipped.** Declared `paused:` external wait + verified `captain-held`; `FM_PAUSE_RESURFACE_SECS` recheck; status+backlog hold records | **Thin today.** Brief “Needs you”. Idea **#43**: named wait + optional `wake_at` + progress log + leases on **tasks-axi** (not Spynel markdown) | **Shipped.** Named waiting; optional RFC3339 `wake_at`; Progress log; claim leases / fencing | **Beat:** keep Spynel wait/`wake_at`/leases; port Sentinel axi habits if dual-store; tip transitions → notify | FM: architecture wait/pause sections. Sentinel: #43. Spynel: `docs/tasks-and-goals.md` Waiting |
| **notify-on-transition** | **Shipped (fleet-local).** Actionable status wakes, merge-outcome emitter, optional Relay public follow-ups; captain-facing via firstmate only | **Thin shipped.** Spawn/land/stop `ui.notify` this session (never `sendUserMessage`). Idea **#44**: tip Grok Bot on wait/land/stop via **#16**, journal send/skip, no model wake solely to notify | **Shipped.** Terminal / actionable-wait → notification-agent job; `notify.origin`; agent journals send/skip/fail in Progress | **Beat:** Spynel transition notify → **Grok Bot** wire (not TG/WA product); origin + journal; Atlas captain path | FM: architecture merge-outcome / Relay. Sentinel: README HIL updates; #44/#16. Spynel: `docs/integrations.md`, tasks-and-goals |
| **heartbeat / check-in** | **Shipped.** Watcher heartbeat backstop; `/bearings`, `/ahoy`; AFK return brief. No user `/check-in` | **Shipped (bash stamp).** `sentinel-heartbeat` metronome; `heartbeatAt`; alert if stale. Explicitly **no** `/check-in`. Away first-slice ≠ full FM AFK | **Shipped.** Primary election: 5s heartbeats / 30s stale; semantic heartbeat worker (ordinary async job); `/status` `/jobs` | **Match:** keep Spynel semantic + election heartbeats; add Herdr socket poke for external visibility; no second check-in UX | FM: architecture Event-driven. Sentinel: README Sidecar / Commands. Spynel: AGENTS.md primary election; tasks-and-goals heartbeat |
| **task board** | **Shipped.** `tasks-axi` + `data/backlog.md`; Bearings four-section; fleet snapshot JSON | **Shipped (companion).** `tasks-axi` required for spawn; `/sentinel brief` four-section (Needs you / Landed / In flight / Queued); WorkerPort metas | **Shipped.** Markdown task/goal OS SoT; `/tasks` `/goals` / CLI; leases; Progress | **Match/beat:** Spynel markdown board day-one; steal axi UX habits if needed; Hub-class visibility later | FM: AGENTS layout, `.tasks.toml`. Sentinel: README brief / Companion. Spynel: tasks-and-goals |
| **harness-agnostic dispatch** | **Shipped.** Co-primary: Claude/Grok/Pi (+ omp, Codex, OpenCode, Cursor…). `config/crew-dispatch.json` NL profiles → harness/model/effort | **omp-centric.** Default backend **omp-task**; `--herdr` opt-in. Claude = guidance-only (no live HubTool). Spawn blocked if `gh-axi`/`tasks-axi`/gh auth missing | **Shipped.** `harness.Harness` + catalog: Codex, Claude Code, Pi, ACP, Agent Zero CLI; channel≠harness | **Beat:** keep Spynel multi-harness; add dispatch profiles / quota arrays from FM; omp optional client not sole runtime | FM: README Recommended harnesses; architecture Dispatch profiles. Sentinel: README Workers. Spynel: harness-compatibility.md, AGENTS.md |
| **Herdr / mux visibility** | **Shipped.** `backend=herdr` (CI lane); presentation spaces; one tab/task; Treehouse worktrees; native busy/events | **Partial.** `--herdr` spawn; peek only (not Agent Hub). omp Hub for omp-task workers | **Absent.** Multi-TUI + `/jobs` status; no Herdr mux panes | **Beat via socket-first:** external Herdr watcher + hooks; Go panes only if socket fails (lock) | FM: `docs/herdr-backend.md`, architecture Runtime backends. Sentinel: README worker `--herdr`. Spynel/FM-551: weak without Go |
| **Grok Bot liaison surface** | **Shipped pattern.** Single captain liaison (`GROK_BOT.md`); Grok verified primary harness; crewmates never address captain | **Locked product:** Atlas = sole Grok Bot front door. Idea **#16** crew↔Sentinel talk (rooms/DMs); not shipped. No second liaison | **Different product.** Own TUI + Telegram/WhatsApp as “one human → one assistant”. Not Atlas | **Match lock:** Atlas remains sole captain door; Spynel TG/WA **not** product front door; optional Grok Bot as Go `channel.Channel` later or localapi client | FM: `GROK_BOT.md`, README. Sentinel: #16/#49; FM-543. Spynel: README / integrations; FM-551 Grok Bot channel = Go |
| **Worktree isolation (extra)** | **Shipped.** [Treehouse](https://github.com/kunchenguid/treehouse) pooled worktrees (tmux/herdr/zellij/cmux); Orca owns own | **Not found** as first-class Treehouse integration in README | **Absent** as Treehouse product; workspace `.spynel/` + harness sandboxes | **Beat:** **GitButler** (`but` CLI via hooks) instead of Treehouse (lock) | FM: architecture Worktrees. Iris lock / #49 / FM-543 |
| **Away / AFK supervision (extra)** | **Shipped.** `/afk` `/quiet`; away contract; sub-supervisor daemon; return brief; wedge alarm | **First slice shipped.** Persist away + live-session inbox auto-deliver; **not** full FM AFK (no unattended spawn). Planned: recap/stow/update (#22) | **Absent** as FM-style AFK product | **Match later:** optional away; do not block Iris MVP | FM: architecture Away mode. Sentinel: README Away / Bar |
| **Extension / hook surface (extra)** | Harness-specific extensions (`.pi/`, `.omp/`, Claude hooks, Grok hooks); scripts under `bin/` | omp plugin extension `extensions/sentinel.ts` + sidecar | **Shipped.** Portable executable hooks (5 names); install via git; fail-closed | **Inherit + use early** for gates/axi/GitButler | Spynel: `docs/extensions.md`; FM-551 |
| **Primary election / multi-TUI (extra)** | Session lock + secondmate homes; restart-proof state | Multi-session away owner lock; shared metronome | **Shipped.** Single primary election; secondaries over loopback API | **Inherit** as Spynel control-plane property | Spynel: AGENTS.md; integrations.md |

**Master table row count:** 13 (9 Atlas-named + 4 extras).

---

## Spynel day-one inheritance

| What | Comes free (Spynel upstream) | Needs hooks / companion CLI | Needs Go (fork/contrib) |
| --- | --- | --- | --- |
| Markdown task/goal OS, leases, `wake_at`, Progress | Yes | — | — |
| Review lifecycle / `harness.reviews` | Yes | — | — |
| Harness plugs (Codex/Claude/Pi/ACP/A0) | Yes | — | New harness adapter |
| Channels TUI / Telegram / WhatsApp | Yes (upstream) | — | **Product lock:** do not use as captain door |
| Primary election, 5s/30s, `/status` `/jobs` | Yes | — | — |
| Semantic heartbeat + notification-agent jobs | Yes | Wire tip target → Grok Bot | Optional `channel.Channel` for Grok Bot |
| Portable hooks (5 lifecycle names) | Yes | Hook pack: no-mistakes, axi, `but`, Herdr socket poke | — |
| Zero-token bash watcher / quota-axi arrays | No | Companion + scripts/hooks | Or native Go watcher later |
| Herdr mux panes | No | **Socket-first external watcher** (lock) | Deep panes only if needed |
| Treehouse worktrees | No | **GitButler** CLI via hooks (lock) | GB-native store later (open Q) |
| Firstmate AFK / secondmates | No | Steal selectively later | — |
| tasks-axi as SoT | No (markdown SoT) | Optional axi companion beside markdown | Replacing SoT = product change |

---

## Gap / steal list (prioritized for Intern → FM-553)

| P | Gap | Steal from | Land on Iris as | Why |
| --- | --- | --- | --- | --- |
| 1 | Ship validation rigor | FM no-mistakes modes + Sentinel #45 | Hook pack + companion CLI; refuse when missing | FM bar Iris must beat; Spynel review ≠ full gate |
| 2 | Zero-token absorb + quota-aware dispatch | FM `fm-watch` + `quota-axi` / crew-dispatch | External watcher or hooks + companion | Avoid token burn; match FM exceed bar |
| 3 | Notify tips to Grok Bot (not TG/WA) | Spynel transition notify + Sentinel #44/#16 | Notification path → Atlas/Grok Bot rooms/DMs; journal send/skip | Lock: Atlas sole door; Spynel notify packaging is stealable |
| 4 | Herdr visibility | FM herdr backend lessons; Sentinel `--herdr` peek | **Socket-first** external client | Lock: not day-one Go panes |
| 5 | GitButler worktrees | Replace FM Treehouse | Hooks shell `but`; GB store later | Lock: GitButler not Treehouse |
| 6 | Wait contracts polish | Spynel already strong; Sentinel #43 axi fields | Keep Spynel wait/`wake_at`/leases; mirror axi UX if dual-track Sentinel still needs it | Mostly inherited; steal field names/habits |
| 7 | Dispatch profiles | FM `crew-dispatch.json` | Config + harness catalog selection | Beat single-harness defaults |
| 8 | Away / AFK | FM full; Sentinel first-slice | Later; don’t block MVP | Nice-to-have |
| 9 | Grok Bot as Spynel channel | FM-551 / #16 | localapi client first; Go `channel.Channel` only if needed | Atlas remains door; avoid second liaison |

**Near-term fork order (aligned FM-543):** stand Iris fork → hook pack (gates / axi / `but`) → Herdr socket client → budget Go only where hooks can’t reach.

---

## Sources

### Firstmate (`kunchenguid/firstmate`)

| Path / URL | Used for |
| --- | --- |
| https://github.com/kunchenguid/firstmate/blob/main/README.md | Features, liaison, zero-token, backends, harnesses |
| https://github.com/kunchenguid/firstmate/blob/main/AGENTS.md | Supervisor contract, layout, tasks-axi, modes |
| https://github.com/kunchenguid/firstmate/blob/main/docs/architecture.md | Watcher, waits, Herdr, Treehouse, no-mistakes, dispatch, AFK |
| https://github.com/kunchenguid/firstmate/blob/main/GROK_BOT.md | Captain-only liaison pattern |
| https://github.com/kunchenguid/firstmate/blob/main/.no-mistakes.yaml | Gate config |
| `.agents/skills/bootstrap-diagnostics/SKILL.md`, `bin/fm-bootstrap.sh` | `gh-axi` / `tasks-axi` / `quota-axi` floors |

### Sentinel (`edheltzel/Sentinel`)

| Path / URL | Used for |
| --- | --- |
| `README.md` (sha `6404612…`) | Shipped concepts: heartbeat, absorb, deliver, workers, Herdr, HIL notify, companion CLIs |
| `AGENTS.md` | Coordinator bar, no `/check-in`, away limits |
| `docs/GETTING-STARTED.md` | First-run / deliver / away honesty |
| https://github.com/edheltzel/Sentinel/issues/49 | Dual-track / Iris destination lock |
| https://github.com/edheltzel/Sentinel/issues/43 | Wait contracts on tasks-axi (idea) |
| https://github.com/edheltzel/Sentinel/issues/44 | Notify-on-transition via Grok Bot (idea) |
| https://github.com/edheltzel/Sentinel/issues/45 | Steal no-mistakes (idea) |
| https://github.com/edheltzel/Sentinel/issues/16 | Grok Bot talk wire (idea) |

### Spynel (`agent0ai/spynel`)

| Path / URL | Used for |
| --- | --- |
| https://github.com/agent0ai/spynel/blob/main/README.md | Product pillars |
| https://github.com/agent0ai/spynel/blob/main/AGENTS.md | Primary election, leases, honesty, harness-neutral |
| https://github.com/agent0ai/spynel/blob/main/docs/extensions.md | Hook names / protocol |
| https://github.com/agent0ai/spynel/blob/main/docs/tasks-and-goals.md | Wait, wake_at, heartbeat, notify |
| https://github.com/agent0ai/spynel/blob/main/docs/integrations.md | Channels, notify agent, multi-TUI |
| https://github.com/agent0ai/spynel/blob/main/docs/harness-compatibility.md | Harness matrix |

### Prior factory reports

| Path | Used for |
| --- | --- |
| `/home/box/agent-data/grok-ship/reports/FM-551-spynel-extensibility.md` | Hooks vs Go; Grok Bot / Herdr hang points |
| `/home/box/agent-data/grok-ship/reports/FM-543-spynel-vs-sentinel.md` | Dual-track lock; steal order |

### Blockers / missing docs

- **FM-542** report file not present under `grok-ship/reports/` (only FM-543/FM-551 Spynel-related); Spynel facts taken from upstream + FM-543 relay.
- Sentinel is **private** to anonymous raw.githubusercontent; content via GitHub MCP (`user-GitHub-xai`).
- Firstmate `AGENTS.md` truncated by MCP size cap; capability claims cross-checked via README + architecture + search_code.
- Sentinel #43/#44/#45/#16 are **ideas** (not shipped) — marked accordingly.
- No Iris fork remote invented; Spynel origin remains `agent0ai/spynel`.

---

## Hand-off tip (Atlas → Intern)

Iris inherits Spynel’s durable task OS, waits/`wake_at`, leases, multi-harness plugs, portable hooks, primary election, and notification-agent packaging for free — but **does not** inherit Firstmate’s no-mistakes ship gate, zero-token bash watcher + quota-array dispatch, Treehouse worktrees, or deep Herdr panes. Exceed bar = hook pack early (no-mistakes + axi + GitButler `but`) + Herdr **socket-first** + Grok Bot notify tips while **Atlas stays the only captain door**; keep Sentinel live on dual track and fold this table into FM-553 without re-fetching.
