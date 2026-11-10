package main

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GameState represents the current state of the game
type GameState int

const (
	MenuState GameState = iota
	PlayingState
	GameOverState
	WinState
)

// Model represents the application state
type Model struct {
	state       GameState
	score       int
	particles   []particle
	text        string
	cursorX     int
	cursorY     int
	width       int
	height      int
	message     string
	shake       float64
	background  []colorParticle
}

// Particle represents a floating particle
type particle struct {
	x, y        float64
	vx, vy      float64
	color       lipgloss.Color
	size        float64
	life        float64
	maxLife     float64
}

// ColorParticle represents a colored background particle
type colorParticle struct {
	x, y        float64
	color       lipgloss.Color
	size        float64
}

// InitialModel creates the initial model
func InitialModel() Model {
	rand.Seed(time.Now().UnixNano())

	return Model{
		state:       MenuState,
		score:       0,
		particles:   make([]particle, 0),
		text:        "Hello, World!",
		cursorX:     0,
		cursorY:     0,
		message:     "Press SPACE to start, Q to quit",
		background:  createBackgroundParticles(),
	}
}

// createBackgroundParticles creates animated background particles
func createBackgroundParticles() []colorParticle {
	particles := make([]colorParticle, 0)
	colors := []lipgloss.Color{
		lipgloss.Color("205"), // orange
		lipgloss.Color("214"), // yellow
		lipgloss.Color("220"), // light yellow
		lipgloss.Color("46"),  // blue
		lipgloss.Color("201"), // pink
	}

	for i := 0; i < 50; i++ {
		particles = append(particles, colorParticle{
			x:     float64(rand.Intn(100)),
			y:     float64(rand.Intn(100)),
			color: colors[rand.Intn(len(colors))],
			size:  float64(rand.Float64() * 2 + 0.5),
		})
	}

	return particles
}

// Initialize sets up the initial application state
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles the application logic
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			// Allow Ctrl+C/ESC to quit
			return m, tea.Quit

		case tea.KeyQ:
			// Q always quits
			return m, tea.Quit

		case tea.KeySpace:
			if m.state == MenuState {
				m.state = PlayingState
				m.message = "Game started! Move with arrow keys, SPACE to pause"
				m.score = 0
				return m, m.updateParticles
			} else if m.state == PlayingState {
				m.state = PausedState
				m.message = "Paused - SPACE to resume, Q to quit"
			} else if m.state == PausedState {
				m.state = PlayingState
				m.message = "Resumed!"
			}

		case tea.KeyArrowUp:
			if m.state == PlayingState {
				m.cursorY = max(0, m.cursorY-1)
				m.score += 10
				m.createParticleBurst(m.cursorX, m.cursorY)
			}

		case tea.KeyArrowDown:
			if m.state == PlayingState {
				m.cursorY = min(m.height-1, m.cursorY+1)
				m.score += 10
				m.createParticleBurst(m.cursorX, m.cursorY)
			}

		case tea.KeyArrowLeft:
			if m.state == PlayingState {
				m.cursorX = max(0, m.cursorX-1)
				m.score += 10
				m.createParticleBurst(m.cursorX, m.cursorY)
			}

		case tea.KeyArrowRight:
			if m.state == PlayingState {
				m.cursorX = min(m.width-1, m.cursorX+1)
				m.score += 10
				m.createParticleBurst(m.cursorX, m.cursorY)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateWindowSize()

	default:
		// Update particles
		m.particles, cmd = m.updateParticles()

		// Update background particles
		m.background = m.updateBackgroundParticles()
	}

	// Add some ambient particles when in playing state
	if m.state == PlayingState && rand.Float64() < 0.1 {
		m.particles = append(m.particles, m.createRandomParticle())
	}

	return m, cmd
}

// updateWindowSize adjusts the model when window size changes
func (m *Model) updateWindowSize() {
	// Keep cursor within bounds
	m.cursorX = min(m.cursorX, m.width-1)
	m.cursorY = min(m.cursorY, m.height-1)
}

