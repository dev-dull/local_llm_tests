# hello-wave 🌊

A "Hello, World!" that refuses to sit still. Built with
[Bubble Tea](https://github.com/charmbracelet/bubbletea) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss).

The greeting rides a rainbow sine wave across a twinkling starfield.
Fire confetti at it. Say hello in eight languages. It's a whole thing.

## Build & run

Requires Go 1.21+.

```sh
go build -o hello-wave .
./hello-wave
```

Or just:

```sh
go run .
```

## Controls

| Key             | Action                        |
| --------------- | ----------------------------- |
| `space`         | Confetti burst 🎉             |
| `←` / `→` (or `h`/`l`) | Cycle greeting language |
| `p`             | Pause / resume the animation  |
| `q` / `esc` / `ctrl+c` | Quit                   |

`q` always quits — even while paused.

## How it works

- A 30 fps `tea.Tick` drives a frame counter; each letter's vertical
  position is `sin(frame·speed + index·freq)`, and its color walks the
  HSL hue wheel offset by index, giving the rolling rainbow.
- Confetti particles get a random radial velocity and a little gravity;
  they die when their lifetime runs out or they fall off screen.
- Everything is composited into a rune grid each frame, then rendered
  with runs of identically-styled cells batched into single Lip Gloss
  calls to keep per-frame allocations sane.
