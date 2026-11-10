# Hello World - Bubbletea Edition

A beautiful and interactive "Hello, World!" application built with Go and [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea).

## Features

- **🎨 Animated Particle System**: Colorful particles that dance across the screen
- **🌟 Phase-Based UI**: Three different animated phases that transition automatically  
- **💫 Interactive Elements**: Press 'h' for help, 'q' to quit, spacebar to restart
- **🎨 Beautiful Terminal UI**: Using charmbracelet/lipgloss for styling
- **✨ Smooth Animations**: 60 FPS animations with particle physics

## Installation

1. Ensure you have Go 1.21 or later installed
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running

Simply run:
```bash
go run main.go
```

Or build and run:
```bash
go build -o hello-bubbletea
./hello-bubbletea
```

## Controls

- **q** or **Ctrl+C**: Quit the application
- **h** or **?**: Toggle help screen
- **Space**: Restart the animation

## Demo Mode

For testing without a terminal environment, use the demo mode:
```bash
go run demo.go
```

## Technical Details

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) TUI framework
- Styled with [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- Features custom particle physics with life cycles
- Responsive terminal UI that adapts to window size
- Color cycling and smooth transitions

## Project Structure

```
.
├── main.go      # Main application with interactive animations
├── demo.go      # Demo mode for testing without terminal
├── go.mod       # Go module dependencies
├── README.md    # This file
└── prompt.md    # Original task instructions
```

## Inspiration

This is a creative take on the classic "Hello, World!" program, demonstrating:
- The power of modern Go terminal UI libraries
- Beautiful animations in constrained terminal environments  
- Interactive programming paradigms
- The joy of creative coding in Go

## Notes

- The interactive version requires a terminal environment with TTY access
- The demo mode provides a preview of the animated output
- Build successfully completes and demo runs properly in headless environments