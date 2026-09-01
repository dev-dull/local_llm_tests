# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** line 425 — `time.Sleep(33 * time.Millisecond)` works to cap frame rate at ~30 fps, but it is a blocking sleep inside the `tea.Cmd` closure. Bubble Tea provides `tea.Tick` for non-blocking, timer-based frame scheduling (see PR 2), which is the more idiomatic approach.
2. **[nit]** The sleep is hardcoded to 33 ms (~30 fps). This is reasonable, but `time.Second/30` would make the intent self-evident.

**Summary:** The change correctly eliminates busy-spinning by introducing a 33 ms delay per frame. It is a functional fix, but `tea.Tick` (PR 2) is the preferred Bubble Tea idiom.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. **[none]**

**Summary:** This is the correct, idiomatic fix. `tea.Tick(time.Second/30, func(time.Time) tea.Msg { return tickerMsg{} })` schedules the next frame via Go's `time.AfterFunc`, avoiding both the blocking sleep of PR 1 and the original busy-spin. It compiles cleanly and does exactly what the description claims.

## PR 3: Add +/- keys to control spin speed

**Verdict:** APPROVE

**Findings:**
1. **[suggestion]** line 483-488 — The `+` key handler lacks a clamp on the upper bound; `spinSpeed` can grow unbounded. A simple cap like `m.spinSpeed = math.Min(m.spinSpeed+0.02, 1.0)` would prevent runaway values, though this is not a defect as-is.
2. **[suggestion]** line 487 — `math.Max(0.01, m.spinSpeed-0.02)` correctly prevents the speed from going below 0.01. The default `reset()` and `main()` initializers both set `spinSpeed = 0.05`, so the field has a well-defined starting value. The help bar is updated consistently.
3. **[nit]** `"="` is handled alongside `"+"` on most keyboards, which is fine for usability, but the comment could clarify that `+` and `=` are treated equivalently.

**Summary:** The PR correctly adds a new `spinSpeed` field, handles `+`/`-` keys with a sensible floor, resets on `r`, and updates the help bar. The only improvement would be an upper bound on `+`.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** line 353 — The old `randInt` used `time.Now().UnixNano() / 1e6` which always returned the same value for calls within the same millisecond, making all initial rain positions correlated. Replacing it with `rand.Intn(max)` fixes the correlation. However, Go 1.20+ seeds `math/rand` automatically, so no explicit `rand.Seed` call is needed. The PR should add `"math/rand"` to the import block and remove the `"time"` import — which it does correctly.
2. **[suggestion]** `rand.Intn` panics if `max <= 0`. The existing callers pass `20` and `3`, so this is not a current risk, but a guard or comment would be prudent for future maintainability.

**Summary:** A clean fix. Removing the broken clock-based "PRNG" and using `math/rand` resolves the lockstep initialization issue.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** lines 561-563 — `lines = lines[:m.height]` truncates from the *end*, not the beginning. The footer (last element) gets dropped while the top lines — including the header — may still be pushed off-screen. To keep the header visible, the code should either discard excess lines from the front (`lines = lines[len(lines)-m.height:]`) or, better yet, clamp individual mode renders to available height (see PR 6's broader concern about View emitting too many lines).
2. **[suggestion]** The check `m.height > 0` prevents a panic on zero-height windows, which is good, but a negative `m.height` would still produce a panic or unexpected behavior.

**Summary:** The fix truncates from the wrong end of the slice. The footer is dropped while the header can still scroll off-screen, so the stated goal is not achieved.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** lines 581-587 — The range loop `for i, c := range m.rainPos` copies each `rainChar` into the loop variable `c`. Modifying `c.y`, `c.col` etc. changes the copy, **not** `m.rainPos[i]`. The slice elements are never updated, so the rain animation will be static (only the first frame's positions are rendered). This must use `m.rainPos[i].y += m.rainPos[i].v` directly (or iterate by index and write back to `m.rainPos[i]`).
2. **[suggestion]** The same logic is duplicated from the View branch (lines 174-180 in the base). Once this is fixed, the two copies should ideally be deduplicated, but that is a future refactor.

**Summary:** The PR's intent is sound — mutating state in `View()` is an anti-pattern. However, the range loop overwrites copies, so the fix does not actually advance the rain positions. It needs to write to `m.rainPos[i]` directly.
