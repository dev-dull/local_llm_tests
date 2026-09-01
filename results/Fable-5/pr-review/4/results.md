# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** `tickerCmd`, the added `time.Sleep(33 * time.Millisecond)` — this works (commands run in their own goroutines, so the sleep blocks neither input handling nor rendering, and the tick→Update→tick loop is now throttled), but it is not the idiomatic mechanism. `tea.Tick(time.Second/30, ...)` exists for exactly this and is what PR 2 does; prefer that.
2. **[nit]** Same hunk — a fixed 33 ms sleep does not account for the time spent processing each frame, so the real frame rate is somewhat below 30 fps and drifts with load. `tea.Tick` aligns to interval boundaries and avoids this.

**Summary:** The change genuinely fixes the busy-spin and is safe, but it reimplements frame scheduling that Bubble Tea already provides; if PR 2 is on the table, merge that instead.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. none

**Summary:** This is the idiomatic fix for the busy-spinning ticker: `tea.Tick(time.Second/30, ...)` schedules the next `tickerMsg` ~33 ms out instead of returning it immediately, capping the loop at roughly 30 fps exactly as described. Correct as-is, and preferable to PR 1's `time.Sleep` approach.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** Tick handler in `Update` (not touched by the diff) — the animation still advances with the hard-coded `m.spun += 0.05`. The new `m.spinSpeed` field is written by the `+`/`-` handlers, `reset`, and `main`, but never read anywhere. Pressing +/- therefore has no visible effect at all; the line must become `m.spun += m.spinSpeed`.
2. **[nit]** Hunk at the key handlers — `+` has no upper clamp, so holding it grows the per-tick step without bound and the spin degenerates into flicker. Consider a `math.Min` cap symmetric with the `0.01` floor on `-`.
3. **[nit]** `reset()` and `main()` hunks — the default `0.05` is now duplicated in two places; a named constant would keep them from drifting.

**Summary:** The plumbing (field, key bindings, clamp, reset, help text) is all present, but the one line that would consume the speed was never changed, so the feature is dead code as shipped.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** APPROVE

**Findings:**
1. **[nit]** `randInt` hunk — `rand.Intn` panics when `max <= 0`. Current callers pass 20 and 3 so this is fine (and the old millisecond-modulo version would have panicked on 0 too), but the doc comment could note the precondition.

**Summary:** Replacing the wall-clock-millisecond hack with `rand.Intn` fixes the described lockstep initialization (consecutive calls within the same millisecond previously returned identical values), and the now-unused `time` import is correctly dropped; since Go 1.20 the global PRNG is auto-seeded, so no `rand.Seed` is needed.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** The added clamp at the end of `View` — truncating to `m.height` keeps the header, but it guarantees the footer separator is never visible: wave/rain/pulse each emit `m.height` content lines, so the total is `m.height + 3` and the last three lines (including the footer) are always cut. The cleaner fix is to give the renderers the available height (`m.height - 3`) so header, content, and footer all fit.
2. **[nit]** Same hunk — the `m.height > 0` guard for the pre-`WindowSizeMsg` state is good; note the clamp does nothing about over-*wide* lines, which is fine for this PR's stated scope but worth remembering.

**Summary:** Does what it claims — the header can no longer be pushed out of view — and the guard against an unset height is correct, but it trades the footer away permanently; sizing the content region would fix both.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** New loop in the tick handler: `for i, c := range m.rainPos { c.y += c.v; ... }` — `c` is a copy of the slice element, so every mutation (`c.y`, `c.col`, the wrap-around reset) is discarded at the end of each iteration. Combined with the (correct) removal of the in-`View` updates, no code updates rain positions any more: the characters stay at their initial negative `y` values and the rain mode renders a permanently empty screen. The loop must index the slice as the old code did (`m.rainPos[i].y += m.rainPos[i].v`, etc.) or take a pointer (`c := &m.rainPos[i]`).
2. **[nit]** Same hunk — gating on `m.mode == modeRain` preserves the base behavior (rain only advances while visible), which is a reasonable choice; just noting it is a behavior the base also had, not a regression.

**Summary:** The refactor's goal is right — mutating model state in `View()` is a real Bubble Tea anti-pattern — but the range-variable copy bug means the moved code is a no-op and rain mode is completely broken; fix the loop to mutate the slice elements.
