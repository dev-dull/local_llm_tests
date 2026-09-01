# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** `tickerCmd`, the added `time.Sleep(33 * time.Millisecond)` — this does fix the busy-spin (the command runs in its own goroutine, so sleeping there is safe and does throttle the loop), but it is the wrong idiom for Bubble Tea. `tea.Tick` exists precisely for frame scheduling and is what PR 2 uses; a hand-rolled sleep also accumulates drift, since the 33 ms is added on top of whatever time Update/View take, so the real frame rate is somewhat below 30 fps and varies with render cost.
2. **[nit]** Same hunk — the sleep runs before the very first `tickerMsg` from `Init` as well, so startup is delayed by one frame. Harmless here, but `tea.Tick`'s deadline-based scheduling avoids thinking about this class of thing.

**Summary:** The change genuinely does what it claims — CPU usage drops from a pegged core to ~30 fps pacing — and has no correctness defect, but `tea.Tick` (PR 2) is the clearly better way to achieve the same goal and should be preferred if both PRs are on the table.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. **[nit]** `tickerCmd` — `tea.Tick` schedules a single message, so continuous animation relies on `Update` re-issuing `tickerCmd()` on every `tickerMsg`. The base code already does exactly that (line `return m, tickerCmd()` in the tick branch), so this works; just noting the frame chain would silently stop if that re-issue were ever removed. If fixed-cadence ticks independent of processing time mattered, `tea.Every` would be an alternative, but for this app `tea.Tick` is exactly right.

**Summary:** This is the idiomatic fix for the busy-spinning event loop: frames are scheduled at `time.Second/30` by the runtime instead of being returned immediately, the `time` import remains needed, and no behavior other than pacing changes. Correct as-is, and preferable to PR 1's `time.Sleep` approach.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** Tick handler in `Update` (unchanged by the diff, line `m.spun += 0.05`) — `spinSpeed` is added, initialized, clamped, reset, and documented, but it is never *read*: the spin angle still advances by the hard-coded `0.05` every tick. Pressing `+`/`-` therefore has no visible effect whatsoever. The line must become `m.spun += m.spinSpeed` for the feature to exist at all.
2. **[suggestion]** `+`/`-` key handling hunk — `+` has no upper clamp, so holding it grows the speed without bound. Symmetric clamping (e.g. `math.Min(0.5, ...)`) would match the description's "clamped to a sensible minimum" spirit on both ends.
3. **[nit]** `reset()` and `main()` — the default `0.05` is now written in two places. Consider a `const defaultSpinSpeed = 0.05`, or have `main` call `m.reset()` (note `reset` also zeroes `mode`/`tick`, which is fine at startup) so the default lives once.
4. **[nit]** `reset()` hunk — `m.mode = 0` was already a pre-existing quirk (reset also switches you out of your current mode), untouched here; fine to leave, just noting the new key help sits next to it in the same reset story.

**Summary:** The scaffolding (field, keys, clamp, reset, help text) is all present, but the one line that would make the feature work — using `spinSpeed` to advance `spun` — was never changed, so the PR does not do what it claims. Must be fixed before merging.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** APPROVE

**Findings:**
1. **[nit]** `randInt` hunk — `rand.Intn` panics for `max <= 0` where the old expression divided-by-zero'd only for `max == 0`; both call sites pass constants 20 and 3, so this is academic. On Go < 1.20 the global `rand` source is unseeded (deterministic across runs); values within a run are still uncorrelated, which is what the bug report is about, and on Go ≥ 1.20 the global source is auto-seeded, so no `rand.Seed` call is needed.

**Summary:** Correct diagnosis and correct fix: the old `randInt` returned the same wall-clock-derived value for calls in the same millisecond, so all raindrops got near-identical `y` and `v`; `rand.Intn` gives independent values. The `time` import removal is right too — after this change nothing else in the file uses `time` (the ticker returns immediately in the base), so it compiles cleanly.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** Truncation hunk at the end of `View` — clamping to `m.height` keeps the header (line 0) visible, which fixes the stated problem, but the footer was appended *last*, so whenever a mode renders `m.height` body lines (wave, rain, and pulse all do: header + blank + `m.height` body rows + footer = `m.height + 3` lines) the footer separator is now always cut off, along with the bottom two body rows. The cleaner fix is to size the body to the available space — have the render functions produce `m.height - 3` rows (or truncate the body slice before appending the footer) — so both header and footer survive.
2. **[nit]** Same hunk — the `m.height > 0` guard is good (avoids `lines[:0]` blanking the pre-first-`WindowSizeMsg` frame).

**Summary:** The change is safe and does keep the header on screen as claimed, so it's mergeable, but it trades the header problem for a permanently clipped footer and hidden bottom rows; reserving vertical space for the chrome in the renderers is the better version of this fix.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** New tick-handler hunk, `for i, c := range m.rainPos` — `c` is a *copy* of each `rainChar`; the assignments to `c.y`, `c.col`, and the wrap-around all mutate the copy and are discarded at the end of each iteration. Since the old in-place update in `View` was removed, rain positions now never advance at all: every drop stays at its initial negative `y`, off-screen forever, and Rain mode renders as a blank screen. The loop must write through the slice (`m.rainPos[i].y += m.rainPos[i].v`, etc., exactly as the removed `View` code did) or iterate by index.
2. **[nit]** Same hunk — gating the update on `m.mode == modeRain` preserves the old behavior (previously the update ran only inside `case modeRain` in `View`), so that part is a faithful move; only the copy semantics are wrong.

**Summary:** The refactoring goal is right — `View` should not mutate model state — but the transplanted loop iterates over value copies, so the rain animation stops working entirely. Fix the loop to mutate the slice elements and this becomes a clean approve.
