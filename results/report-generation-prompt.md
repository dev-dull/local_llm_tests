# Task: Generate result reports for local LLM test runs

You are generating evaluation reports for prompts that were run against locally
hosted LLM models. Work through the directory structure described below and
produce one `README.md` report per model **per test**.

## Directory layout

- `prompts/<test-name>/prompt.md` — the original prompt that was given to each
  model.
- `results/<model>/` — one directory per model (the directory name is the model
  name).
- `results/<model>/<test-name>/<n>/` — numbered directories (`1/`, `2/`, `3/`,
  …), each holding one independent run of that prompt by that model.
- `results/<model>/<test-name>/README.md` — the report you generate: one per
  model per test, covering all of that model's runs of that test.
- `results/.rubric/<test-name>.md` — the scoring rubric for each test (see
  Scoring below). Not a model directory; skip it when iterating models.

## Procedure

For each model directory under `results/`, and each `<test-name>`
subdirectory within it:

1. **If `results/<model>/<test-name>/README.md` already exists, skip that
   test for that model entirely.** Do not read the existing file — assume it
   is complete and move on. A test directory containing only run directories
   with prompt copies (no model output) is scaffolded but not yet executed:
   skip it too, without writing a report.
2. Read the matching `prompts/<test-name>/prompt.md` so you understand what
   the model was asked to do and what the strict requirements are, and the
   test's rubric in `results/.rubric/<test-name>.md` for how to evaluate it.
3. For each numbered run directory, evaluate the result as the rubric
   directs. For build-and-run tests that means:
   - What did the model produce (files, completeness)?
   - Does it build? Actually build it (e.g. `go build` for Go projects) and
     record success or the exact failure.
   - Does it run, and does it satisfy the prompt's strict requirements (e.g.
     for `hello-go-bubbletea`: does `q` quit)?
   - Note anything interesting: creative choices, bugs, hallucinated
     dependencies, deviations from the prompt.
4. Score the run using the scoring rules below.
5. Write the report as `results/<model>/<test-name>/README.md` using the
   template below.

Do not modify anything inside the numbered run directories except as needed to
build (e.g. running `go mod tidy`); never edit the model's source code. If a
run doesn't build or run, report that honestly — do not fix it.

## Scoring

Each run gets a **score out of 10**, and each test gets a **combined average
score** across its runs, so models can be compared at a glance.

Every test has its own rubric in `results/.rubric/<test-name>.md`. Before
scoring a test's runs, read that rubric and score strictly against it. A
rubric may also adjust the report for its test — different summary-table
columns, whether GIFs apply, where the model's output lives — and those
instructions override the generic template below. If a test has no rubric
file, do not invent one: leave its Score column and averages blank and note
the missing rubric in the report.

Rules that apply to every test:

- Report the score as a breakdown plus total, e.g. `3 + 2 + 3 + 1 = 9/10`, in
  each run's subsection, and put the total in the summary table's Score
  column.
- The **combined average score** for a test is the arithmetic mean of its runs'
  totals, rounded to one decimal place (e.g. `7.3/10`). State it directly
  below the test's summary table as `**Average score: <x.y>/10**`.
- Cross-test and cross-model comparisons live in the index at
  `results/README.md`, not in per-test reports; a report covers exactly one
  model on exactly one test.
- Apply the rubric identically across models and runs; the score column is the
  factual verdict — keep subjective impressions out of it and in the prose.

## Recording GIFs for CLI / TUI results

For prompts whose result is a command-line or terminal UI program, capture a
short demo GIF of each successfully running result and embed it in the report.

- Preferred tool: `vhs` (Charmbracelet). **Verify it is installed first**
  (`command -v vhs`) — do not assume, as these tests may run on other machines.
  If `vhs` is missing, check for alternatives (e.g. `asciinema` + `agg`,
  `t-rec`); if no recording tool is available, skip GIFs and state so in the
  report rather than installing anything.
- With `vhs`: write a `.tape` script per run that launches the program, waits
  long enough to show it off (a few seconds — include any interaction the
  program invites, such as key presses), then exercises the quit key.
- Save the GIF and its `.tape` script inside the run's numbered directory as
  `demo.gif` and `demo.tape`, and embed it in the README with a relative path
  from the test directory, e.g. `![run 1](1/demo.gif)`.
- If a run fails to build or crashes at startup, do not record a GIF; report
  the failure instead.

## Report template

Every `README.md` must follow this structure so reports look and feel
consistent across models and tests:

```markdown
# <model name> — <test-name>

Results for locally hosted model `<model name>` on the `<test-name>` test.

> One-sentence summary of what `prompts/<test-name>/prompt.md` asked for,
> including its strict requirement(s).

| Run | Builds | Runs | Meets requirements | Score | Notes |
|-----|--------|------|--------------------|-------|-------|
| 1   | ✅/❌  | ✅/❌ | ✅/❌             | n/10  | one-liner |
| …   |        |      |                    |       |         |

**Average score: <x.y>/10**

## Run 1

![run 1](1/demo.gif)

Two to five sentences: what the model built, what's notable (creativity,
correctness, style), and any failures with the relevant error message.
End with the score breakdown, e.g. `Score: 3 + 2 + 3 + 1 = 9/10`.

## Run 2
…
```

Rules for consistency, including with future test prompts:

- One `## Run <n>` section per numbered directory, in numeric order.
- Always include the summary table, even for a single run.
- Use the summary-table columns the test's rubric specifies; the columns shown
  above are the build-and-run default.
- The GIF line appears only for CLI/TUI tests where a recording was made; for
  non-interactive prompts (or failed runs), substitute a fenced code block
  with representative output, or omit.
- Keep judgments factual and tied to the prompt's requirements; put subjective
  impressions in the Notes/prose, clearly framed as impressions.
