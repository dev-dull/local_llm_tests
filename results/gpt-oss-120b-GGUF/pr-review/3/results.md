# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** main.go:330 – Adding a fixed `time.Sleep(33 * time.Millisecond)` caps the frame rate but introduces a hard‑coded delay; using Bubble Tea’s `tea.Tick` would be a more idiomatic way to schedule frames.

**Summary:** The change stops the busy‑spin and reduces CPU usage, but a more flexible tick‑based approach would be preferable.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. none

**Summary:** Replaces the manual sleep with Bubble Tea’s `tea.Tick`, providing a clean, configurable frame rate and eliminating the need for a hard‑coded delay.

## PR 3: Add +/- keys to control spin speed

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** main.go:483‑488 – The new `+` and `-` keys adjust `spinSpeed`, but the default speed is still hard‑coded in `reset()` and `main()`. Consider exposing a configurable default via a flag or constant.
2. **[nit]** main.go:462‑476 – Added `spinSpeed` field increases struct size; no functional issue but could be documented.

**Summary:** Introduces useful speed controls and updates help text, but the default speed could be made configurable.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** APPROVE

**Findings:**
1. none

**Summary:** Replaces the deterministic `randInt` implementation with `math/rand.Intn`, fixing correlated rain motion.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** APPROVE

**Findings:**
1. none

**Summary:** Clamps the output to the terminal height, ensuring the header remains visible.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** APPROVE

**Findings:**
1. none

**Summary:** Restores purity of `View()` by moving state mutation to `Update()`, aligning with Bubble Tea best practices.
