# Results

Evaluation reports for locally hosted models, one report per model per test
in `<model>/<test-name>/README.md`. Scoring rubrics live in `.rubric/`; the
report format and procedure are defined in
[report-generation-prompt.md](report-generation-prompt.md).

| Model | Test reports | Overall average |
|-------|--------------|-----------------|
| Fable-5 | [hello-go-bubbletea](Fable-5/hello-go-bubbletea/README.md) (10.0), [pr-review](Fable-5/pr-review/README.md) (10.0) | 10.0/10 |
| Qwen3.6-35B-A3B | [hello-go-bubbletea](Qwen3.6-35B-A3B/hello-go-bubbletea/README.md) (7.3), [pr-review](Qwen3.6-35B-A3B/pr-review/README.md) (7.0) | 7.2/10 |
| Qwen3-Coder-Next-UD-Q4_K_M | [hello-go-bubbletea](Qwen3-Coder-Next-UD-Q4_K_M/hello-go-bubbletea/README.md) (5.3), [pr-review](Qwen3-Coder-Next-UD-Q4_K_M/pr-review/README.md) (4.8) | 5.1/10 |

*(Fable-5's pr-review score carries a disclosure: the review corpus was
authored with Fable-5 assistance, so a same-family advantage is possible.)*

## Conclusions: Qwen3-Coder-Next vs Qwen3.6-35B

Is Qwen3-Coder-Next actually better than Qwen3.6-35B for typical coding
tasks? On this evidence, **no — Qwen3.6-35B wins both with and without
creativity factored in**, and the gap is widest when creativity is factored
*out*. The pr-review test, added later, independently confirms the ordering
on a task with no creativity axis at all.

### Factoring creativity in: Qwen3.6-35B wins clearly (7.3 vs 5.3)

Qwen3.6 aimed higher in every run — four switchable animation modes, a
starfield with a typewriter reveal and wind control, a steerable bouncing
bubble with particle bursts — and still scored two points higher.
Qwen3-Coder's designs were the blandest of the three models tested: a loading
spinner, a static bordered card, a box that was *supposed* to orbit the
screen but never moved. There is no "less flashy but more solid" trade-off to
invoke; it lost on both axes.

### Factoring creativity out: Qwen3.6-35B still wins on pure engineering

Judged only on build/run/requirement mechanics:

- **Hard failures.** Qwen3-Coder: 2 of 4 runs (run 2 doesn't build; run 1
  fails the sole strict requirement — `q`, and even ctrl+c, cannot terminate
  it). Qwen3.6: 1 of 4 (run 3's startup panic).
- **Hallucination — the key differentiator.** Qwen3-Coder invented three
  nonexistent Go modules in a `go.mod`, a nonexistent `lipgloss.FontSize`
  API, and mis-called `tea.Tick`. Qwen3.6 hallucinated nothing across four
  runs: every dependency real, every API call valid.
- **Bug class.** Qwen3-Coder's bugs are framework-semantics failures:
  mutating a model copy inside a closure so its animation never runs,
  fabricating `spinner.TickMsg` values by hand, a quit path that hangs the
  event loop. Qwen3.6's bugs are ordinary logic slips — a negative
  `strings.Repeat` before the first resize, a delay-free tick command that
  busy-loops the CPU, sloppy centering math — the kind a human catches in one
  debug pass.

For typical coding tasks, "does it build, does the API exist, does the event
loop behave" matters far more than flair, and that is precisely where
Qwen3-Coder-Next was weakest: all four of its runs contained at least one
significant correctness bug, versus roughly one and a half for Qwen3.6.

### pr-review: the gap holds on code review too (7.0 vs 4.8)

The review test — six PRs against a buggy Bubble Tea app, two of them
containing planted bugs that compile and run — removes creativity from the
equation entirely, and Qwen3.6-35B still wins. It caught the PR 6 range-copy
bug in 3 of 4 runs and produced the best local review of the batch (9/10,
every verdict defensible). Qwen3-Coder-Next caught it in only 2 of 4 and
**approved the broken refactor twice** — once with no findings at all, once
while hallucinating a behavior change that doesn't exist — and repeatedly
invented compile errors in PRs it was told compile.

The models share one failure: neither ever earned the precision point.
Both demand `rand.Seed` where Go 1.20+ needs none, and Qwen3.6's worst run
reproduced verbatim the "sleep blocks the whole UI" misconception the answer
key anticipates. Local-model reviews here flag real bugs at best two-thirds
of the time and pad every review with at least one confident falsehood.

### Caveats

Small sample: four runs of a single, fairly niche task (Bubble Tea TUIs
stress knowledge of one library's idioms). The UD-Q4_K_M quantization may be
degrading the Coder model disproportionately. And Qwen3.6's one crash came in
its most ambitious run — its failure rate partly bought more attempted
functionality.
