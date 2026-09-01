# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking] tickerCmd() — sleep placed before return, not after** — The `time.Sleep(33 * time.Millisecond)` is placed before `return tickerMsg{}`, which means the sleep happens *before* the message is sent, not between frames. This creates a blocking delay in the command execution but doesn't properly throttle the event loop. The sleep should be placed after the message is returned, or better yet, use `tea.Tick` for proper frame scheduling.

**Summary:** While the intent to cap frame rate is correct, the implementation is flawed — sleeping before returning the message doesn't achieve proper throttling of the animation loop.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. **none**

**Summary:** This is the correct approach — using `tea.Tick(time.Second/30, ...)` properly integrates with Bubble Tea's event loop and will schedule frames at ~30 fps without busy-spinning.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking] Update() handler — spinSpeed not used in renderSpin()** — The PR adds `spinSpeed` field and key handlers to modify it, but `renderSpin()` still uses the hardcoded value `0.05` for advancement. The line `m.spun += 0.05` should use `m.spinSpeed` instead.
2. **[suggestion] Update() — should handle "=" separately from "+"** — The PR handles both "+" and "=" the same way, but "=" typically doesn't require Shift on many keyboards. Consider if this is intentional or if "=" should be handled separately for international keyboards.

**Summary:** The PR has a critical bug where the `spinSpeed` variable is modified but never actually used in the renderSpin function, making the feature non-functional.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[nit] missing `rand.Seed()` call** — While `rand.Intn()` is used correctly, Go's `math/rand` package produces the same sequence by default unless seeded. Add `rand.Seed(time.Now().UnixNano())` at program startup (e.g., in `main()`) to ensure different random sequences on each run.
2. **[nit] `time` import still needed elsewhere** — The `time` import is removed entirely, but it's used elsewhere in the file (e.g., `time.Now()` in `randInt` was replaced, but `time` may still be needed elsewhere). Verify no other code relies on the time package after removal.

**Summary:** Using `math/rand` is the correct fix for the correlated random values. Just ensure the PRNG is properly seeded for different behavior on each run.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking] Clamping logic is incorrect** — The PR clamps `lines` to `m.height`, but this includes the header, separator lines, and footer in the count. This means if the terminal has 24 rows and the output has 26 lines total (including header/footer), it would cut off the last 2 rows of actual content but still potentially push content off. The clamping should only apply to the content area, not the fixed header/footer.
2. **[blocking] Clamping happens after rendering each mode** — The rain mode and pulse mode already generate `m.height` lines. Adding header (2 lines) + separator (1 line) + footer (1 line) = 42 lines for a 38-line terminal, then cutting to 38 would still lose content. The issue is that some modes can generate more than `m.height` lines.

**Summary:** The clamping approach is too simplistic — it clamps the entire output including header/footer rather than just the content area, which could still result in important UI elements being pushed off-screen or the content area being truncated incorrectly.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** APPROVE

**Findings:**
1. **none**

**Summary:** This is correct — moving rain state updates from the pure `View()` function into `Update()`'s tick handling follows the Bubble Tea pattern of having `View()` be a pure function of the model state.
