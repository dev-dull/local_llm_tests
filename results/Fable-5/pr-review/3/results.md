# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** `tickerCmd`, the added `time.Sleep(33 * time.Millisecond)` — this does fix the busy-spin (each tick command now parks its goroutine for 33ms before delivering the next `tickerMsg`, and commands run off the event loop so input stays responsive), but it is the non-idiomatic way to do it in Bubble Tea. `tea.Tick(time.Second/30, ...)` expresses the same thing declaratively, is what the framework provides for exactly this, and is trivially swappable here (see PR 2).
2. **[nit]** Same hunk — the effective frame rate is slightly below 30 fps because the 33ms sleep is added on top of per-frame processing time rather than scheduling frames on a fixed cadence. Harmless for this app, but `tea.Tick` handles this too.

**Summary:** The change genuinely eliminates the CPU spin and does what it claims with no functional defect, but PR 2's `tea.Tick` version of the same fix is strictly better; if both are open, prefer that one.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. **[nit]** `tickerCmd` — `tea.Tick` measures its delay from when the command runs, so frame timing drifts slightly with processing time rather than firing on exact 1/30s boundaries. Irrelevant for a toy animation; noting only for completeness.

**Summary:** This is the idiomatic fix for the busy loop: `tea.Tick(time.Second/30, ...)` schedules one delayed `tickerMsg`, and the existing `Update` handler already re-issues `tickerCmd()` on every tick, so the animation keeps running at ~30 fps. Correct as-is, and preferable to PR 1.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** `Update`, tick handling — the PR adds, clamps, resets, and documents `spinSpeed`, but never uses it. The tick handler still reads `m.spun += 0.05` (untouched by the diff), so pressing `+`/`-` changes a field that has no effect on the animation. The line must become `m.spun += m.spinSpeed`. As shipped, the feature simply doesn't work, while the help bar claims it does.
2. **[suggestion]** `case "+", "=":` hunk — `+` has no upper clamp, so holding it grows `spinSpeed` without bound; once the (missing) wiring is added, the spin would eventually alias into visual chaos. Mirror the lower clamp with something like `math.Min(0.5, ...)`.
3. **[nit]** `reset()` and `main()` both set the default `0.05` — fine, but consider a named constant so the default, the reset value, and (after fixing finding 1) the tick increment can't drift apart.

**Summary:** All the scaffolding is present (field, key handling, clamp, reset, help text, initial value) but the one line that would make the keys actually do anything was never changed, so the feature is inert. Must be fixed before merge.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** APPROVE

**Findings:**
1. **[nit]** `randInt` — on Go < 1.20 the global `math/rand` source is unseeded, so the rain layout is identical on every run (though still uncorrelated within a run, which is what the PR claims to fix). If the module targets an older Go, add a seed; on Go 1.20+ the global source is auto-seeded and this is moot.
2. **[nit]** `randInt` — `rand.Intn` panics for `max <= 0`, but so did the old modulo (division by zero), and both call sites pass positive constants, so this is parity, not a regression.

**Summary:** Correct and complete: the wall-clock-millisecond hack really did return near-identical values for consecutive calls (the whole init loop runs within one millisecond), and `rand.Intn` fixes that. Dropping the now-unused `time` import is right — `randInt` was its only user.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** The added clamp at the end of `View()` — truncating from the bottom keeps the header visible (the stated goal) but silently drops the footer in essentially every frame: `renderWave`, `renderRainStatic`, and `renderPulse` each return exactly `m.height` lines, so with the 2 header lines and 1 footer line the view is always `m.height + 3` lines and the clamp always cuts the last 3 (footer plus the bottom two content rows). The real fix is to give the render functions a content height of `m.height - 3` so header, content, and footer all fit; the clamp is then a safety net rather than the mechanism.
2. **[nit]** Same hunk — the `m.height > 0` guard is correct for the pre-`WindowSizeMsg` frame where height is still 0; good defensive touch.

**Summary:** The change does what the description promises — the header and keybinding help can no longer scroll out of sight — and introduces no defect, but it papers over the underlying sizing problem and permanently hides the footer and bottom content rows. Mergeable as a stopgap with the sizing follow-up noted.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** `Update`, tick handler — `for i, c := range m.rainPos` iterates over *copies* of each `rainChar`; the mutations `c.y += c.v`, `c.col = ...`, and `c.y = -5` are applied to the loop-local copy and discarded each iteration. Since the in-place update loop was correctly removed from `View()`, nothing updates `m.rainPos` at all: every drop stays at its initial negative `y` forever and rain mode renders as a blank screen. The loop must index the slice directly (`m.rainPos[i].y += m.rainPos[i].v`, etc.), as the deleted `View()` code did.
2. **[nit]** Same hunk — gating the update on `m.mode == modeRain` faithfully preserves the old behavior (positions only advanced while rain mode was on screen), so that part of the port is fine.

**Summary:** The refactor's goal is right — `View()` should not mutate model state — and the removal from `View()` is correct, but the relocated loop mutates range-variable copies instead of the slice elements, which completely breaks the rain animation. One-line-per-statement fix required before merge.
