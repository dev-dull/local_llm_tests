# Hello, World! - Bubbletea Edition

A delightful "Hello, World!" application built with Go and [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea), featuring colorful animations, particle effects, and interactive elements.

## Features

- 🎨 **Colorful Text Animation**: The "Hello, World!" message features animated colors that cycle through a vibrant palette
- ✨ **Particle Effects**: Background particles float across the screen for a dynamic atmosphere
- ✦ **Sparkle Effects**: Sparkles randomly appear and rotate for extra visual flair
- 🎮 **Interactive Controls**:
  - Press `q` to quit the application
  - Press ` ` (spacebar) to toggle animation effects
  - Press `r` to reset the application
- 🚀 **Responsive Design**: Adapts to terminal size changes

## Requirements

- Go 1.21 or later
- Terminal with Unicode support

## Installation

1. Clone or download this project
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Application

```bash
go run main.go
```

Or build it first:
```bash
go build -o hello-bubbletea main.go
./hello-bubbletea
```

## Controls

- **q** - Quit the application (always works!)
- **Space** - Toggle animation effects
- **r** - Reset the application

## Technical Details

This application demonstrates:
- Bubbletea's reactive programming model
- Custom model state management
- Terminal-based animation techniques
- Color and styling with lipgloss
- Particle system implementation

The particle effects and animations are purely visual enhancements that create an engaging user experience while maintaining the simplicity of a "Hello, World!" demonstration.

## License

This project is part of the "Hello, World!" collection and follows the same open source principles.