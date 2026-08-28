# Hello, Bubble Tea! 🧋

An animated "Hello, World!" terminal app built with
[Bubble Tea](https://github.com/charmbracelet/bubbletea) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss).

A big block-letter **HELLO, WORLD!** banner rides a sine wave while a rainbow
scrolls across it, confetti drifts down behind it, and a greeting underneath
cycles through a dozen languages every few seconds.

## Run it

Requires Go 1.21+.

```sh
go run .
```

Or build a binary:

```sh
go build -o hello .
./hello
```

## Controls

| Key            | Action                                      |
| -------------- | ------------------------------------------- |
| `q` / `ctrl+c` | Quit (always)                               |
| `space`        | Toggle party mode (more confetti, faster!)  |
| `enter` / `→`  | Next greeting language                      |
| `←`            | Previous greeting language                  |

## How it works

- The banner is drawn from a tiny hand-rolled 5-row block font onto a character
  canvas; each column's vertical offset follows a sine wave and its color comes
  from an HSL→RGB rainbow that scrolls with the frame counter.
- Confetti particles are spawned at the top each tick, fall with individual
  velocities, and sway side to side as they drop.
- A 30 fps `tea.Tick` drives the animation; all state lives in the Bubble Tea
  model and every frame is a pure render of that state.
