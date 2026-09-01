# Task: Generate result reports for local LLM test runs

You are generating evaluation reports for prompts that were run against locally
hosted LLM models. Work through the directory structure described below and
produce one `README.md` report per model directory.

## Directory layout

- `prompts/<test-name>/prompt.md` — the original prompt that was given to each
  model.
- `results/<model>/` — one directory per model (the directory name is the model
  name).
- `results/<model>/<test-name>/<n>/` — numbered directories (`1/`, `2/`, `3/`,
  …), each holding one independent run of that prompt by that model.
- `results/<model>/README.md` — the report you generate. One per model,
  covering all of that model's tests and runs.

## Procedure

For each model directory under `results/`:

1. **If `results/<model>/README.md` already exists, skip this model
   entirely.** Do not read the existing file — assume it is complete and move
   on to the next model directory.
2. For each `<test-name>` subdirectory, read the matching
   `prompts/<test-name>/prompt.md` so you understand what the model was asked
   to do and what the strict requirements are.
3. For each numbered run directory, evaluate the result:
   - What did the model produce (files, completeness)?
   - Does it build? Actually build it (e.g. `go build` for Go projects) and
     record success or the exact failure.
   - Does it run, and does it satisfy the prompt's strict requirements (e.g.
     for `hello-go-bubbletea`: does `q` quit)?
   - Note anything interesting: creative choices, bugs, hallucinated
     dependencies, deviations from the prompt.
4. Score the run using the scoring rules below.
5. Write the report as `results/<model>/README.md` using the template below.

Do not modify anything inside the numbered run directories except as needed to
build (e.g. running `go mod tidy`); never edit the model's source code. If a
run doesn't build or run, report that honestly — do not fix it.

## Scoring

Each run gets a **score out of 10**, and each test gets a **combined average
score** across its runs, so models can be compared at a glance.

Not every test is a build-and-run test. If a test's directory in `prompts/`
contains an `answer-key.md`, score that test's runs with the key's rubric and
follow its instructions (including any changes to the summary-table columns)
in place of the generic rubric below. Never include an answer key in any
model prompt.

Score each run against this rubric:

- **Builds (0–3):** 3 = builds cleanly as delivered; 2 = builds only after
  routine fixes you are allowed to make (e.g. `go mod tidy`); 0 = does not
  build. (1 is unused; building is nearly all-or-nothing.)
- **Runs (0–2):** 2 = starts and runs without crashing; 1 = runs but with
  visible errors or glitches; 0 = crashes at startup or doesn't run.
- **Meets requirements (0–3):** 3 = satisfies every strict requirement in the
  test's `prompt.md`; deduct 1 per missed requirement, down to 0.
- **Quality (0–2):** 2 = clean, idiomatic code with no notable defects;
  1 = works but with sloppy or non-idiomatic code, dead code, or hallucinated/
  unnecessary dependencies; 0 = significant bugs or badly structured code even
  if requirements are technically met.

Rules:

- A run that does not build is capped at its Builds points (score ≤ 3 total);
  do not award Runs, Requirements, or Quality points for code you could not
  execute.
- Report the score as a breakdown plus total, e.g. `3 + 2 + 3 + 1 = 9/10`, in
  each run's subsection, and put the total in the summary table's Score
  column.
- The **combined average score** for a test is the arithmetic mean of its runs'
  totals, rounded to one decimal place (e.g. `7.3/10`). State it directly
  below the test's summary table as `**Average score: <x.y>/10**`.
- If a model has multiple tests, end the report with an `## Overall` section
  giving the model's overall average: the mean of the per-test average scores,
  rounded to one decimal place.
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
  `demo.gif` and `demo.tape`, and embed it in the README with a relative path,
  e.g. `![run 1](hello-go-bubbletea/1/demo.gif)`.
- If a run fails to build or crashes at startup, do not record a GIF; report
  the failure instead.

## Report template

Every `README.md` must follow this structure so reports look and feel
consistent across models:

```markdown
# <model directory name>

Results for locally hosted model `<model name>`.

## <test-name>

> One-sentence summary of what `prompts/<test-name>/prompt.md` asked for,
> including its strict requirement(s).

| Run | Builds | Runs | Meets requirements | Score | Notes |
|-----|--------|------|--------------------|-------|-------|
| 1   | ✅/❌  | ✅/❌ | ✅/❌             | n/10  | one-liner |
| …   |        |      |                    |       |         |

**Average score: <x.y>/10**

### Run 1

![run 1](<test-name>/1/demo.gif)

Two to five sentences: what the model built, what's notable (creativity,
correctness, style), and any failures with the relevant error message.
End with the score breakdown, e.g. `Score: 3 + 2 + 3 + 1 = 9/10`.

### Run 2
…

## Overall

**Overall average score: <x.y>/10** (only when the model has more than one
test; mean of the per-test averages.)
```

Rules for consistency, including with future test prompts:

- One `## <test-name>` section per test, in alphabetical order; one
  `### Run <n>` subsection per numbered directory, in numeric order.
- Always include the summary table, even for a single run.
- The GIF line appears only for CLI/TUI tests where a recording was made; for
  non-interactive prompts (or failed runs), substitute a fenced code block
  with representative output, or omit.
- Keep judgments factual and tied to the prompt's requirements; put subjective
  impressions in the Notes/prose, clearly framed as impressions.
