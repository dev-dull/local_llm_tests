# Fable-5 — pr-review

Results for locally hosted model `Fable-5` on the `pr-review` test.

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

## Run 1

Every verdict lands in the acceptable range, both planted bugs are caught with
the mechanism named and the one-line fix given (`m.spun += m.spinSpeed`; index
the slice instead of mutating the `range` copy), and both improvement catches
are made — `tea.Tick` suggested on PR 1, and PR 5's clamp shown to always cut
the footer plus two body rows. Extras are accurate too: the one-frame startup
delay from PR 1's sleep, and a correct Go 1.20+ seeding analysis on PR 4 that
avoids the `rand.Seed` trap.
Score: 3 + 4 + 2 + 1 = 10/10

## Run 2

Same complete result: all verdicts defensible (PR 4 gets APPROVE WITH
SUGGESTIONS on valid, hedged version-dependent notes), both blocking bugs
named with fixes, both improvements caught. The only blemish is a nit claiming
`tea.Tick` "truncates to the interval boundary" — behavior that belongs to
`tea.Every` — which is too minor and tangential to count as a significant
hallucination.
Score: 3 + 4 + 2 + 1 = 10/10

## Run 3

All verdicts in range and all four key catches made, with notably exact
arithmetic on PR 5 (view is always `m.height + 3` lines, so the clamp always
cuts the footer plus the bottom two content rows) and a correct observation
that PR 4's whole init loop runs within one millisecond, which is precisely
why the old `randInt` produced lockstep values.
Score: 3 + 4 + 2 + 1 = 10/10

## Run 4

The most concise of the four, and still complete: both bugs caught with
mechanisms and fixes, both improvements flagged, no hallucinated findings, and
PR 2 approved cleanly with no invented caveats.
Score: 3 + 4 + 2 + 1 = 10/10