// updateParticles updates all particles
func (m Model) updateParticles() ( []particle, tea.Cmd ) {
	newParticles := make([]particle, 0)

	for i := len(m.particles) - 1; i >= 0; i-- {
		p := m.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.life++

		// Apply gravity and friction
		p.vy += 0.1
		p.vx *= 0.99
		p.vy *= 0.99

		// Remove dead particles
		if p.life > p.maxLife {
			continue
		}

		newParticles = append(newParticles, p)
	}

	return newParticles, nil
}

// updateBackgroundParticles updates background particles
func (m *Model) updateBackgroundParticles() []colorParticle {
	for i := range m.background {
		m.background[i].x += 0.1
		if m.background[i].x > float64(m.width) {
			m.background[i].x = 0
		}
		m.background[i].y += 0.05
		if m.background[i].y > float64(m.height) {
			m.background[i].y = 0
		}
	}
	return m.background
}

// createParticleBurst creates a burst of particles at the given position
func (m *Model) createParticleBurst(x, y int) {
	for i := 0; i < 10; i++ {
		xf := float64(x) + rand.Float64()*2 - 1
		yf := float64(y) + rand.Float64()*2 - 1

		vx := rand.Float64()*4 - 2
		vy := rand.Float64()*4 - 2

		color := lipgloss.Color(fmt.Sprintf("%d", 205+rand.Intn(20)))

		m.particles = append(m.particles, particle{
			x:        xf,
			y:        yf,
			vx:       vx,
			vy:       vy,
			color:    color,
			size:     rand.Float64()*2 + 1,
			life:     0,
			maxLife:  60 + rand.Float64()*40,
		})
	}
}

// createRandomParticle creates a random particle
func (m *Model) createRandomParticle() particle {
	return particle{
		x:        rand.Float64() * float64(m.width),
		y:        rand.Float64() * float64(m.height),
		vx:       rand.Float64()*2 - 1,
		vy:       rand.Float64()*2 - 1,
		color:    lipgloss.Color(fmt.Sprintf("%d", 205+rand.Intn(30))),
		size:     rand.Float64()*3 + 1,
		life:     0,
		maxLife:  120 + rand.Float64()*80,
	}
}

// View renders the current state of the application
func (m Model) View() string {
	// Set up styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Render

	instructionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Render

	scoreStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true).
		Render

	// Render background particles
	output := ""
	for _, p := range m.background {
		output += p.color.Render(fmt.Sprintf("%s", repeatString(" ", int(p.size))))
	}

	// Render main content based on state
	switch m.state {
	case MenuState:
		title := titleStyle(m.text)
		instructions := instructionStyle(m.message)

		// Add some particles around the title
		centerX := (m.width - len(title)) / 2
		centerY := m.height / 2

		for i := 0; i < 20; i++ {
			offset := float64(i) * 0.5
			output += fmt.Sprintf("%s\n", m.renderParticleAt(centerX+offset, centerY))
		}

		output += fmt.Sprintf("\n\n%s\n\n", title)
		output += fmt.Sprintf("%s\n", instructions)

	case PlayingState:
		// Render all particles
		output = m.renderAllParticles()

		// Add cursor
		if m.cursorX < m.width && m.cursorY < m.height {
			cursorChar := "*"
			output += lipgloss.NewStyle().
				Foreground(lipgloss.Color("220")).
				Bold(true).
				Render(cursorChar)
		}

		// Add UI overlay
		ui := fmt.Sprintf("%s | %s | %s",
			scoreStyle(fmt.Sprintf("Score: %d", m.score)),
			instructionStyle("Arrow keys to move"),
			instructionStyle("SPACE to pause"))

		output += "\n\n" + ui

	case PausedState:
		// Render particles but dim them
		output = m.renderAllParticlesDimmed()

		// Add pause message
		pausedStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214")).
			Italic(true).
			Render

		output += fmt.Sprintf("\n\n%s\n", pausedStyle("PAUSED"))
		output += fmt.Sprintf("%s\n", instructionStyle("SPACE to resume, Q to quit"))

	case GameOverState:
		output = m.renderAllParticles()

		gameOverStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196")).
			Italic(true).
			Render

		output += fmt.Sprintf("\n\n%s\n", gameOverStyle("GAME OVER"))
		output += fmt.Sprintf("%s\n", instructionStyle(fmt.Sprintf("Final Score: %d", m.score)))
		output += fmt.Sprintf("%s\n", instructionStyle("Q to quit"))

	case WinState:
		output = m.renderAllParticles()

		winStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46")).
			Italic(true).
			Render

		output += fmt.Sprintf("\n\n%s\n", winStyle("YOU WIN!"))
		output += fmt.Sprintf("%s\n", instructionStyle(fmt.Sprintf("Final Score: %d", m.score)))
		output += fmt.Sprintf("%s\n", instructionStyle("Q to quit"))
	}

	return output
}

