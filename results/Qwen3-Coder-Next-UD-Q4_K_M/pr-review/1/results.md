# PR review

## PR 1: Fix 100% CPU usage during animation

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking] tickerCmd() — the sleep is in the wrong place** — The `time.Sleep(33 * time.Millisecond)` is inside the function that returns the message, but it doesn't block the cmd from returning. The cmd returns immediately with `tickerMsg{}` after sleeping, but the next tick is scheduled without any delay mechanism. This doesn't actually cap the frame rate; it just adds unnecessary delay before each message is sent without preventing the tight loop.

**Summary:** The fix doesn't achieve its stated goal. The sleep happens before returning the message, but `tickerCmd()` is called again immediately after receiving the message, creating the same busy loop pattern. This should use `tea.Tick` like PR 2 does, or add a proper delay mechanism.

## PR 2: Use tea.Tick for frame scheduling

**Verdict:** APPROVE

**Findings:**
1. **[nit] main.go:343 — should import `time`** — The PR uses `time.Second/30` but doesn't import the `time` package. This will cause a compilation error.

**Summary:** The approach of using `tea.Tick` is correct and follows Bubble Tea's recommended pattern for frame scheduling. The only issue is the missing `time` import, which is an easy fix.

## PR 3: Add +/- keys to control spin speed

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking] main.go:373 — spinSpeed is never used in renderSpin()** — The PR adds `spinSpeed` as a field and allows it to be modified via +/- keys, but the `renderSpin()` function still hardcodes `0.05` in the line `m.spun += 0.05`. The actual spin advancement should use `m.spinSpeed` instead.
2. **[blocking] main.go:138 — update for modeSpin missing spinSpeed** — In the Update function's ticker handling, `m.spun += 0.05` should be `m.spun += m.spinSpeed` for the spin mode to respect the user's speed setting.

**Summary:** The PR introduces spin speed control but fails to actually use the new field in the animation logic. The renderSpin function and Update ticker handler both hardcode the value 0.05, making the +/- keys ineffective.

## PR 4: Fix correlated random values in rain initialization

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking] main.go:352 — missing `rand.Seed()` call** — Using `math/rand.Intn()` without seeding will produce the same sequence of "random" numbers on every run. While better than the wall-clock based approach, it will still produce identical rain patterns across program restarts.
2. **[suggestion] main.go:352 — consider using `math/rand/v2`** — The newer `math/rand/v2` package seeds itself automatically and is preferred for new code.

**Summary:** The switch to `math/rand` is an improvement over the time-based PRNG, but without seeding the random generator, all runs will produce identical rain patterns. The PR fixes the correlation issue within a single run but doesn't ensure variety across runs.

## PR 5: Fix header being pushed off-screen by tall views

**Verdict:** REQUEST CHANGES

**Findings:**
1. **[blocking] main.go:159 — clamping logic is incorrect** — The clamping `lines = lines[:m.height]` happens AFTER the footer is appended, but the header is at the BEGINNING of the `lines` slice. When lines are sliced to `m.height`, the footer is kept and the header gets trimmed off, which is the opposite of what's intended.
2. **[suggestion] Should trim empty lines from end first** — Before clamping, the code should remove trailing empty strings that are just padding, not content.

**Summary:** The clamping logic trims from the end of the slice, but the header is at the beginning. This will keep the footer and cut off the header, which is exactly the problem the PR claims to solve. The fix should either trim from the beginning when exceeding height, or trim trailing empty lines first.

## PR 6: Move rain state updates out of View() and into Update()

**Verdict:** APPROVE WITH SUGGESTIONS

**Findings:**
1. **[suggestion] main.go:582 — rain updates use shadow copy** — The loop `for i, c := range m.rainPos` iterates over copies of the struct values. While the code modifies `c.y`, this only modifies the copy, not the original in the slice. The line `c.y += c.v` should be `m.rainPos[i].y += m.rainPos[i].v` or the loop should use `for i := range m.rainPos` and access via index.
2. **[nit] main.go:583 — color update happens every tick but uses same formula** — The line `c.col = rainbowStyle(float64(i) / float64(len(m.rainPos)))` recalculates the same color every frame unnecessarily since `i` and `len(m.rainPos)` don't change.

**Summary:** The architectural improvement of moving rain updates from View() to Update() is correct and follows Bubble Tea best practices. However, there's a subtle bug in the range loop - iterating over copies means the position updates won't persist to the actual slice elements.
