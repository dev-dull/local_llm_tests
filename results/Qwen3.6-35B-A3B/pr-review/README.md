# Qwen3.6-35B-A3B — pr-review

Results for locally hosted model `Qwen3.6-35B-A3B` on the `pr-review` test.

> Review six PRs against a Bubble Tea app (two real fixes, two working-but-
> suboptimal, two that compile and run with genuine bugs) and write a
> structured verdict-plus-findings review to `results.md`.

| Run | Verdicts OK | Caught PR 3 | Caught PR 6 | Score | Notes |
|-----|-------------|-------------|-------------|-------|-------|
| 1   | ❌          | ✅          | ✅          | 8/10  | Both bugs caught; demands `rand.Seed` on PR 4 |
| 2   | ❌          | ❌          | ✅          | 6/10  | Approves PR 3's dead feature outright |
| 3   | ✅          | ✅          | ✅          | 9/10  | Best run; invents a `rand.Default.Intn` API |
| 4   | ❌          | ✅          | ❌          | 5/10  | Falls for both anticipated wrong claims |

**Average score: 7.0/10**

## Run 1

A strong review: both planted bugs are caught with the mechanism named and the
right fix, and both improvement catches land (`tea.Tick` on PR 1, the dropped
footer on PR 5). The one verdict miss is PR 4, rejected for lacking an
explicit `rand.Seed` — unnecessary on the auto-seeded Go 1.20+ toolchain.
Precision suffers from side claims: that Wave/Pulse render "exactly four
lines" and never overflow (all modes emit `m.height` rows), a bogus warning
that the `range` loop's shadowed `i` "could cause a compile error", and a
wrong assertion that `main`'s initializer is redundant because reset runs
first.
Score: 2 + 4 + 2 + 0 = 8/10

## Run 2

Catches PR 6's range-copy bug cleanly, but approves PR 3 outright — even
praising its "well-defined starting value" — missing that `spinSpeed` is
never read and the feature is dead. On PR 5 it correctly notes the footer is
dropped (earning the catch) but then claims the header "may still be pushed
off-screen", which is the one problem the clamp actually does fix.
Score: 2 + 2 + 2 + 0 = 6/10

## Run 3

The best local-model run of the batch: every verdict in the acceptable range,
both bugs caught with fixes, both improvements flagged. The only cost is
precision: it recommends a nonexistent `rand.Default.Intn(max)` as "the
idiomatic form" on Go 1.22+, and claims a raw sleep "bypasses Bubble Tea's
tick bookkeeping", which is invented framework behavior.
Score: 3 + 4 + 2 + 0 = 9/10

## Run 4

Falls into both traps the corpus plants: PR 1 is rejected with the claim that
`time.Sleep` in a command "freezes the entire terminal" (commands run in their
own goroutine — the exact hallucination the answer key anticipates), and PR 6
is approved by inventing a write-back (`m.rainPos[i] = c`) that appears
nowhere in the diff, thereby missing the range-copy bug. PR 3's dead field and
PR 5's footer loss are still caught.
Score: 1 + 2 + 2 + 0 = 5/10
