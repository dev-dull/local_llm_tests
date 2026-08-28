# Hello, World! ✨

A small terminal celebration built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

The greeting bounces on a sine wave, colored by a moving rainbow, over a field
of twinkling stars. Cycle through "Hello, World!" in 18 languages — each switch
fires a burst of confetti.

## Run

Requires Go 1.21+ and a terminal with color support (truecolor looks best).

```sh
go run .
```

Or build a binary:

```sh
go build -o hello .
./hello
```

## Keys

| Key | Action |
| --- | --- |
| `space`, `→`, `l` | Next language (with confetti) |
| `←`, `h` | Previous language (with confetti) |
| `enter` | Confetti, just because |
| `q` (also `esc`, `ctrl+c`) | Quit |
