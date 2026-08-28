# Results

Evaluation reports for locally hosted models, one report per model directory.
Scoring rubric and report format are defined in [prompt.md](prompt.md).

| Model | Tests | Average score | Report |
|-------|-------|---------------|--------|
| Fable-5 | hello-go-bubbletea | 10.0/10 | [Fable-5/README.md](Fable-5/README.md) |
| Qwen3.6-35B-A3B | hello-go-bubbletea | 7.3/10 | [Qwen3.6-35B-A3B/README.md](Qwen3.6-35B-A3B/README.md) |
| Qwen3-Coder-Next-UD-Q4_K_M | hello-go-bubbletea | 5.3/10 | [Qwen3-Coder-Next-UD-Q4_K_M/README.md](Qwen3-Coder-Next-UD-Q4_K_M/README.md) |

## Conclusions: Qwen3-Coder-Next vs Qwen3.6-35B

Is Qwen3-Coder-Next actually better than Qwen3.6-35B for typical coding
tasks? On this evidence, **no — Qwen3.6-35B wins both with and without
creativity factored in**, and the gap is widest when creativity is factored
*out*.

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

### Caveats

Small sample: four runs of a single, fairly niche task (Bubble Tea TUIs
stress knowledge of one library's idioms). The UD-Q4_K_M quantization may be
degrading the Coder model disproportionately. And Qwen3.6's one crash came in
its most ambitious run — its failure rate partly bought more attempted
functionality.
