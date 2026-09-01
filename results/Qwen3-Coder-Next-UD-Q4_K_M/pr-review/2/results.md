# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[nit]** line 425 — `time.Sleep(33 * time.Millisecond)` is a reasonable quick fix, but it blocks the entire goroutine. Bubble Tea provides `tea.Tick` for non-blocking frame scheduling.
2. **[suggestion]** The 33ms sleep results in ~30fps, but using `tea.Tick` would be more idiomatic and integrate better with Bubble Tea's event loop.

**Summary:** The change does cap CPU usage by adding a delay, but it's a brute-force approach that blocks the command goroutine. Using `tea.Tick` would be more idiomatic.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. **[none]** 

**Summary:** This is the correct, idiomatic solution. `tea.Tick` integrates properly with Bubble Tea's event loop, schedules frames at ~30fps, and doesn't block the goroutine. It fixes the busy-spinning issue more cleanly than PR 1.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** line 487 — `math.Max(0.01, m.spinSpeed-0.02)` is called, but `spinSpeed` starts at 0 (default float64 value), so after one press of `-`, it becomes 0.01. However, the initial value is never set in `reset()` until line 496, and in `main()` it's set to 0.05. The issue is that `spinSpeed` defaults to 0, so the first `-` press clamps it to 0.01, but then the first `+` press makes it 0.03. This is inconsistent behavior.
2. **[suggestion]** line 487 — Using `math.Max` requires importing `math`, but `math` is already imported, so this is fine. However, a simple conditional or `if spinSpeed < 0.01 { spinSpeed = 0.01 }` would be clearer.
3. **[nit]** line 483-488 — The key bindings for `+` and `-` should also handle `Shift+-` (which produces `+` on most keyboards) and the equals sign for consistency. The code handles `=` separately, which is good.

**Summary:** The PR has a logic issue where `spinSpeed` starts at 0 but gets clamped to 0.01 on first decrease, making the initial behavior inconsistent. The default value should be set before any key handling occurs.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** line 534 — The `time` import is removed, but `time.Now()` is still used in `randInt()` at line 400. This will cause a compilation error.
2. **[blocking]** line 543 — After removing the `time` import, the old implementation at line 400 no longer exists, but the new `rand.Intn(max)` is correct and will work.

**Summary:** The PR accidentally removes the `time` import while keeping code that uses it. The fix for `randInt()` is correct, but the code won't compile without the import removal being consistent. The `math/rand` package should be seeded with `rand.Seed(time.Now().UnixNano())` in `main()` for proper randomness.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** line 561-563 — The clamping logic `lines = lines[:m.height]` will drop the last lines, not the first. This means the footer (which comes last) gets cut off, not the header. To keep the header visible, the code should drop from the beginning: `lines = lines[len(lines)-m.height:]`.
2. **[suggestion]** The check `m.height > 0` is good defensive programming, but the order of operations means that if `m.height` is 0 or negative, no lines are displayed at all.

**Summary:** The clamping logic removes lines from the end (footer) instead of the beginning (top), which is the opposite of what's intended. The header will still be pushed off-screen by tall content.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** line 582 — The rain character `c` is a copy, not a reference. Modifying `c.y` has no effect on `m.rainPos[i]`. The code should modify `m.rainPos[i].y` directly.
2. **[blocking]** line 583 — Similarly, `c.col = rainbowStyle(...)` modifies the copy, not the original. This means rain colors won't update correctly.
3. **[suggestion]** The fix should use `m.rainPos[i].y += m.rainPos[i].v` instead of using the range variable.

**Summary:** The rain animation code uses value copies in the range loop, so modifications to `c` don't affect the actual rain positions. The code won't work as intended without fixing the loop to modify the slice elements directly.
