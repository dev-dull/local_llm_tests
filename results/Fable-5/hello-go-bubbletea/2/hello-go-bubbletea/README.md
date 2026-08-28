# hello-go-bubbletea

"Hello, World!" — but it rides a rainbow sine wave through a field of
twinkling sparkles, courtesy of [Bubble Tea](https://github.com/charmbracelet/bubbletea)
and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

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

## Keys

| Key     | Action                                              |
| ------- | --------------------------------------------------- |
| `q`     | Quit (always — also `esc` and `ctrl+c`)             |
| `enter` | Cycle greetings (9 languages)                       |
| `tab`   | Cycle color themes (rainbow, ocean, sunset, neon)   |
| `space` | Pause / resume time                                 |

## How it works

Each frame (30 fps via `tea.Tick`) the letters of the greeting are placed
on a character grid at a vertical offset of `sin(phase + i·0.35)`, giving a
traveling wave. Each letter's color is computed in HSL — the hue sweeps
across the text and drifts over time — then converted to a truecolor hex
for Lip Gloss. Sparkles spawn at random grid positions and fade through
three glyph/brightness stages (`✦ → ✧ → ·`) over their lifetime.

Best viewed in a truecolor terminal (most modern terminals qualify);
elsewhere colors gracefully degrade to the nearest ANSI palette entry.
