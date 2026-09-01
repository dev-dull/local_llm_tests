# Fable-5

Results for locally hosted model `Fable-5`.

## hello-go-bubbletea

> Build a "Hello, World!" app in Go using `charmbracelet/bubbletea`, with creative
> flair (animation, color, interactivity); the one strict requirement is that
> pressing `q` must always quit.

| Run | Builds | Runs | Meets requirements | Score | Notes |
|-----|--------|------|--------------------|-------|-------|
| 1   | ✅     | ✅   | ✅                 | 10/10 | Block-letter rainbow banner with confetti and party mode |
| 2   | ✅     | ✅   | ✅                 | 10/10 | Sine-wave greeting with sparkles, themes, and pause |
| 3   | ✅     | ✅   | ✅                 | 10/10 | Starfield with confetti bursts and 18 languages, runewidth-aware |
| 4   | ✅     | ✅   | ✅                 | 10/10 | Wave greeting over starfield with confetti and pause |

**Average score: 10.0/10**

### Run 1

![run 1](hello-go-bubbletea/1/demo.gif)

The model built a 30 fps TUI that renders "HELLO, WORLD!" in a hand-drawn 5-row
block font, with each column riding a sine wave under a scrolling rainbow, over
falling confetti; space toggles a faster "party mode" and enter/arrows cycle
through 12 languages. It builds and runs cleanly as delivered, `q` (and
ctrl+c) quits, and it prints "Goodbye, World! 👋" on exit. The code is
idiomatic Bubble Tea (immutable model, tick commands, cell-grid renderer) with
a correct hand-rolled HSL→hex converter and no stray dependencies.
Score: 3 + 2 + 3 + 2 = 10/10

### Run 2

![run 2](hello-go-bubbletea/2/demo.gif)

Delivered as a nested `hello-go-bubbletea/` project with its own README: the
greeting rides a rainbow sine wave over twinkling sparkles, with enter cycling
9 languages, tab cycling 4 color themes (rainbow/ocean/sunset/neon), and space
freezing time. Builds and runs cleanly; `q`, esc, and ctrl+c all quit. Notable
quality touch: the grid renderer batches unstyled cells so escape codes aren't
emitted per blank character, and it degrades gracefully on tiny terminals.
Score: 3 + 2 + 3 + 2 = 10/10

### Run 3

![run 3](hello-go-bubbletea/3/demo.gif)

A bouncing rainbow greeting over a twinkling starfield with radial confetti
bursts (enter) and 18 languages including Japanese, Chinese, Arabic, and Thai —
and it correctly pulls in `mattn/go-runewidth` to handle double-width CJK
glyphs so rows stay aligned, the only run in the whole batch to consider this.
Builds and runs cleanly; `q` quits. The canvas `put()` helper that claims the
trailing cell of double-width runes is a genuinely careful detail.
Score: 3 + 2 + 3 + 2 = 10/10

### Run 4

![run 4](hello-go-bubbletea/4/demo.gif)

"hello-wave": the greeting rides a rainbow sine wave over a starfield, with
space firing gravity-affected confetti, ←/→ cycling 8 languages, and p pausing.
Builds and runs cleanly; `q`, esc, and ctrl+c quit. The renderer batches runs
of identically-styled cells before invoking lipgloss — an efficiency touch the
comments call out — and the code is clean throughout. Conceptually it overlaps
heavily with runs 2 and 3 (wave + starfield + confetti), so the four runs are
polished variations on one theme rather than four distinct ideas.
Score: 3 + 2 + 3 + 2 = 10/10

## pr-review

> Review six PRs against a Bubble Tea app (two real fixes, two working-but-
> suboptimal, two that compile and run with genuine bugs) and write a
> structured verdict-plus-findings review to `results.md`.

*(Disclosure: the pr-review corpus was authored with Fable-5 assistance, so a
same-family advantage is possible in these runs.)*

| Run | Verdicts OK | Caught PR 3 | Caught PR 6 | Score | Notes |
|-----|-------------|-------------|-------------|-------|-------|
| 1   | ✅          | ✅          | ✅          | 10/10 | All catches, precise line-level reasoning |
| 2   | ✅          | ✅          | ✅          | 10/10 | All catches; one borderline tea.Tick drift claim |
| 3   | ✅          | ✅          | ✅          | 10/10 | All catches, exact overflow arithmetic |
| 4   | ✅          | ✅          | ✅          | 10/10 | All catches, clean and concise |

**Average score: 10.0/10**

### Run 1

Every verdict lands in the acceptable range, both planted bugs are caught with
the mechanism named and the one-line fix given (`m.spun += m.spinSpeed`; index
the slice instead of mutating the `range` copy), and both improvement catches
are made — `tea.Tick` suggested on PR 1, and PR 5's clamp shown to always cut
the footer plus two body rows. Extras are accurate too: the one-frame startup
delay from PR 1's sleep, and a correct Go 1.20+ seeding analysis on PR 4 that
avoids the `rand.Seed` trap.
Score: 3 + 4 + 2 + 1 = 10/10

### Run 2

Same complete result: all verdicts defensible (PR 4 gets APPROVE WITH
SUGGESTIONS on valid, hedged version-dependent notes), both blocking bugs
named with fixes, both improvements caught. The only blemish is a nit claiming
`tea.Tick` "truncates to the interval boundary" — behavior that belongs to
`tea.Every` — which is too minor and tangential to count as a significant
hallucination.
Score: 3 + 4 + 2 + 1 = 10/10

### Run 3

All verdicts in range and all four key catches made, with notably exact
arithmetic on PR 5 (view is always `m.height + 3` lines, so the clamp always
cuts the footer plus the bottom two content rows) and a correct observation
that PR 4's whole init loop runs within one millisecond, which is precisely
why the old `randInt` produced lockstep values.
Score: 3 + 4 + 2 + 1 = 10/10

### Run 4

The most concise of the four, and still complete: both bugs caught with
mechanisms and fixes, both improvements flagged, no hallucinated findings, and
PR 2 approved cleanly with no invented caveats.
Score: 3 + 4 + 2 + 1 = 10/10

## Overall

**Overall average score: 10.0/10** (hello-go-bubbletea 10.0, pr-review 10.0)
