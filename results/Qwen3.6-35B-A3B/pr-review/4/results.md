# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking] tickerCmd() — time.Sleep blocks the Bubble Tea event loop** — The `tea.Cmd` returned by `Init()` is a `func() tea.Msg`. Bubble Tea calls this function synchronously on its main goroutine; any `time.Sleep(33 * time.Millisecond)` inside it freezes the entire terminal, blocking key presses, window resizes, and quit signals for 33 ms every frame. The app would feel sluggish or unresponsive. Use `tea.Tick` instead (already shown in PR 2).

**Summary:** The goal (eliminate busy-spin) is correct, but `time.Sleep` in a `tea.Cmd` is a blocking call that freezes the TUI. PR 2 shows the proper fix via `tea.Tick`.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. **[nit] time.Second/30 is integer division yielding 33 ms, not 33⅓ ms** — Go integer division truncates, so the actual frame rate is ~30.3 fps rather than exactly 30. This is cosmetically irrelevant, but `33*time.Millisecond + 11*time.Microsecond` (or `33*time.Millisecond` with a comment noting the approximation) would be clearer.

**Summary:** `tea.Tick` is the idiomatic Bubble Tea mechanism for periodic animation frames. It yields control between ticks so input is responsive. The integer-division truncation is harmless.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking] Update() tick branch — `m.spun += 0.05` is never changed to use `m.spinSpeed`** (line ~3 of the `case " ":` block / the `tickerMsg` branch at base line 136) — The PR adds the `spinSpeed` field, key handlers, and reset logic, but leaves the hardcoded `0.05` increment untouched. The `spinSpeed` field is initialized and clamped but never read, so +/- keys have zero effect on the animation.
2. **[suggestion] Update() — `+` and `=` keys increase without a ceiling** — `m.spinSpeed += 0.02` has no upper bound. Users can set arbitrarily high spin speeds, producing a blur. Consider clamping to a maximum (e.g., `m.spinSpeed = min(m.spinSpeed+0.02, 0.5)`).

**Summary:** The new field, key handlers, and reset logic are correct in isolation, but the PR is incomplete: the tick branch still uses the hardcoded `0.05` instead of `m.spinSpeed`, so the feature doesn't work. Fixing the increment to `m.spun += m.spinSpeed` is all that's needed.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion] No explicit seeding of `math/rand`** — The PR replaces the hand-rolled `randInt` with `rand.Intn(max)` but doesn't call `rand.Seed`. In Go 1.x the global `math/rand` auto-seeds from `time.Now()` on first use, so this works in practice, but an explicit `rand.Seed(time.Now().UnixNano())` at init time or `rand.New(rand.NewSource(...))` would be more robust and portable across Go versions.
2. **[nit] Removing `"time"` from imports but `time` is still used** — Wait, `time` is no longer used after removing the `randInt` body, so the import removal is correct. Good hygiene.

**Summary:** The change correctly fixes the correlation problem by using a proper PRNG. Add an explicit seed call for robustness.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking] Simple truncation discards footer and separator on small terminals** (diff at base line ~149) — When `len(lines) > m.height` the code slices to `lines[:m.height]`, which may cut off the separator and footer entirely. A terminal with fewer rows than the minimum content (title + blank + animation line + separator + footer = at least ~5 rows) would show a header with no context. The fix should check whether the terminal is *large enough* for any meaningful output and return early if not.
2. **[suggestion] Compute minimum required lines and return early when `m.height < minRequired`** — Calculate the minimum rows the app needs (title bar, blank line, at least one animation row, separator, footer ≈ 5) and bail out with a message like "Terminal too small" instead of silently cropping content.

**Summary:** The intent is right, but truncating lines is too aggressive — it can remove the separator and footer. A size guard with an early return is a cleaner solution.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** APPROVE

**Findings:**
1. **[nit] Range-over-copy loop vs index loop** — `for i, c := range m.rainPos` iterates over copies of each `rainChar`; the assignments (`c.y += c.v`, etc.) modify those copies, which are then written back with `m.rainPos[i] = c`. This is correct and idiomatic Go, but using the index-based `for i := range m.rainPos { m.rainPos[i].y += ... }` (which was the original View code) is slightly more direct. Not a problem, just a style note.

**Summary:** Moving side-effectful state updates from `View` (which should be pure) to `Update` is the right architectural fix for Bubble Tea. The logic is identical to the original, just in the correct place.
