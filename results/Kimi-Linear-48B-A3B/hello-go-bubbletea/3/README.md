# Hello, World! - Bubbletea Edition

A creative and interactive "Hello, World!" application built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea), featuring particle animations, color effects, and engaging interactivity.

## Features

- **Interactive gameplay**: Control a cursor with arrow keys to score points
- **Particle system**: Beautiful animated particles that respond to your movements
- **Multiple game states**: Menu, Playing, Paused, Game Over, and Win states
- **Dynamic UI**: Real-time score display and helpful instructions
- **Visual effects**: Particle bursts, color gradients, and smooth animations
- **Always-quitting**: Press `q` or `Ctrl+C` to quit at any time

## Requirements

- Go 1.21 or later
- Terminal with ANSI color support

## Installation

1. Clone or download this project
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Application

```bash
# Run the application
go run main.go

# Or build and run
go build -o hello-go-bubbletea
./hello-go-bubbletea
```

## Controls

- **Arrow Keys**: Move the cursor and score points
- **SPACE**: Pause/unpause the game
- **Q**: Quit the application (always works!)
- **Ctrl+C**: Force quit

## How to Play

1. Start the game from the menu by pressing SPACE
2. Use arrow keys to move the cursor around the screen
3. Each move scores points and creates a particle burst effect
4. Press SPACE to pause/resume the game
5. Press Q to quit at any time

## Game States

- **Menu**: Welcome screen with instructions
- **Playing**: Active gameplay with particle effects
- **Paused**: Game suspended, UI overlay shown
- **Game Over**: End state when you quit
- **Win**: Victory celebration (achieved by reaching certain score)

## Technical Details

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) for terminal UI
- Uses [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling and colors
- Particle system with physics (gravity, friction, life cycles)
- Responsive design that adapts to terminal size
- Smooth animations and visual feedback

## Creative Elements

- Animated background particles
- Colorful particle bursts on movement
- Dynamic lighting effects
- Smooth transitions between game states
- Professional UI styling with gradients and shadows

Enjoy exploring the interactive world of Bubbletea! 🎮✨