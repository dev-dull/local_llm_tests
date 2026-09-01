# Qwen3.6-35B-A3B — hello-go-bubbletea

Results for locally hosted model `Qwen3.6-35B-A3B` on the `hello-go-bubbletea` test.

> Build a "Hello, World!" app in Go using `charmbracelet/bubbletea`, with creative
> flair (animation, color, interactivity); the one strict requirement is that
> pressing `q` must always quit.

| Run | Builds | Runs | Meets requirements | Score | Notes |
|-----|--------|------|--------------------|-------|-------|
| 1   | ✅     | ✅   | ✅                 | 7/10  | Four animation modes; `q` quits, but layout overflows and it busy-loops the CPU |
| 2   | ✅     | ✅   | ✅                 | 8/10  | Bouncing boxed greeting; works, with a box-border rendering artifact |
| 3   | ✅     | ❌   | ❌                 | 5/10  | Panics at startup: `strings.Repeat` with negative count in View |
| 4   | ✅     | ✅   | ✅                 | 9/10  | Bouncing bubble with particle bursts; works, some dead code |

**Average score: 7.3/10**

## Run 1

![run 1](1/demo.gif)

The most ambitious run of the batch: space cycles four animation modes (wave,
rain, spin, pulse), and the spin mode — the greeting orbiting in a rotating
ring of colored letters — is genuinely striking. It builds and runs, and `q`
quits cleanly. But the View emits more lines than the terminal has, so the
header with the controls is pushed off-screen, and its "ticker" command
returns a message immediately with no delay, making the event loop a
CPU-pegging busy loop. Rain positions are mutated inside `View()`, and the
"random" helper is just the current millisecond modulo n, so values correlate.
Score: 3 + 1 + 3 + 0 = 7/10

## Run 2

![run 2](2/demo.gif)

A rainbow-bordered box containing a rainbow "Hello, World!" bounces around the
screen, steerable with the arrow keys, with a pause toggle, a status light,
and a decorative gradient progress bar. Builds and runs; `q` and ctrl+c quit.
The box is positioned by prefixing spaces to only the first of its three
lines, so a detached piece of border floats at the box's x-offset while the
rest hugs the left margin — a visible artifact — and centering uses `len()` on
ANSI-styled strings, so alignment drifts. Amusingly, it imports the `progress`
bubble yet hand-rolls everything else.
Score: 3 + 1 + 3 + 1 = 8/10

## Run 3

```text
Caught panic:

goroutine 1 [running]:
...
strings.Repeat({...}, 0x0?)
        .../strings/strings.go:628
main.(*model).View(0x14000160000)
        .../Qwen3.6-35B-A3B/hello-go-bubbletea/3/main.go:239
```

Builds cleanly but crashes at startup, reproducibly: `View()` runs before the
first `WindowSizeMsg` arrives, so `strings.Repeat("─", min(m.width-2, 80))`
is called with a negative count and panics. The design on paper — a starfield
with a typewriter reveal, per-character bobbing, wind control, and emoji
particle bursts — is the most creative of this model's four runs, but none of
it can be seen, and the strict `q`-quits requirement can't be satisfied by a
program that dies before accepting input. No GIF was recorded per the
crash-at-startup rule.
Score: 3 + 0 + 2 + 0 = 5/10

## Run 4

![run 4](4/demo.gif)

A bubble (🫧) drifts and bounces around a tinted background while a tri-color
"Hello, World!" sits center-screen; space fires emoji particle bursts with a
screen-shake, WASD/arrows steer, r resets, and a spinner animates in the
footer. Builds and runs without visible errors, and `q` quits (though ctrl+c
is notably not bound — only `q`/esc). Quality is middling: it hand-rolls raw
ANSI escape codes alongside lipgloss, the `greetingPhase`/`entranceProgress`/
`maxLife` fields are dead, and `intMax`/`intMin` reimplement builtins.
Score: 3 + 2 + 3 + 1 = 9/10
