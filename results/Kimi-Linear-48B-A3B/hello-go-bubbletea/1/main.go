package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type Model struct {
	cursor      int
	choices     []string
	showHelp    bool
	selected    map[int]bool
	animation   float64
	particles   []particle
	colorIndex  int
	showColors  bool
	inputMode   string
	inputText   string
	showInput   bool
	score       int
	highScore   int
	gameState   string // "menu", "game", "highscore"
}

type particle struct {
	x, y        float64
	vx, vy      float64
	color       string
	size        float64
	life        float64
	maxLife     float64
}

// Generate particles for background animation
func generateParticles(count int) []particle {
	particles := make([]particle, count)
	colors := []string{"red", "blue", "green", "yellow", "magenta", "cyan"}

	for i := range particles {
		particles[i] = particle{
			x:        float64(rand.Intn(80) + 10),
			y:        float64(rand.Intn(20) + 10),
			vx:       rand.Float64()*2 - 1,
			vy:       rand.Float64()*2 - 1,
			color:    colors[rand.Intn(len(colors))],
			size:     rand.Float64()*0.5 + 0.5,
			life:     rand.Float64() * 100,
			maxLife:  rand.Float64() * 100 + 100,
		}
	}
	return particles
}

// Initializer
func initialModel() Model {
	rand.Seed(time.Now().UnixNano())

	// Set high score from environment or defaults
	highScore := 0
	if hs := os.Getenv("BUBBLE_HIGH_SCORE"); hs != "" {
		fmt.Sscanf(hs, "%d", &highScore)
	}

	return Model{
		cursor:     0,
		choices:    []string{"Play Game", "Show Colors", "Input Mode", "High Score", "Quit"},
		selected:   make(map[int]bool),
		animation:  0,
		particles:  generateParticles(50),
		colorIndex: 0,
		showColors: false,
		inputMode:  "name",
		inputText:  "",
		showInput:  false,
		score:      0,
		highScore:  highScore,
		gameState:  "menu",
	}
}

// Model implementation for Bubbletea
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles state transitions and logic
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case " ", "enter", "\r":
			switch m.cursor {
			case 0: // Play Game
				m.gameState = "game"
				m.score = 0
				return m, m.startGame()

			case 1: // Show Colors
				m.showColors = !m.showColors
				if m.showColors {
					m.choices[1] = "Hide Colors"
				} else {
					m.choices[1] = "Show Colors"
				}

			case 2: // Input Mode
				m.showInput = true
				m.inputText = ""
				m.gameState = "input"

			case 3: // High Score
				m.gameState = "highscore"

			case 4: // Quit
				return m, tea.Quit
			}

		case "esc":
			if m.showInput {
				m.showInput = false
				m.gameState = "menu"
			}
		}

	case tea.WindowSizeMsg:
		// Handle window resize if needed

	case tickMsg:
		if m.gameState == "game" {
			// Update game logic
			m.updateGame()

			// Update particles
			for i := range m.particles {
				m.updateParticle(&m.particles[i])
			}
		}

	case inputMsg:
		if m.showInput {
			m.inputText = string(msg.Runes)

			if msg.String() == "\r" {
				m.showInput = false
				m.gameState = "menu"
				// Save high score
				if m.score > m.highScore {
					m.highScore = m.score
					os.Setenv("BUBBLE_HIGH_SCORE", fmt.Sprintf("%d", m.highScore))
				}
			}
		}
	}

	// Update animation
	if m.gameState == "game" {
		m.animation += 0.1

		// Cycle through colors
		m.colorIndex = int(m.animation) % 6

		// Update score
		m.score = int(m.animation) * 10
	}

	return m, cmd
}

// View renders the UI
func (m Model) View() string {
	if m.showColors {
		return m.renderColorDemo()
	}

	if m.showInput {
		return m.renderInputMode()
	}

	if m.gameState == "game" {
		return m.renderGame()
	}

	if m.gameState == "highscore" {
		return m.renderHighScore()
	}

	return m.renderMenu()
}

// Render color demonstration
func (m Model) renderColorDemo() string {
	styles := []lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("red")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("blue")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("green")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("yellow")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("magenta")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("cyan")),
	}

	var b strings.Builder
	b.WriteString("Color Demo - Press ESC to return\n\n")

	for i, style := range styles {
		if i == m.colorIndex {
			b.WriteString(style.Render("◆ " + style.Render("Color " + string(rune('A'+i))) + " ◆"))
		} else {
			b.WriteString(style.Render("  " + style.Render("Color " + string(rune('A'+i))) + "  "))
		}
		if i < len(styles)-1 {
			b.WriteString(" | ")
		}
	}

	return b.String()
}

