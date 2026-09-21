This in-repo file (`docs/iris-vision.md`) is the living Iris scope/spec SoT. `box` grok-ship/reports is no longer parallel law.

# FM-553 — Iris conversion VISION / spec

**Kind:** Product vision + conversion spec (report only — no code, no PR)  
**Date:** 2026-09-17 (America/New_York)  
**Amended:** FM-555 exceed inventory; **FM-558** blinking-eyeball UX; **FM-559** fork home `edheltzel/Iris`; **FM-560** Phase 1 rename CLI/binary/package `spynel` → `iris`; **FM-561** workspace `.spynel/` retained + OS state migrated to `$HOME/.agents/Iris`; **FM-562** hard cut no `spynel` binary/alias; **FM-563** Go module/imports → iris; **FM-564** Phase 2 starts with GitButler `but`; **FM-565** then no-mistakes; **FM-566** full Phase 2 order; **FM-567** TG/WA optional side channels; **FM-568** BigMac `~/Developer/Iris`; **FM-569** npm → iris; **FM-570** dual-track Sentinel; **FM-571** cutover parked + npm install; **FM-572** OMP first-run docs; **FM-573** `@edheltzel/iris`; **FM-574** Iris room `f57aa087`; **FM-575** Jev intent; **FM-576** Phase 2 step (4) Jev/Typesafe routing; **Brainstorm-0919** Jev seats triage+notify → Sentinel#50
**Upstream:** FM-543 (dual-track), FM-551 (extensibility), FM-554 (exceed inventory), Sentinel [#49](https://github.com/edheltzel/Sentinel/issues/49)  
**Checkout (BigMac):** `/Users/ed/Developer/Iris` (**FM-568** — renamed from Spynel). Remotes: `origin=edheltzel/Iris`, `upstream=agent0ai/spynel` (push disabled), `old-spynel=edheltzel/spynel`  
**Rename:** destination product = **Iris** (Spynel-base fork/adapt). Upstream remains Spynel until fork exists.
**Fork home (locked FM-559 / grill 2026-09-18):** [`edheltzel/Iris`](https://github.com/edheltzel/Iris) — **Atlas creates/retargets**; Intern does not invent remotes.

---

## 1) Product one-liner + non-goals

### One-liner

**Iris** is a harness-agnostic wake hub and agent control plane: durable tasks, wait contracts, leases, and notify — with **Atlas as the only Grok Bot front door**, portable **hooks** for gates/CLIs, **GitButler** for ship VCS, and **Herdr** as a socket-first watcher — built by adapting Spynel so Iris **matches and exceeds** Firstmate + Sentinel overall ability.

### What Iris is

- A **Go-primary daemon** (Spynel shape) that owns task lifecycle, harness dispatch, and inspectable wait/land state.
- An **extension surface** of trusted executable hooks early (gates, axi, `but`, no-mistakes) without waiting on a full plugin ABI.
- The **long-term destination** of the Ideas dual-track lock (**FM-570**): **keep hacking Sentinel** (**no freeze**, **no sunset date**). **Cutover criteria parked (FM-571)** — not invented yet.

### What Iris is not

- A second captain / liaison beside Atlas.
- Telegram/WhatsApp as the **captain door** or required poke path (**FM-567:** optional **side channels** only; Grok Bot via **Atlas** is SoT).
- An omp-only plugin that stays forever tied to one harness seat.
- A day-one deep Herdr pane Hub in Go (socket-first watcher first).
- Treehouse as default VCS (GitButler instead).
- A hard kill of Sentinel on day one of the fork.

### Non-goals (Phase 0–1)

| Non-goal | Why |
| --- | --- |
| Invent / create Iris GitHub remote | **Locked:** `edheltzel/Iris` (FM-559). **Atlas** creates/retargets — Intern does not |
| Replace Atlas with Iris chat packaging | No second liaison |
| GB-native task store as Phase 0 | Hook-shell `but` first; native store is a later spike |
| Full Go Herdr pane ownership day one | Socket + hooks first; Go only if that fails |
| Freeze / sunset Sentinel on a date | **Locked (FM-570):** keep hacking Sentinel — no freeze, no sunset date yet |
| Invent cutover criteria now | **Parked (FM-571)** — unpark later |
| Day-one install = Go binary only | **Locked (FM-571):** **npm primary**, Go binary **secondary** |
| Iris = omp-only plugin | **No** — harness-agnostic; **OMP-biased onboarding only (FM-572)** |
| Invent Jev/Typesafe API details in VISION | **No (FM-575)** — intent only until phase + real docs |
| Absorb-only into omp Sentinel as the product lock | Superseded by dual track / Iris destination |
| Upstream Spynel TG/WA as **required** / captain door | **Locked (FM-567):** TG/WA remain **optional side channels** only; **Atlas sole Grok Bot captain door** — never a second liaison |
| Firstmate-full AFK as Iris MVP | Match later; don’t block cutover bar |
| Second mascot / channel face for Iris | Marks-only for channel avatar; blinking eyeball is in-product motion only (FM-558) |
| Soft `spynel` binary/alias after Phase 1 | **Hard cut (FM-562)** — `iris` only |

---

---

## Dual track (FM-570 / FM-571)

**Keep hacking Sentinel** while Iris ships (**no freeze**, **no sunset date** yet).

- **Cutover criteria: parked (FM-571)** — do not invent a Phase 2 “done → cutover” checklist until captain unparks it.
- Dual track continues; happy-with-Iris / Phase 2 progress still inform the later cutover conversation.

---

---

## Day-one install (FM-571)

| Channel | Role |
| --- | --- |
| **npm** | **Primary** day-one install path — publish as **`@edheltzel/iris`** (FM-573) |
| **Go binary** | **Secondary** (supported, not the lead) |

Do not invent a third primary installer for v1.

---

---

## First-run docs / examples (FM-572 — confirmed FM-573)

| Rule | Lock |
| --- | --- |
| **Product** | Harness-**agnostic** (keep Spynel multi-harness registry) |
| **Onboarding bias** | First-run docs and examples **default to OMP** |
| **Other harnesses** | Remain **supported** (Claude Code, Pi, Codex, ACP, …) — documented after OMP path, not dropped |

Do not rewrite Iris as an omp-only plugin. Do not lead README/try-it with a non-OMP harness.

---

---

## Crew channel (FM-574)

| | |
| --- | --- |
| **Iris room** | id `f57aa087` |
| **Spynel room** | **Superseded** — captain deletes from sidebar |
| **npm publish** | `@edheltzel/iris` (FM-573 / FM-574) |

Talk Iris factory work in the Iris room going forward.

---

---

## Dynamic model routing — Jev / Typesafe (FM-575)

**Product intent (locked):** Iris will include **dynamic model routing via Jev** ([Typesafe AI](https://typesafe.ai/) — related watch note: [Rosetta2 #177](https://github.com/edheltzel/Rosetta2/issues/177)).

| | |
| --- | --- |
| **What** | Route model/harness choices using Jev-style typed decisions with calibrated confidence (not “chat LLM as router”) |
| **What not** | Do **not** invent API shapes, endpoints, SDKs, or integration code in this VISION |
| **Phase timing** | **Locked (FM-576):** Phase 2 step **(4)** — after Herdr + Grok Bot notify |

Steal from Rosetta2 #177 + Typesafe docs at Phase 2 (4). **No invented API** in VISION or early ships.

### Additional Jev seats (Brainstorm 2026-09-19)

Beyond routing, captain locked these build-time seats (not a live chat loop). **Recall ingest scoring is first overall** ([Recall#291](https://github.com/edheltzel/Recall/issues/291)); Iris seats follow.

| Priority | Seat | Decision |
| --- | --- | --- |
| 1 | Model routing | Already Phase 2 (4) |
| 2 | Wait / land triage | Classify ready, stalled, or needs human |
| 3 | Notify gate | Atlas vs optional side channel vs stay quiet |

GitHub track: [Iris#4](https://github.com/edheltzel/Iris/issues/4) (moved from [Sentinel#50](https://github.com/edheltzel/Sentinel/issues/50)).

---

## 2) Rename / rebrand plan (Spynel → Iris)

| Layer | Plan |
| --- | --- |
| **Product name** | **Iris** everywhere captain-facing and in factory docs |
| **BigMac checkout** | **Locked (FM-568):** `/Users/ed/Developer/Iris` (was Spynel). Remotes unchanged: origin Iris, upstream agent0ai/spynel, old-spynel |
| **Binary / CLI / npm package** | **Locked (FM-560 + FM-562 + FM-569):** Phase 1 renames **`spynel` → `iris`** including **npm package name(s)**. **Hard cut** — no `spynel` binary/alias. Same rename ship as the Go module and OS-state path split |
| **Go module / imports** | **Locked (FM-563):** Phase 1 renames Go module path + imports to **iris** in the **same** rename ship as CLI/npm/config (not a follow-up epic) |
| **Workspace / OS state** | **Locked (FM-561):** workspace config and data stay in **`.spynel/`**. Installation identity, process records, and speech models live under **`$HOME/.agents/Iris`**, with a one-shot migration from legacy `UserConfigDir` / `UserCacheDir` `iris` or `spynel` directories. Do not migrate workspace config to a user-global Iris directory |
| **Hooks / env** | Track Spynel hook names (`SPYNEL_HOOK`, `.spynel-extension.yaml`) until a thin rename pass; do not block hooks on rebrand |
| **Docs / VISION** | This file + later `VISION.md` in the Iris repo once minted |
| **GitHub** | **Locked:** `edheltzel/Iris` (FM-559). Atlas creates/retargets; leave mint to Atlas |
| **Channel / room** | Spynel room may rename to Iris when Atlas flips it; factory FM ids stay |

**Rebrand ship (thin, after fork remote exists):** rename module paths, binary, install scripts, README lead — separate from capability ships. Prefer one focused PR for names, not mixed with hooks.

---

---

## UX / branding lock — blinking eyeball (FM-558)

**Captain lock (Spynel channel, 2026-09-17):** busy / loading / working animations use a **blinking eyeball** (Iris mark).

| Surface | Apply |
| --- | --- |
| **TUI spinner / status** | Replace generic spinner with blinking-eyeball animation while Iris is busy, loading, or working |
| **Optional Grok Bot status affordance** | Same mark if we show “Iris working” in a bot/status chip — still **Atlas sole door**; this is status chrome, not a second liaison |
| **Channel / room avatar** | **Non-goal:** do **not** invent a second mascot face for the Spynel/Iris channel. Marks-only rule still stands (silhouette + grid + oxide tick) — eyeball blink is **in-product motion**, not a new crew face |

**Steal later:** wire the mark into `cmd`/TUI status paths after fork mint; keep animation assets thin and reusable. Do not block Phase 0/1 on art.

---

## 3) Top exceed gaps (from FM-554 — fold first)

These three close the Firstmate + Sentinel exceed bar. **Phase 2 order locked (FM-564–566 + FM-576):** **(0) GitButler `but` → (1) no-mistakes → (2) zero-token/quota → (3) Herdr + Grok Bot notify → (4) Jev/Typesafe dynamic model routing** (Atlas sole door for notify).

| P | Gap | Firstmate / Sentinel today | Spynel day-one | Iris land |
| --- | --- | --- | --- | --- |
| **1** | **no-mistakes ship gate** | FM: shipped modes + `no-mistakes axi` + refuse. Sentinel: honor flag only; idea **#45** | Review lifecycle ≠ full gate. Hooks can shell gate (FM-551) | **Hook pack early** — `harness.after` / `task.completed` runs gate CLI; refuse dispatch when missing; modes match FM |
| **2** | **Zero-token absorb + quota-aware dispatch** | FM: `fm-watch.sh` + `quota-axi` dispatch arrays. Sentinel: absorb idle ticks; `quota-axi` companion floor | Semantic heartbeat only; **no** quota-axi | Companion watcher + hooks (or Go later): absorb benign wakes; quota arrays steer dispatch |
| **3** | **Herdr socket-first + Grok Bot transition notify** | FM: Herdr backend + fleet notify. Sentinel: `--herdr` peek; thin ui.notify; ideas **#44/#16** | Transition notify + journal exist; TG/WA channels; **no** Herdr mux | **Socket-first** Herdr external watcher + hooks. Spynel notify packaging → **Atlas/Grok Bot** tips (origin + send/skip journal). **Atlas sole door** — not TG/WA product |

**Near-term fork order (FM-543 + FM-554 + FM-564–566 + FM-576):** stand Iris fork → Phase 1 rename → Phase 2: **(0) GitButler `but` → (1) no-mistakes → (2) zero-token/quota → (3) Herdr + Grok Bot notify → (4) Jev/Typesafe dynamic model routing** → budget Go only where hooks can’t reach.

---

## 4) Capability matrix — Iris vs Firstmate vs Sentinel

Legend: **M** = must-match · **E** = must-exceed · **D** = drop from Spynel product · **L** = later · **S** = steal into Iris · **I** = inherit Spynel day-one

Full evidence table: [`docs/exceed-inventory.md`](exceed-inventory.md).

| Capability | Firstmate | Sentinel (live) | Spynel / Iris inherit | Iris target |
| --- | --- | --- | --- | --- |
| Human front door | Atlas / DA | Atlas; Sentinel plumbing | Own TUI + TG/WA | **M** Atlas sole Grok Bot captain door (**FM-567**) · TG/WA **optional side channels** only · **D** Spynel as second captain |
| Control-plane shape | Distro + skills + panes | omp plugin | Go daemon + Markdown tasks | **E** Go hub harness-agnostic · omp optional client |
| Task store | tasks-axi | tasks-axi companion | Markdown task OS **I** | **I/M** Markdown SoT day-one · **S** axi wait/lease/progress habits via hooks |
| Wait contracts | Declared pause + resurface | Thin “Needs you”; idea **#43** | Named wait, `wake_at`, Progress, leases **I** | **E/I** Keep Spynel wait board; tip transitions → notify |
| Leases / stale | — | Filed #43 | 5s / 30s **I** | **M/I** |
| Notify on transition | Fleet-local wakes | Thin ui.notify; **#44/#16** | Notify-agent + origin + journal **I** | **E** Retarget tips → **Atlas/Grok Bot** · TG/WA **optional side** only (**FM-567**) |
| Grok Bot agent talk | Single liaison pattern | #16 idea | N/A (own chat) | **M** Atlas wire · **L** Go `channel.Channel` only if localapi fails |
| **no-mistakes / ship gate** | **Shipped** | Idea **#45** | Review ≠ gate | **E** Hook pack + refuse when CLI missing (**top gap #1**) |
| **Zero-token / quota** | **Shipped** watch + quota-axi | Absorb + quota companion | Absent / different | **E** Watcher + quota-aware dispatch (**top gap #2**) |
| Git / worktrees | Treehouse + backends | — | Absent Treehouse | **E** **GitButler** via `but` hooks · **D** Treehouse default |
| Heartbeat | Watcher + bearings | Bash stamp | Semantic + primary election **I** | **M/I** + Herdr socket poke |
| Herdr / mux panes | Herdr backend shipped | `--herdr` peek | Multi-TUI `/jobs` only | **E** Socket-first watcher (**top gap #3**) · **L** Go panes |
| Harness adapters | Multi + dispatch profiles | omp-centric | Codex/Claude/Pi/ACP **I** | **E/I** Keep multi-harness · **FM-572:** first-run docs/examples **default OMP**; others supported · **S** FM dispatch profiles |
| Dynamic model routing (Jev) | — | — | Absent | **Intent (FM-575) + Phase 2 (4) (FM-576):** Jev/Typesafe routing; no API invented |
| Extensibility | Harness hooks/skills | omp plugin | Portable hooks **I** | **E/I** Hooks early for CLI gates |
| Agent Hub / list | Fleet visibility | Hub peek + #17 | `/jobs` + TUI **I** | **E** Iris owns hub at cutover |
| Review-from-artifact | RedTeam social | RedTeam | Review stage **I** | **L** after wait/notify |
| Away / AFK | Full AFK | First slice | Absent | **L** |
| Goal vs task | Multi-FM epics | Social | Goals **I** | **L** |
| Primary election | Session locks | Away owner lock | 5s/30s election **I** | **I** |

### Must-match (bar to call Iris “ready to cut over”)

1. Atlas remains the only Grok Bot front door (no second liaison).
2. Durable wait + land state a human can inspect (Spynel wait/`wake_at`/leases + tips).
3. Harness dispatch for captain’s primary agents (Claude Code + Pi + ≥1 more).
4. Hooks install path for **no-mistakes / axi / `but`** without core rewrite.
5. GitButler on the ship path (hook-driven), not Treehouse-default.
6. Herdr status via socket watcher (even if deep panes are later).
7. Zero-token absorb path + quota-aware dispatch at least at Firstmate bar.

### Must-exceed (why Iris beats Firstmate + Sentinel together)

1. **Harness-agnostic hub** — not omp-plugin-shaped forever.
2. **First-class wait contracts + leases** in the control plane (inherited Spynel strength).
3. **Notify-on-transition with origin/journal** tipped through Atlas (steal Spynel packaging + Sentinel #44 intent).
4. **Portable hook ecosystem** for Firstmate CLIs (no-mistakes, axi, `but`, quota watch).
5. **Clear cutover story** — Sentinel live until Iris wins on feel.

### Drop (do not rebuild)

- Spynel “one chat runs the crew” as a rival to Atlas.
- TG/WA as required product notify or second captain door (**FM-567:** optional side channels OK).
- Treehouse as default.
- Day-one Go Herdr Hub / multi-window primary-election complexity beyond what Spynel already has.
- Absorb-only “forever omp Sentinel” as the destination.

### Spynel day-one inheritance (summary — detail in FM-554)

| Free | Needs hooks / companion | Needs Go |
| --- | --- | --- |
| Markdown tasks/goals, leases, `wake_at`, Progress | no-mistakes, axi, `but`, Herdr socket poke, zero-token/quota scripts | New harness adapter; optional Grok Bot `channel.Channel` |
| Review lifecycle, multi-harness catalog | Notify tip retarget → Atlas | Deep Herdr panes if socket fails |
| Primary election, semantic heartbeat, notify-agent jobs | — | — |
| Portable hooks (5 names) | Hook pack = Phase 2 | — |

---

## 5) Phased conversion plan (smallest ships; hooks first)

Dual track always: **Sentinel ships continue** on their own FM ids. Iris phases below are the rebuild track.

### Phase 0 — Spec locked (this document)

- VISION accepted by captain (or listed HIL Qs answered).
- FM-554 exceed inventory folded (FM-555).
- No fork remote invented.

**Done:** see §7.

### Phase 1 — Fork mint + rename (Atlas-owned create; rename locked)

- **Fork home locked:** `edheltzel/Iris`.
- **Atlas** creates/retargets remote (Intern does not).
- Thin fork from `agent0ai/spynel` already live as `edheltzel/Iris`. BigMac path: **`/Users/ed/Developer/Iris`** (FM-568).
- **Rename locked (FM-560 / FM-561 / FM-562 / FM-563 / FM-569 / grill 2026-09-18):** Phase 1 ships **one rename ship** (or tightly stacked thin PRs) covering (one rename ship):
  1. **CLI / binary / npm package** `spynel` → **`iris`** (FM-560 + **FM-569** npm).
  2. **Workspace / OS state (FM-561):** keep workspace config and data in **`.spynel/`**. Put installation identity, process records, and speech models under **`$HOME/.agents/Iris`**, with a one-shot migration from legacy `UserConfigDir` / `UserCacheDir` `iris` or `spynel` directories.
  3. **Hard cut (FM-562):** after Phase 1 lands, there is **no** `spynel` binary and **no** `spynel` alias/compat shim — callers use `iris` only.
  4. **Go module / import paths (FM-563):** rename module + imports to **iris** in that **same** Phase 1 rename ship — not a later epic.
  Prefer thin rename-focused PR(s) on the Iris tip — do **not** defer past Phase 1. Hook env names (e.g. `SPYNEL_HOOK`) may follow in the same or next thin PR; leftover `SPYNEL_*` env strings are not a license to keep a `spynel` binary/alias or old Go module path.

### Phase 2 — Hook pack + exceed gaps (full order — FM-564–566 / FM-576)

Smallest ships, one concern each. **Order locked (grill 2026-09-18):**

1. **Extension install proof** — hello hook; doctor green (prerequisite).
2. **(0) GitButler `but` hooks** (FM-564) — shell `but` on task/harness events (not GB-native store).
3. **(1) no-mistakes ship gate** (FM-565) — `harness.after` / `task.completed` → gate CLI; refuse when missing; journal skip/fail.
4. **(2) Zero-token / quota** (FM-566) — absorb benign wakes; quota-aware dispatch (script/hook companion).
5. **(3) Herdr socket-first + Grok Bot transition notify** (FM-566) — external Herdr watcher + hooks; on wait/land tip **Atlas/Grok Bot** (origin + send/skip journal). **Atlas sole door** — no TG/WA product.
6. **(4) Jev / Typesafe dynamic model routing** (FM-575 intent + **FM-576** phase) — after Herdr/notify; **do not invent API** until real Typesafe/Jev docs land (Rosetta2 #177).
7. **tasks-axi bridge hook** — claim/complete side effects (Markdown OS remains SoT); may interleave after `(0)` if needed, but does **not** jump ahead of `(0)`–`(1)`.

### Phase 3 — Wait / lease / progress polish

1. Confirm Spynel wait/`wake_at`/leases under Iris branding.
2. Progress handoff log; align field names with Sentinel #43/#44 where useful.
3. Dispatch profiles (steal FM `crew-dispatch.json` shape).

### Phase 4 — Grok Bot / Atlas wire deepen

1. Prefer **localapi/CLI client** from Atlas/crew.
2. Only if insufficient: thin Go `channel.Channel` — separate ship; RedTeam + captain HIL.

### Phase 5 — Cutover readiness

1. **Cutover criteria: parked (FM-571)** — do not invent the green checklist yet; Phase 2 work still ships in locked order.
2. Captain happy-with-Iris bar (subjective + demo) when cutover is unparked.
3. **Until then:** keep hacking Sentinel (**no freeze**, **no sunset date**).

**Ship discipline (standing):** small frequent commits; individual focused PRs; no large multi-concern PRs.

---

## 6) Open HIL questions

| # | Question | Notes |
| --- | --- | --- |
| 1 | **Iris GitHub owner / repo name?** | **Locked:** `edheltzel/Iris` (grill 2026-09-18 / FM-559). Atlas creates/retargets. |
| 2 | Binary name: `iris` vs keep `spynel` until cutover? | **Locked (FM-560):** rename to `iris` in Phase 1 (CLI/binary/package) |
| 2b | Workspace / OS state split | **Locked (FM-561):** workspace stays `.spynel/`; identity, processes, and speech use `$HOME/.agents/Iris` with one-shot legacy user-dir migration |
| 2c | Keep `spynel` alias after rename? | **Locked (FM-562):** hard cut — no binary/alias |
| 2d | Defer Go module rename? | **Locked (FM-563):** same Phase 1 ship as CLI/config |
| 2e | Defer npm package rename? | **Locked (FM-569):** same Phase 1 ship as CLI/Go/config |
| 2f | npm publish scope/name? | **Locked (FM-573):** `@edheltzel/iris` |
| 3 | When does Sentinel stop net-new features? | **Locked (FM-570):** keep hacking — no freeze/sunset. **Cutover criteria parked (FM-571)** |
| 4 | GitButler: hook-shell `but` only for v1, or schedule GB-native store spike? | FM-543 open Q |
| 5 | Heartbeat: spike semantic+Herdr on Phase 2/3, or keep Sentinel bash stamp until measured? | |
| 6 | Grok Bot: localapi/client-only for Phase 4, or budget Go channel? | FM-551: hooks ≠ chat ABI |
| 7 | Task SoT long-term: Markdown forever, or eventual axi-unified store? | Bridge via hooks first |
| 8 | Room rename Spynel → Iris? | **Locked (FM-574):** Iris room `f57aa087`; Spynel superseded |
| 9 | Zero-token: port FM `fm-watch` shape vs thin Iris-native script? | FM-554 gap #2 |
| 10 | After no-mistakes: zero-token vs Herdr+notify? | **Locked (FM-566):** (2) zero-token/quota → (3) Herdr+Grok Bot notify |
| 11 | TG/WA as product notify? | **Locked (FM-567):** optional side channels; Atlas sole Grok Bot captain door |
| 12 | Day-one install shape? | **Locked (FM-571):** npm **primary**, Go binary **secondary** |
| 13 | Cutover criteria checklist? | **Parked (FM-571)** |
| 14 | First-run default harness? | **Locked (FM-572):** OMP in docs/examples; others supported |
| 15 | Jev/Typesafe dynamic model routing — which phase? | **Locked (FM-576):** Phase 2 step **(4)** — see Rosetta2 #177; no API invented |

---

## 7) Done definition — Phase 0 (spec locked)

Phase 0 is **done** when all of the following are true:

1. This file exists at `docs/iris-vision.md` as the living SoT. `box` grok-ship/reports is not parallel law.
2. Product one-liner, non-goals, rename plan, capability matrix, phased plan, and open HIL Qs are present.
3. **FM-554 exceed inventory is folded** — companion provenance is `docs/exceed-inventory.md`.
3b. **FM-558 UX lock present** — blinking eyeball for busy/loading/working; no second channel mascot face. — especially top gaps: (1) no-mistakes hook gate (2) zero-token + quota dispatch (3) Herdr socket-first + Grok Bot transition notify with Atlas sole door.
4. Dual-track lock is explicit: Sentinel live + Iris destination.
5. Fork home is locked as `edheltzel/Iris` (FM-559); **create/retarget left to Atlas** — Intern does not invent remotes.
6. Atlas has tips against **FM-553** / **FM-555**; captain can accept Phase 0 from this file.

**Phase 0 exit → Phase 1:** fork home locked (`edheltzel/Iris`); Atlas creates/retargets; Phase 1 rename ship: CLI/binary/**npm** → `iris` / publish **`@edheltzel/iris`** (FM-560+569+**573**), workspace `.spynel/` retained with OS state migrated to `$HOME/.agents/Iris` (FM-561), hard cut no `spynel` binary/alias (FM-562), Go module/imports → iris (FM-563).

---

## References

| Doc / issue | Path / URL |
| --- | --- |
| Dual-track gap map | Historical factory report FM-543 (not in this repo). Tracking: Sentinel #49 |
| Extensibility scout | Historical factory report FM-551 (not in this repo) |
| Exceed inventory | [`docs/exceed-inventory.md`](exceed-inventory.md) |
| Dual-track tracking | https://github.com/edheltzel/Sentinel/issues/49 |
| Wait / notify / no-mistakes / Grok Bot | Sentinel #43 · #44 · #45 · #16 |
| Upstream Spynel | https://github.com/agent0ai/spynel |
| BigMac checkout | `/Users/ed/Developer/Iris` (FM-568) |
