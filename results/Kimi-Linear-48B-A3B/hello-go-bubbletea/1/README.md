# Hello Go Bubbletea World! 🎮

A creative "Hello, World!" application built with Go and [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) that features animations, colors, interactivity, and a fun menu system.

## Features

- **Interactive Menu System**: Navigate with arrow keys (up/down or j/k)
- **Animated Title**: Color-cycling title with underline effect
- **Color Demo**: Press "Show Colors" to see a colorful demonstration
- **Input Mode**: Type your name and see it displayed
- **Game Screen**: Animated background with particles and score tracking
- **High Score System**: Persistent high score using environment variables
- **Particle Animation**: Background particles that bounce around the screen
- **Smooth Transitions**: Clean state management between different screens

## Requirements

- Go 1.21 or later
- The following Go packages:
  - `github.com/charmbracelet/bubbletea` v0.24.2
  - `github.com/charmbracelet/lipgloss` v0.9.1

## Installation

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Build the application**:
   ```bash
   go build -o hello-go-bubbletea main.go
   ```

3. **Run the application**:
   ```bash
   ./hello-go-bubbletea
   ```

## Usage

### Main Menu Options:
1. **Play Game** - Start the interactive game mode
2. **Show Colors** - Display a colorful demonstration
3. **Input Mode** - Type your name and press ENTER
4. **High Score** - View current high score
5. **Quit** - Exit the application (press 'q' or Ctrl+C)

### Game Controls:
- **Arrow Keys** or **j/k**: Navigate menus
- **Space** or **Enter**: Select menu option
- **q** or **Ctrl+C**: Quit anytime

### Game Features:
- Animated particle background
- Real-time score tracking
- High score persistence (saved in environment variable)
- Colorful visual effects

## Environment Variables

- `BUBBLE_HIGH_SCORE`: Sets the initial high score value (optional)

## Creative Elements

- **Animated Title**: The main title "🎮 Hello Go Bubbletea World! 🎮" cycles through colors
- **Particle System**: Background particles with physics simulation
- **Color Cycling**: Dynamic color changes throughout the interface
- **Border Effects**: Rounded borders and highlights for visual appeal
- **Smooth State Management**: Clean transitions between different application states

## Architecture

The application uses a simple model-based architecture:

- `Model`: Holds all application state
- `Update()`: Handles state transitions and logic
- `View()`: Renders the UI based on current state
- Custom message types for different interactions

## Example Output

```
🎮 Hello Go Bubbletea World! 🎮

  Play Game    
  Show Colors  
  Input Mode   
  High Score   
  Quit         

Press 'q' or 'ctrl+c' to quit
```

## Fun Facts

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea), a powerful TUI framework for Go
- Uses [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling
- Features a complete game loop with score tracking
- Demonstrates Go's concurrency and state management capabilities
- Includes particle physics for engaging visual effects

Enjoy exploring the colorful world of Go TUI development! 🌈