// Render input mode
func (m Model) renderInputMode() string {
	var b strings.Builder
	b.WriteString("Input Mode - Type your name and press ENTER\n\n")
	b.WriteString("Name: ")
	b.WriteString(m.inputText)
	if m.showInput {
		b.WriteString("_")
	}
	b.WriteString("\n\nPress ESC to cancel")
	return b.String()
}

// Render game screen
func (m Model) renderGame() string {
	var b strings.Builder

	// Game UI
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("🎮 BUBBLE GAME 🎮"))
	b.WriteString("\n\n")

	// Score display
	scoreStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("green")).
		Bold(true).
		Padding(1, 2)

	b.WriteString(scoreStyle.Render(fmt.Sprintf("Score: %d", m.score)))
	b.WriteString("\n")
	b.WriteString(scoreStyle.Render(fmt.Sprintf("High Score: %d", m.highScore)))
	b.WriteString("\n\n")

	b.WriteString("Collect points by exploring! Press 'q' to quit anytime.\n")
	b.WriteString("Use arrow keys to navigate menus in the main screen.\n")

	return b.String()
}

// Render high score screen
func (m Model) renderHighScore() string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().Bold(true).Render("🏆 HIGH SCORE 🏆"))
	b.WriteString("\n\n")

	// Display top scores (simplified)
	b.WriteString(fmt.Sprintf("Current Score: %d\n", m.score))
	b.WriteString(fmt.Sprintf("High Score: %d", m.highScore))

	if m.score > 0 {
		b.WriteString("\n\n🎉 New High Score! 🎉")
	}

	b.WriteString("\n\nPress any key to continue...")

	return b.String()
}

// Render main menu
func (m Model) renderMenu() string {
	// Create title with animation
	title := m.animateTitle()

	var b strings.Builder
	b.WriteString(title)
	b.WriteString("\n\n")

	// Create menu with selection highlighting
	for i, choice := range m.choices {
		if i == m.cursor {
			// Selected item - highlight with color
			color := m.getColorForIndex(i)
			selectedStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(color)).
				Bold(true).
				Background(lipgloss.Color("blue")).
				Padding(1, 2)

			b.WriteString(selectedStyle.Render("▶ " + choice + " ◀"))
		} else {
			// Unselected item - dimmed
			dimStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("gray")).
				Padding(1, 2)

			b.WriteString(dimStyle.Render("  " + choice + "  "))
		}

		if i < len(m.choices)-1 {
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\nPress 'q' or 'ctrl+c' to quit")

	return b.String()
}

// Animate title with color cycling
func (m Model) animateTitle() string {
	colors := []string{"red", "blue", "green", "yellow", "magenta", "cyan"}
	color := colors[m.colorIndex]

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Underline(true).
		MarginLeft(10).
		MarginRight(10)

	return titleStyle.Render("🎮 Hello Go Bubbletea World! 🎮")
}

// Get color for menu index
func (m Model) getColorForIndex(index int) string {
	colors := []string{"red", "blue", "green", "yellow", "magenta", "cyan"}
	return colors[index%len(colors)]
}

// Update particle physics
func (m *Model) updateParticle(p *particle) {
	p.x += p.vx
	p.y += p.vy
	p.life++

	// Bounce off edges
	if p.x <= 0 || p.x >= 80 {
		p.vx = -p.vx
	}
	if p.y <= 0 || p.y >= 24 {
		p.vy = -p.vy
	}

	// Reset particle if life exceeded
	if p.life > p.maxLife {
		p.x = float64(rand.Intn(80) + 10)
		p.y = float64(rand.Intn(20) + 10)
		p.life = 0
		p.maxLife = rand.Float64() * 100 + 100
	}
}

// Update game logic
func (m *Model) updateGame() {
	// Simple game progression
	// In a real game, this would handle actual game mechanics
}

// Start game command
func (m Model) startGame() tea.Cmd {
	return func() tea.Msg {
		return inputMsg{}
	}
}

// Helper function for tick messages
type tickMsg time.Time

type inputMsg struct {
	Runes []rune
}

// String method for inputMsg
func (i inputMsg) String() string {
	return fmt.Sprintf("input: %s", string(i.Runes))
}

// Main function
func main() {
	// Set up high score environment variable
	highScore := os.Getenv("BUBBLE_HIGH_SCORE")
	if highScore == "" {
		highScore = "0"
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}