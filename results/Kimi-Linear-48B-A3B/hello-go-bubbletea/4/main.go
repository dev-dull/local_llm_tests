package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type Model struct {
	quitting       bool
	message        string
	colorIndex     int
	animatingText  []string
	currentFrame   int
	showAnimation  bool
	particles      []particle
	sparkles       []sparkle
	cursorPosition int
}

type particle struct {
	x, y        int
	speed       float64
	direction   float64
	color       lipgloss.Color
	size        int
	life        int
	maxLife     int
}

type sparkle struct {
	x, y        int
	size        int
	direction   float64
	speed       float64
	rotation    float64
	life        int
	maxLife     int
	color       lipgloss.Color
}

// Model represents the application state
type Model struct {
	quitting       bool
	message        string
	colorIndex     int
	animatingText  []string
	currentFrame   int
	showAnimation  bool
	particles      []particle
	sparkles       []sparkle
	cursorPosition int
}

type particle struct {
	x, y        int
	speed       float64
	direction   float64
	color       lipgloss.Color
	size        int
	life        int
	maxLife     int
}

type sparkle struct {
	x, y        int
	size        int
	direction   float64
	speed       float64
	rotation    float64
	life        int
	maxLife     int
	color       lipgloss.Color
}

// Initializer for the model
func initialModel() Model {
	m := Model{
		message:        "Hello, World!",
		colorIndex:     0,
		animatingText:  []string{},
		currentFrame:   0,
		showAnimation:  true,
		particles:      make([]particle, 0),
		sparkles:       make([]sparkle, 0),
		cursorPosition: 0,
	}

	// Initialize animating text with colors
	colors := []lipgloss.Color{
		lipgloss.Color("205"),
		lipgloss.Color("213"),
		lipgloss.Color("39"),
		lipgloss.Color("213"),
		lipgloss.Color("39"),
		lipgloss.Color("205"),
	}

	for i, char := range m.message {
		m.animatingText = append(m.animatingText, string(char))
		m.particles = append(m.particles, particle{
			x:         i * 4,
			y:         10,
			speed:     float64(0.5 + float64(i%3)*0.2),
			direction: float64(i%360),
			color:     colors[i%len(colors)],
			size:      1,
			life:      0,
			maxLife:   30,
		})
	}

	return m
}

// Initialize particles for effect
func (m *Model) initParticles() {
	// Create background particles
	for i := 0; i < 20; i++ {
		m.particles = append(m.particles, particle{
			x:         5 + (i%40),
			y:         15 + (i%20),
			speed:     0.1 + float64(i%10)*0.05,
			direction: float64(i%360),
			color:     lipgloss.Color("240"),
			size:      1,
			life:      i * 2,
			maxLife:   100,
		})
	}

	// Create sparkle effects
	for i := 0; i < 10; i++ {
		m.sparkles = append(m.sparkles, sparkle{
			x:         10 + (i%30),
			y:         8 + (i%15),
			size:      2,
			direction: float64(i%360),
			speed:     0.2 + float64(i%5)*0.1,
			rotation:  0,
			life:      i * 5,
			maxLife:   50,
			color:     lipgloss.Color("226"),
		})
	}
}

// Update handles state transitions
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case " ", "enter":
			m.showAnimation = !m.showAnimation
			return m, nil

		case "r":
			return initialModel(), nil
		}

	case tea.WindowSizeMsg:
		return m, nil

	case time.TickMsg:
		// Update animation
		m.currentFrame++

		// Update particles
		for i := range m.particles {
			p := &m.particles[i]
			p.x += int(float64(p.speed) * float64(time.Now().UnixNano()%1000/1000))
			p.y += int(p.speed)
			p.life++

			// Reset particle if it goes off screen
			if p.x > 80 || p.y > 24 || p.life > p.maxLife {
				p.x = 5 + (i % 40)
				p.y = 15 + (i % 20)
				p.life = 0
			}
		}

		// Update sparkles
		for i := range m.sparkles {
			s := &m.sparkles[i]
			s.rotation += s.speed
			s.life++

			// Reset sparkle if it goes off screen
			if s.x > 80 || s.y > 24 || s.life > s.maxLife {
				s.x = 10 + (i % 30)
				s.y = 8 + (i % 15)
				s.life = 0
			}
		}

		// Update color animation
		m.colorIndex = (m.colorIndex + 1) % 6

		return m, nil
	}

	return m, nil
}

// View renders the application
func (m Model) View() string {
	// Create styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(1, 2)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(2)

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		MarginTop(4)

	// Create the view
	var view strings.Builder

	// Title
	view.WriteString(titleStyle.Render("🚀 Hello, World! Bubbletea Edition") + "\n\n")

	// Animation indicator
	if m.showAnimation {
		view.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Render("✨ Animation Active") + "\n\n")
	}

	// Main message with color animation
	coloredMessage := ""
	for i, char := range m.message {
		colorIndex := i % len(m.message)
		color := lipgloss.Color(fmt.Sprintf("%d", 205+colorIndex*8))
		coloredChar := lipgloss.NewStyle().Foreground(color).Render(string(char))
		coloredMessage += coloredChar
	}
	view.WriteString(coloredMessage + "\n\n")

	// Instructions
	if m.showAnimation {
		view.WriteString(instructionStyle.Render("Press 'q' to quit • Press ' ' to toggle animation • Press 'r' to reset") + "\n\n")
	} else {
		view.WriteString(instructionStyle.Render("Press 'q' to quit • Press ' ' to start animation • Press 'r' to reset") + "\n\n")
	}

	// Particle effects (if enabled)
	if m.showAnimation {
		for _, p := range m.particles {
			if p.life < p.maxLife {
				particleStyle := lipgloss.NewStyle().
					Foreground(p.color).
					Background(lipgloss.Color("black")).
					Render("•")
				// Simple positioning for particles
				if p.x < 80 && p.y < 24 {
					view.WriteString(particleStyle)
				}
			}
		}

		for _, s := range m.sparkles {
			if s.life < s.maxLife {
				sparkleStyle := lipgloss.NewStyle().
					Foreground(s.color).
					Render("✦")
				if s.x < 80 && s.y < 24 {
					view.WriteString(sparkleStyle)
				}
			}
		}
	}

	// Status line
	status := ""
	if m.quitting {
		status = "\n\n👋 Goodbye!"
	} else {
		status = "\n\n💡 Press 'q' to quit anytime!"
	}

	view.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(status))

	return view.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's an error: %v", err)
		return
	}
}