# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** `tickerCmd`, the added `time.Sleep(33 * time.Millisecond)` — sleeping inside the command does cap the loop at ~30 fps (commands run on their own goroutine, so input handling is not blocked), but it is the wrong idiom for Bubble Tea. `tea.Tick(time.Second/30, ...)` exists for exactly this purpose (see PR 2) and is what the framework documents for frame scheduling.
2. **[nit]** same hunk — the sleep does not account for time spent updating/rendering the frame, so the effective rate is slightly below 30 fps and drifts with load. `tea.Tick` truncates to the interval boundary and drifts less.

**Summary:** The change genuinely fixes the busy-spin and is safe to merge, but PR 2 achieves the same goal the idiomatic way; if both are open, prefer PR 2 and close this one.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. none

**Summary:** This is the canonical fix: `tea.Tick(time.Second/30, ...)` schedules each frame with a real timer instead of returning a message immediately, eliminating the busy loop. The re-arm pattern (returning `tickerCmd()` from the tick branch of `Update`) already exists in the base code, so a single `tea.Tick` per frame is exactly right.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** tick handler in `Update` (not touched by the diff) — the diff adds and adjusts `m.spinSpeed`, but the line that advances the spin is still the hardcoded `m.spun += 0.05`. Nothing ever reads `spinSpeed`, so pressing +/- has no visible effect whatsoever. The tick handler must become `m.spun += m.spinSpeed` for the feature to work.
2. **[suggestion]** `case "+", "=":` hunk — `spinSpeed` has a minimum clamp (0.01) on `-` but no maximum on `+`; holding `+` lets the speed grow without bound. Clamping to a sane ceiling (e.g. `math.Min(1.0, ...)`) would match the "clamped to a sensible minimum" spirit of the description.
3. **[nit]** same hunk — the +/- keys are handled in every mode, not just Spin. Harmless, but the help text says "spin speed", so either gate them on `m.mode == modeSpin` or leave as-is knowingly.

**Summary:** The plumbing (field, default in `main` and `reset`, help text) is all present, but the one line that would consume the new field was never changed, so the PR does not do what it claims — the animation speed is unchangeable. Must be fixed before merge.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** `randInt` hunk — on Go versions before 1.20 the global `math/rand` source is deterministically seeded, so every run of the program would get the identical rain layout (though no longer correlated within a run). If the module targets Go ≥ 1.20 this is a non-issue (the global source is auto-seeded); otherwise add a seed, or better, use `math/rand/v2` (`rand.IntN`) which is always well-seeded.
2. **[nit]** same hunk — `rand.Intn` panics for `max <= 0`. All current call sites pass positive constants (20, 3), so this is fine today, but the old implementation tolerated any nonzero max; worth a doc comment on `randInt` if you want to keep the contract explicit.

**Summary:** The diagnosis is correct — the old millisecond-clock `randInt` returned the same value for every call within the same millisecond, which is exactly why all rain characters started in lockstep — and switching to a real PRNG fixes it, including the reuse in `reset()`. Only version-dependent seeding is worth a look before merge.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** the added clamp after `lines = append(lines, foot)` — truncating from the bottom keeps the header visible, but it means the footer separator (appended last) is the first thing cut, and it will be cut on *every* overflowing frame, since e.g. `renderWave` alone already emits `m.height` lines and the view adds header + blank + footer on top of that (`m.height + 3` total). So with this PR the footer effectively never renders. A cleaner fix is to give the renderers the available height (`m.height - 3`) so the frame fits, or to keep the footer by trimming to `m.height-1` lines and re-appending it.
2. **[nit]** same hunk — the root cause is that `renderSpin` sizes its grid from the radius (up to ~27 rows) with no regard for `m.height`; this clamp treats the symptom. Fine as a safety net, but worth a follow-up.

**Summary:** The clamp is correct, guards the pre-`WindowSizeMsg` case with `m.height > 0`, and achieves the stated goal of keeping the header on screen; it just silently sacrifices the footer whenever it kicks in, which the renderers-aware fix would avoid.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** added loop in the tick handler, `for i, c := range m.rainPos` — `c` is a *copy* of the slice element, so `c.y += c.v`, `c.col = ...`, and `c.y = -5` mutate the copy and are discarded at the end of each iteration. Combined with removing the in-`View` update, rain characters are frozen at their initial off-screen negative `y` positions forever: rain mode now renders nothing at all. The loop must index the slice as the old code did (`m.rainPos[i].y += m.rainPos[i].v`, etc.) or take a pointer (`c := &m.rainPos[i]`).

**Summary:** The refactor's goal is right — `View` should be pure and the base code's in-`View` mutation is a real wart — but the moved loop iterates over element copies, so the state update is lost and the rain animation is completely broken. Fix the range loop to mutate the slice elements before merging.
