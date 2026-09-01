# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion]** `tickerCmd()`, inner return func — `time.Sleep(33 * time.Millisecond)` inside the closure does throttle the loop to ~30 fps and fixes the busy-spin. However, the idiomatic Bubble Tea approach is `tea.Tick` (which PR 2 demonstrates). A raw sleep in a `tea.Cmd` works but bypasses Bubble Tea's internal tick bookkeeping, which can cause frame-skipping or dropped ticks if the render pipeline lags.

**Summary:** The CPU fix is correct and sufficient for the stated goal. Prefer `tea.Tick` (PR 2) for a cleaner, framework-native solution, but this change is mergeable as-is.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. **none**

**Summary:** Replaces the raw `time.Sleep` with `tea.Tick(time.Second/30, ...)` — the idiomatic Bubble Tea mechanism for throttling animation frames. Correct, minimal, and directly addresses the CPU-spinning problem.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** `Update`, `"+"` / `"-"` cases (lines +87–91 of diff) — The PR adds a `spinSpeed` field and key handlers, but the key handlers never affect the animation. In `Update`, the tick branch unconditionally does `m.spun += 0.05` (line +107 of base). The PR never uses `m.spinSpeed` there, so the animation always advances at the hardcoded rate regardless of user input. The feature is non-functional as written.

**Summary:** The scaffolding (new field, key bindings, reset init, help text) is well-structured, but the core logic is missing: `m.spun += 0.05` needs to become `m.spun += m.spinSpeed` for the feature to work.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** APPROVE

**Findings:**
1. **[suggestion]** `randInt` and imports — `rand.Intn` works correctly in Go 1.20+ (the auto-seeded global `rand` package is safe to use directly). In Go 1.22+ the idiomatic form is `rand.Default.Intn(max)` to be explicit. This is a minor style note; the PR correctly replaces the wall-clock modulo with a real PRNG and removes the now-unused `time` import.

**Summary:** The PR correctly fixes the correlation problem. Using the global `math/rand` is fine for a small demo app — no seed, no concurrency concerns.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** `View`, lines after footer render — The guard checks `len(lines) > m.height` but then truncates with `lines = lines[:m.height]` which drops the footer entirely (the footer is the last element appended to `lines`). The header is saved but the help footer is lost. More importantly, the truncation happens *after* `lines` already has the footer appended, so the first `m.height` lines include the header but exclude the footer.

**Summary:** The intent (clamp output to window height) is correct, but truncating `lines[:m.height]` silently drops the footer. A better approach would be to cap the content area before building it, or to preserve the last two lines (header + footer) and truncate the middle if needed.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking]** `Update`, tick branch, `for i, c := range m.rainPos` — The range loop copies each `rainChar` into `c`. Mutating `c.y` and `c.col` changes only the copy; `m.rainPos[i]` is never updated. The rain will appear stationary forever because the model state is never mutated. Must use `for i := range m.rainPos` and index into `m.rainPos[i]` to mutate in place (matching the original View code's working approach).

**Summary:** The architectural intent — moving stateful updates out of `View` — is correct and the right direction. However, the copy-semantics of `range` make the change non-functional as written. Fixing the loop to index into `m.rainPos` directly resolves the issue.
