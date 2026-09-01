# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** APPROVE

**Findings:**
1. **[nit]** line 425 — The comment says "Cap the frame rate at roughly 30 fps" but `time.Sleep(33 * time.Millisecond)` gives approximately 30.3 fps (1000/33 ≈ 30.3). For truly ~30 fps, use `33*time.Millisecond` is acceptable but `34*time.Millisecond` would give 29.4 fps. This is a minor nitpick.

**Summary:** The PR correctly addresses the CPU usage issue by adding a sleep to prevent busy-spinning. The implementation works and will cap the frame rate at approximately 30 fps as claimed.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. **[suggestion]** line 443 — The diff shows `time.Second/30` which is correct for ~30 fps. However, there's a potential issue: `time.Second/30` evaluates to `33333333` nanoseconds, which is exactly 33.3...ms. This is fine, but note that `tea.Tick` schedules based on wall-clock time and may accumulate drift over time. For a more robust solution, consider using `time.NewTicker` or accepting that `tea.Tick` handles this internally.

**Summary:** This PR provides a proper solution using Bubble Tea's built-in `tea.Tick` mechanism, which is cleaner and more idiomatic than manual sleep. It achieves the same goal as PR 1 but in a more principled way using the framework's intended API.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** line 487 — The spin speed decrement uses `math.Max(0.01, m.spinSpeed-0.02)`, but `math.Max` returns a `float64`. Since `spinSpeed` is a `float64`, this should work, but the comparison `0.01` as a float literal should be `0.01` (which it is). However, there's a potential issue: when `spinSpeed` is already at or below 0.01, the user cannot reset it to the default of 0.05. A better approach would be to use `math.Max(0.01, m.spinSpeed-0.02)` but this actually does clamp correctly. The real issue is on line 484-485: the `+=` operation has no upper bound, so spinSpeed can grow unboundedly large.
2. **[suggestion]** line 483-488 — The PR adds `+` and `-` key handling but the `spinSpeed` is only used in the `renderSpin` function (line 364 in the diff shows `m.spun += 0.05` but this should use `m.spinSpeed` instead). Looking at line 364 in the original code: `m.spun += 0.05` — this hard-coded value should be `m.spinSpeed`.

**Summary:** The PR has a critical bug: the spin speed update in `renderSpin` still uses the hardcoded `0.05` instead of `m.spinSpeed`, making the +/- keys ineffective. Additionally, there's no upper bound on increasing spin speed. The PR also fails to mention that the `renderSpin` function needs to be updated to use the new field.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** line 531-534 — The diff removes `time` import but the `randInt` function is still called from `main.go` line 400 (original) which uses `time.Now().UnixNano()`. Wait, looking more carefully at the diff: the PR changes the implementation of `randInt` at line 542-543 to use `rand.Intn(max)`, which is correct. However, the `time` import is removed entirely, but `time` is still used elsewhere in the file (e.g., line 61 in imports, and `tickerCmd` uses time implicitly via the ticker). Actually, looking at the imports: `time` is imported at line 61, and the PR removes it at line 534. But `time` is used in `tickerCmd` (line 378-381 in original) for... actually no, `tickerCmd` doesn't use `time` directly. Let me re-check: the only use of `time` in the original is in `randInt`. So this change is actually correct.
2. **[blocking]** line 531 — The import of `math/rand` is added, but `math/rand` was not imported before. However, the original `time` import is removed. The issue is that `math/rand` needs to be seeded. Without calling `rand.Seed()`, the PRNG will produce the same sequence every run. The original code, while correlated across calls in the same program run, at least varied between runs due to wall-clock dependency. With `math/rand` without seeding, it produces the same sequence on every run.

**Summary:** The PR correctly identifies the correlation issue but introduces a new problem: `math/rand` produces the same sequence on every program run without seeding. To fix this, either use `math/rand` with a seed (e.g., `rand.Seed(time.Now().UnixNano())` in `main()` or use `crypto/rand` for cryptographically secure randomness (overkill here) or stick with a better time-based PRNG. A better fix would be to use `math/rand` with proper seeding.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** line 561-563 — The clamping logic `lines = lines[:m.height]` truncates the output, but this means the footer (with the separator) may be cut off. The header should be preserved, not just any arbitrary lines. A better approach would be to trim blank lines from the bottom first, or to ensure the header and footer always appear. For example, if `m.height = 20` and `lines` has 25 lines, truncating to 20 would lose the footer. The header (info bar) and footer (separator) should be preserved.
2. **[suggestion]** The current approach is a blunt instrument. A better solution would be to:
   - Keep the header (lines 0-1: info bar and blank)
   - Keep the footer (the separator line)
   - Only truncate content lines from the middle/bottom
   - Or use `lines = append(lines[:m.height-1], lines[len(lines)-1])` to keep first n-1 lines plus the last line (footer)

**Summary:** The PR's approach of blindly truncating to `m.height` lines will likely cut off important content including the footer. The header preservation goal is correct, but the implementation is too aggressive and will remove the footer and potentially useful information.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** line 583 — In the rain update code, `c.col = rainbowStyle(float64(i) / float64(len(m.rainPos)))` recalculates the color for each rain character every frame. This means colors will cycle even though the characters aren't moving through different positions. The original code calculated color based on index only, which gave each character a fixed color. The recalculation here is actually a change in behavior — colors will now cycle over time even for stationary characters. This might be intentional as a visual effect, but it's worth noting.

**Summary:** This PR correctly moves the rain animation logic from `View()` to `Update()`, which is the proper separation of concerns for a Bubble Tea application (View should be pure). The change fixes a fundamental architectural issue where the view function was mutating state. The one minor consideration is that colors now cycle each frame, which is a visual change from the original behavior.
