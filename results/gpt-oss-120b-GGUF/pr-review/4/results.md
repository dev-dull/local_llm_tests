# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** APPROVE

**Findings:**
1. none

**Summary:** The added `time.Sleep(33 * time.Millisecond)` in `tickerCmd` caps the frame rate at roughly 30 fps, effectively reducing CPU usage without affecting functionality.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. none

**Summary:** Replacing the manual ticker with `tea.Tick(time.Second/30, ...)` cleanly schedules frames at 30 fps using Bubble Tea's built‑in mechanism, achieving the same CPU‑saving goal as PR 1.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** `model.spinSpeed` is introduced and modified by the new key bindings, but it is never used in the spin rendering logic, so the feature does not actually affect spin speed.

**Summary:** The key handling and UI updates are in place, but without integrating `spinSpeed` into the animation (e.g., incrementing `m.spun` by `m.spinSpeed`), the change has no observable effect. Update the render logic to use `spinSpeed` and reset it appropriately.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** APPROVE

**Findings:**
1. none

**Summary:** Switching `randInt` to use `math/rand.Intn` removes correlation between successive random values, improving the visual randomness of the rain animation.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** APPROVE

**Findings:**
1. none

**Summary:** Adding a height check and truncating `lines` to `m.height` ensures the header remains visible even when the view generates many lines.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** APPROVE

**Findings:**
1. none

**Summary:** Relocating rain position updates to the `Update` function restores the pure‑function contract of `View()` and keeps the animation logic correctly placed in the update loop.
