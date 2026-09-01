# Rubric: hello-go-bubbletea

A build-and-run test: the model delivered a Go/Bubble Tea program. Actually
build and run each result as described in `results/prompt.md`.

Score each run out of 10:

- **Builds (0–3):** 3 = builds cleanly as delivered; 2 = builds only after
  routine fixes you are allowed to make (e.g. `go mod tidy`); 0 = does not
  build. (1 is unused; building is nearly all-or-nothing.)
- **Runs (0–2):** 2 = starts and runs without crashing; 1 = runs but with
  visible errors or glitches; 0 = crashes at startup or doesn't run.
- **Meets requirements (0–3):** 3 = satisfies every strict requirement in the
  test's `prompt.md` (here, above all: `q` always quits); deduct 1 per missed
  requirement, down to 0.
- **Quality (0–2):** 2 = clean, idiomatic code with no notable defects;
  1 = works but with sloppy or non-idiomatic code, dead code, or hallucinated/
  unnecessary dependencies; 0 = significant bugs or badly structured code even
  if requirements are technically met.

Rules:

- A run that does not build is capped at its Builds points (score ≤ 3 total);
  do not award Runs, Requirements, or Quality points for code you could not
  execute.

Report specifics:

- Summary-table columns: `| Run | Builds | Runs | Meets requirements | Score | Notes |`.
- This is a TUI test: record demo GIFs per the GIF instructions in
  `results/prompt.md`.
