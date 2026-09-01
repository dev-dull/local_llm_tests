# Answer key: pr-review

**Never include this file in a model prompt.** It is for the human running the
tests and for the Claude session generating result reports. The model under
test sees only `prompt.md`; save its response as `review.md` in the run's
numbered directory.

The base program is `results/Qwen3.6-35B-A3B/hello-go-bubbletea/1/main.go`
verbatim. All six PRs apply cleanly to it and compile (verified with Go 1.22
against the run's `go.mod`).

## The PRs

### PR 1: Fix 100% CPU usage during animation — *works, but could be better*

Adds `time.Sleep(33 * time.Millisecond)` inside the ticker command's closure.
This genuinely fixes the busy loop (commands run in their own goroutine, so
the sleep doesn't block input) and caps the frame rate near 30 fps.

- Expected verdict: **APPROVE WITH SUGGESTIONS** (a plain APPROVE with the
  suggestion buried in findings is acceptable; REQUEST CHANGES is too harsh).
- The key insight to credit: Bubble Tea has a purpose-built API for this —
  `tea.Tick` (or `tea.Every`) — and sleeping in a command is the non-idiomatic
  workaround. Secondary credit: frame timing drifts because the sleep doesn't
  account for processing time.
- A review that claims the sleep blocks keyboard input or the UI is **wrong**
  (hallucinated finding).

### PR 2: Use tea.Tick for frame scheduling — *real fix*

Replaces the immediate-return command with
`tea.Tick(time.Second/30, func(time.Time) tea.Msg { return tickerMsg{} })`.
This is the idiomatic fix for the same busy-loop issue.

- Expected verdict: **APPROVE** (or APPROVE WITH SUGGESTIONS with only minor
  notes, e.g. `tea.Tick` aligns to wall-clock intervals so frames can drift
  slightly vs. `tea.Every`; naming a constant for the frame rate).
- REQUEST CHANGES here is a scoring miss — there is no defect.

### PR 3: Add +/- keys to control spin speed — *runs, but genuinely broken*

Adds a `spinSpeed` field, `+`/`=`/`-` key handlers that mutate it, reset and
initialization to 0.05, and updated help text. **The feature is dead code:
the spin advance in `Update` is still the hardcoded `m.spun += 0.05` — the
`spinSpeed` field is written everywhere and read nowhere.** The keys visibly
do nothing.

- Expected verdict: **REQUEST CHANGES**.
- Must catch: `spinSpeed` never affects `m.spun`; the fix is
  `m.spun += m.spinSpeed`.
- Nice extra: `=`'s handler means unshifted `+` on most layouts also works —
  that's intentional, not a bug.

### PR 4: Fix correlated random values in rain initialization — *real fix*

Replaces the millisecond-modulo `randInt` with `rand.Intn` from `math/rand`,
and drops the now-unused `time` import.

- Expected verdict: **APPROVE**.
- Acceptable minor notes: `rand.Intn` panics for `max <= 0` (both existing
  call sites pass constants > 0, so not blocking); global-rand seeding
  concerns are obsolete as of Go 1.20+ (the module targets Go 1.22), so a
  review demanding `rand.Seed` is **wrong**.
- Note the `time` import removal is *required* (it would otherwise be an
  unused import and fail to compile) — a review calling the import removal a
  mistake is wrong.

### PR 5: Fix header being pushed off-screen by tall views — *works, but could be better*

Truncates `lines` to `m.height` at the end of `View`. The header genuinely
stays visible now, but every animation mode emits `m.height` body rows, so the
output is always `m.height + 3` lines — **the truncation silently drops the
footer separator and the last two body rows, every frame**. The right fix is
to size the body to the space actually available (`m.height - 3`).

- Expected verdict: **APPROVE WITH SUGGESTIONS** or **REQUEST CHANGES** —
  either is defensible *only if* the review identifies that the footer (and
  bottom body rows) are now cut off; an APPROVE, or a verdict without that
  observation, is a miss.

### PR 6: Move rain state updates out of View() and into Update() — *runs, but genuinely broken*

Moves the rain-position advance from `View` into `Update`'s tick handling —
but the new loop is `for i, c := range m.rainPos { c.y += c.v; ... }`, which
mutates the **range copy**, not the slice element. Positions and colors never
change: rain characters stay at their initial negative `y` forever and Rain
mode renders a blank screen. Compiles and runs without errors.

- Expected verdict: **REQUEST CHANGES**.
- Must catch: the range-variable copy bug (fix: index with `m.rainPos[i]`, or
  take a pointer).
- Nice extras: the move also drops the gating `m.mode == modeRain` semantics
  change (positions previously advanced only while Rain mode was displayed —
  behavior here is equivalent, so this is a nit at most).

## Scoring a run (out of 10)

Score the model's `review.md` with this rubric **instead of** the generic
builds/runs rubric:

- **Verdicts (0–3):** start at 3; −1 for each PR whose verdict is outside the
  acceptable range listed above (floor 0).
- **Blocking catches (0–4):** +2 for identifying the dead `spinSpeed` field in
  PR 3; +2 for identifying the range-copy bug in PR 6. Credit requires naming
  the actual mechanism, not just "this looks wrong".
- **Improvement catches (0–2):** +1 for suggesting `tea.Tick`/`tea.Every` on
  PR 1; +1 for spotting the dropped footer/body rows on PR 5.
- **Precision (0–1):** 1 if the review contains no significant hallucinated
  findings (e.g. the wrong claims called out above, invented APIs, nonexistent
  compile errors); 0 otherwise.

Report the breakdown per run, e.g. `Score: 3 + 4 + 1 + 1 = 9/10`. In the
report's summary table, replace the Builds/Runs/Meets requirements columns
with: `| Run | Verdicts OK | Caught PR 3 | Caught PR 6 | Score | Notes |`.
A run that ignores the required output format still gets scored on substance,
but note the deviation in prose. No GIFs for this test.