// renderAllParticles renders all particles
func (m Model) renderAllParticles() string {
	output := ""

	// Sort particles by y-coordinate (back to front)
	type particleRender struct {
		p    particle
		y    float64
	}
	particlesToRender := make([]particleRender, len(m.particles))
	for i, p := range m.particles {
		particlesToRender[i] = particleRender{p, p.y}
	}

	// Sort by y (simple bubble sort)
	for i := 0; i < len(particlesToRender); i++ {
		for j := i + 1; j < len(particlesToRender); j++ {
			if particlesToRender[i].y > particlesToRender[j].y {
				particlesToRender[i], particlesToRender[j] = particlesToRender[j], particlesToRender[i]
			}
		}
	}

	// Render sorted particles
	for _, pr := range particlesToRender {
		output += m.renderParticle(pr.p)
	}

	return output
}

// renderAllParticlesDimmed renders all particles with reduced intensity
func (m Model) renderAllParticlesDimmed() string {
	output := ""

	// Sort particles by y-coordinate
	type particleRender struct {
		p    particle
		y    float64
	}
	particlesToRender := make([]particleRender, len(m.particles))
	for i, p := range m.particles {
		particlesToRender[i] = particleRender{p, p.y}
	}

	// Sort by y
	for i := 0; i < len(particlesToRender); i++ {
		for j := i + 1; j < len(particlesToRender); j++ {
			if particlesToRender[i].y > particlesToRender[j].y {
				particlesToRender[i], particlesToRender[j] = particlesToRender[j], particlesToRender[i]
			}
		}
	}

	// Render sorted particles with dimmed colors
	for _, pr := range particlesToRender {
		dimmedColor := pr.p.color
		// Reduce color intensity
		dimmedColor = lipgloss.NewStyle().Foreground(dimmedColor).Render("")
		output += dimmedColor.Render(fmt.Sprintf("%s", repeatString(" ", int(pr.p.size))))
	}

	return output
}

// renderParticle renders a single particle
func (m Model) renderParticle(p particle) string {
	// Calculate opacity based on life
	opacity := float64(p.life) / p.maxLife

	// Get color with opacity
	color := p.color
	if opacity < 1.0 {
		// Create a more transparent version
		color = lipgloss.NewStyle().
			Foreground(color).
			Alpha(lipgloss.WithOpacity(opacity)).
			Render("")
		// Note: Simple approach - in real implementation you'd need proper alpha support
	}

	return color.Render(fmt.Sprintf("%s", repeatString(" ", int(p.size))))
}

// renderParticleAt renders a particle at absolute coordinates
func (m Model) renderParticleAt(x, y float64) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Render("●")
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func main() {
	// Set up the program
	p := tea.NewProgram(InitialModel(), tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's an error: %v", err)
		os.Exit(1)
	}
}