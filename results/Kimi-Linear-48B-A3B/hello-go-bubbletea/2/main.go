package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type Model struct {
	phase          int
	colorIndex     int
	particles      []particle
	showHelp       bool
	animationFrame int
	quitConfirmed  bool
}

type particle struct {
	x, y       int
	speed      float64
	color      lipgloss.Color
	life       int
	maxLife    int
	direction  float64
}

// Define colors
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(1, 2).
		Margin(2, 0)

	particleStyle = func(c lipgloss.Color) lipgloss.Style {
		return lipgloss.NewStyle().
			Foreground(c).
			Bold(true).
			Background(lipgloss.Color("0"))
	}

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		Margin(1, 0)

	progressStyle = func(progress float64) lipgloss.Style {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true).
			Width(int(float64(40) * progress))
		}
)

// Initializer
func initialModel() Model {
	return Model{
		phase:          0,
		colorIndex:     0,
		particles:      createParticles(),
		showHelp:       true,
		animationFrame: 0,
		quitConfirmed:  false,
	}
}

// Create particles for animation
func createParticles() []particle {
	particles := make([]particle, 0, 50)
	colors := []lipgloss.Color{
		lipgloss.Color("205"), // Red
		lipgloss.Color("220"), // Orange
		lipgloss.Color("226"), // Yellow
		lipgloss.Color("57"),  // Green
		lipgloss.Color("81"),  // Blue
		lipgloss.Color("123"), // Purple
	}

	for i := 0; i < 50; i++ {
		particles = append(particles, particle{
			x:         40 + (i%10)*4,
			y:         10 + (i/10)*3,
			speed:     float64(1+i%3)*0.1 + 0.05,
			color:     colors[i%len(colors)],
			life:      i % 10,
			maxLife:   10,
			direction: float64(i) * 0.5,
		})
	}
	return particles
}

// Update handles state transitions and animations
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitConfirmed = true
			return m, tea.Quit
		case "h", "?":
			m.showHelp = !m.showHelp
			return m, nil
		case " ":
			// Space to restart animation
			m.particles = createParticles()
			m.phase = 0
			return m, nil
		}

	case time.Time:
		// Handle timer ticks
		m.animationFrame++

		// Update particles
		for i := range m.particles {
			p := &m.particles[i]

			// Update position
			p.x += int(float64(m.animationFrame)*p.speed + math.Sin(p.direction))
			p.y += int(math.Cos(p.direction) * 0.5)

			// Update life
			p.life++
			if p.life > p.maxLife {
				p.life = 0
				p.x = 40 + (i%10)*4
				p.y = 10 + (i/10)*3
			}
		}

		// Color cycling
		m.colorIndex = (m.colorIndex + 1) % 6

		// Phase transitions
		if m.animationFrame%60 == 0 {
			m.phase = (m.phase + 1) % 3
		}

		// Keep the timer running
		return m, tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
			return time.Time{}
		})
	}

	return m, nil
}

// Init is a required method for bubbletea models
func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// View renders the UI
func (m Model) View() string {
	var s string

	// Main content based on phase
	switch m.phase {
	case 0:
		s = renderPhase0(m)
	case 1:
		s = renderPhase1(m)
	case 2:
		s = renderPhase2(m)
	}

	// Always show help
	if m.showHelp {
		s += helpStyle.Render("Press 'h' for help, 'q' to quit, space to restart • Progress: 100%")
	}

	return s
}

// Render different phases
func renderPhase0(m Model) string {
	// Welcome animation
	s := ""
	colors := []string{"205", "220", "226", "57", "81", "123"}

	for i, c := range colors {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(c)).
			Bold(true).
			Italic(true)

		if i == len(colors)-1 {
			s += style.Render("Welcome to Bubbletea!") + "\n\n"
		} else {
			s += style.Render("✨ ") + "\n"
		}
	}

	s += progressStyle(1.0).Render("Complete") + "\n"
	return s
}

func renderPhase1(m Model) string {
	// Interactive animation
	s := ""

	// Title with animation
	title := fmt.Sprintf("🌟 Particle Symphony #%d", m.animationFrame%10)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Margin(1, 0)
	s += titleStyle.Render(title) + "\n\n"

	// Particles
	for _, p := range m.particles {
		if p.life < p.maxLife {
			particleStyle := particleStyle(p.color)
			s += particleStyle.Render("●")
		} else {
			s += " "
		}
	}
	s += "\n\n"

	s += progressStyle(0.7 + math.Sin(float64(m.animationFrame)*0.1)*0.3).Render("Active") + "\n"
	return s
}

func renderPhase2(m Model) string {
	// Celebration phase
	s := ""

	// Celebration title
	celebrateStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46")).
		Background(lipgloss.Color("0"))
	s += celebrateStyle.Render("🎉 SUCCESS! 🎉") + "\n\n"

	// Colorful border
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("226"))

	content := "You've successfully created a beautiful Go application with animations!\n\n"
	content += "Features:\n"
	content += "• Colorful particle system\n"
	content += "• Smooth animations\n"
	content += "• Interactive controls\n"
	content += "• Phase-based UI\n"
	content += "• Help system\n"
	content += "\n"
	content += "Press 'h' for help, 'q' to quit, space to restart!"

	s += borderStyle.Render(content) + "\n\n"

	// Progress indicator
	s += progressStyle(1.0).Render("Final Phase") + "\n"
	return s
}

func main() {
	// Check if bubbletea is installed
	cmd := exec.Command("go", "list", "-m", "github.com/charmbracelet/bubbletea")
	if err := cmd.Run(); err != nil {
		fmt.Println("Installing bubbletea...")
		installCmd := exec.Command("go", "get", "github.com/charmbracelet/bubbletea@latest")
		if err := installCmd.Run(); err != nil {
			fmt.Printf("Failed to install bubbletea: %v\n", err)
			os.Exit(1)
		}
	}

	// Create model
	m := initialModel()

	// Create program
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}