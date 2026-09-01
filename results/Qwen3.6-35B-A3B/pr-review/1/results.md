# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** `tickerCmd()`, line returning `time.Sleep(33 * time.Millisecond)` — sleeping inside the Cmd's goroutine blocks that goroutine for 33 ms each frame, which does reduce CPU, but it's not idiomatic Bubble Tea. The correct mechanism is `tea.Tick`, which schedules the next frame through the Bubble Tea runtime. PR 2 implements this properly.

**Summary:** The change correctly addresses the busy-spin by adding a 33 ms delay (~30 fps), and it works. However, blocking a goroutine with `time.Sleep` is the wrong tool for Bubble Tea — use `tea.Tick` instead (see PR 2).

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. **[none]** — The diff cleanly replaces the busy-loop Cmd with `tea.Tick(time.Second/30, …)`, which is the idiomatic Bubble Tea way to schedule periodic animation frames. No defects.

**Summary:** This is the right fix. `tea.Tick` schedules the next tick through Bubble Tea's event loop rather than blocking a goroutine, giving the runtime control over scheduling and keeping the model update path non-blocking.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** `Update()` lines adding `"+", "="` and `"-"` case handlers (adds `spinSpeed`) — the PR adds a `spinSpeed` field, key handlers to adjust it, a reset-default, and help-bar text, but **never reads `spinSpeed` in the rendering path**. `renderSpin()` is never called — the spin mode increments `m.spun += 0.05` directly on line 138 of the base `Update()` and `renderSpin()` is called from `View()`. The speed control has zero effect.
2. **[suggestion]** `main()` init `m := &model{spinSpeed: 0.05}` — fine, but redundant with the `reset()` default since `Init()` runs before the first frame anyway. No harm done.

**Summary:** The feature is structurally complete (field, keys, reset, help text) but functionally empty because `spinSpeed` is never consumed by the tick logic or any render function. The PR must wire `spinSpeed` into the per-tick `m.spun += 0.05` line (replace the hard-coded `0.05` with `m.spinSpeed`) for the feature to work.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** `randInt()` — replaces the wall-clock-based `randInt` with `rand.Intn(max)`, but never calls `rand.Seed()` (or sets up a seed elsewhere). In Go 1.20 and earlier, `math/rand`'s global source is seeded with 1 by default, so `rand.Intn` would return the **same sequence every run**, making the rain initialization perfectly deterministic and defeating the purpose of the fix. Go 1.20 auto-seeds the global source, but the code should be explicit — call `rand.Seed(time.Now().UnixNano())` in `init()` or `main()`.
2. **[suggestion]** Removes the `time` import but leaves `math` — `math` is still used elsewhere (wave, spin, pulse), so that's correct.

**Summary:** Replacing the correlated wall-clock PRNG with `rand.Intn` is the right idea, but without an explicit seed the global source is predictable (Go < 1.20) or at least unclear. Add an explicit `rand.Seed` call (or gate on `rand.New(rand.NewSource(...))`) to guarantee the fix works across Go versions.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** `View()` lines after appending the footer — clamping `lines` to `m.height` with `lines = lines[:m.height]` is a blunt instrument. It keeps the header but discards the footer and any trailing visual content. For Wave/Pulse modes the output is exactly four lines (center + shadow/status) so it never overflows and the clamp is dead code. It only triggers for Spin (≈2·radius+1 lines) and Rain (one char per column). A tighter fix would trim the middle blank rows while preserving header + footer, or clamp the render functions themselves to produce at most `m.height − 2` content lines.

**Summary:** The clamp prevents the header from scrolling off-screen, but it also silently drops the footer. Consider a smarter trim that preserves the footer, or clamp within each render function.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** `Update()`, tick handler, `for i, c := range m.rainPos` — the PR moves the rain state mutation out of `View()` (correct goal) but iterates **by value** (`for i, c := range m.rainPos`). Each `c` is a **copy** of `rainChar`; mutating `c.y`, `c.col` modifies only the copy. The original `m.rainPos[i]` is untouched, so the rain state never actually advances. It must be `for i := range m.rainPos` followed by `m.rainPos[i].y += m.rainPos[i].v` (the pattern used in the original View code that this PR is removing).
2. **[suggestion]** `Update()`, `for i, c := range` — same shadowing issue: `i` is already in scope from outer for-ranges, so `for i, c := range` could also cause a compile error in some contexts. Even if it compiles, the value semantics are the real blocker.

**Summary:** The architectural change (state in Update, View pure) is correct, but the implementation has a copy-vs-reference bug that makes the mutation a no-op. Fix the range to `for i := range m.rainPos` and mutate `m.rainPos[i]` directly.
