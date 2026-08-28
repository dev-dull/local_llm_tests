package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Types ────────────────────────────────────────────────────────────────────

type tickMsg time.Time

type particle struct {
	x, y     float64
	vx, vy   float64
	life     float64
	maxLife  float64
	emoji    string
}

type model struct {
	width, height   int
	x, y            float64 // ball position
	vx, vy          float64 // ball velocity
	rainbow         float64 // color cycling progress
	spinner         int
	showParticles   bool
	message         string
	shakeOffset     float64
	shakeDecay      float64
	particles       []particle
	greetingPhase   int // 0=entrance, 1=full, 2=pulse
	entranceProgress float64
}

// ── Init ─────────────────────────────────────────────────────────────────────

func initialModel() model {
	return model{
		x:               30,
		y:               10,
		vx:              0.3,
		vy:              0.15,
		message:         "Hello, World!",
		greetingPhase:   0,
		entranceProgress: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.WindowSize(),
		tickerCmd(),
	)
}

// ── Update ───────────────────────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.x > float64(msg.Width-4) {
			m.x = float64(msg.Width - 4)
			m.vx = -m.vx
		}
		if m.y > float64(msg.Height-6) {
			m.y = float64(msg.Height - 6)
			m.vy = -m.vy
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case " ":
			m.showParticles = !m.showParticles
			if m.showParticles {
				m.shakeOffset = 5
				m.shakeDecay = 5
				m.spawnParticles()
			}
			return m, nil
		case "up", "w":
			m.vy = -0.4
		case "down", "s":
			m.vy = 0.4
		case "left", "a":
			m.vx = -0.4
		case "right", "d":
			m.vx = 0.4
		case "r":
			m.x = float64(m.width / 2)
			m.y = float64(m.height / 2)
			m.vx = 0.3 * (1 + rand.Float64())
			m.vy = 0.15 * (1 + rand.Float64())
			m.showParticles = false
			m.particles = nil
			return m, nil
		}

	case tickMsg:
		m.tick()
		return m, tickerCmd()
	}

	return m, nil
}

// ── Animation tick ───────────────────────────────────────────────────────────

func (m *model) tick() {
	m.rainbow += 0.05
	m.spinner = (m.spinner + 1) % 8

	if m.shakeDecay > 0 {
		m.shakeOffset *= 0.9
		m.shakeDecay -= 0.1
	}

	// Move ball
	m.x += m.vx
	m.y += m.vy

	// Bounce off walls
	if m.x <= 2 || m.x >= float64(m.width-4) {
		m.vx = -m.vx
		m.x = clampF(m.x, 2, float64(m.width-4))
	}
	if m.y <= 1 || m.y >= float64(m.height-6) {
		m.vy = -m.vy
		m.y = clampF(m.y, 1, float64(m.height-6))
	}

	// Update particles
	for i := len(m.particles) - 1; i >= 0; i-- {
		p := &m.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += 0.02
		p.life -= 0.02
		if p.life <= 0 {
			m.particles = append(m.particles[:i], m.particles[i+1:]...)
		}
	}
}

// ── Commands ─────────────────────────────────────────────────────────────────

func tickerCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg(time.Now())
	})
}

func (m *model) spawnParticles() {
	emojis := []string{"✨", "🌟", "💫", "⭐", "🎉", "🎊", "💥", "🔥"}
	for i := 0; i < 20; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := 0.5 + rand.Float64()*1.5
		m.particles = append(m.particles, particle{
			x: m.x, y: m.y,
			vx: math.Cos(angle) * speed,
			vy: math.Sin(angle) * speed - 1.0,
			life: 1.0, maxLife: 1.0,
			emoji: emojis[rand.Intn(len(emojis))],
		})
	}
}

// ── Color helpers ────────────────────────────────────────────────────────────

// hslToAnsi256 converts HSL hue (0-1) to an ANSI 256-color escape sequence.
// Saturation and lightness are fixed for a vibrant rainbow.
func hslToAnsi256(h float64) string {
	h = math.Mod(h, 1.0)
	hIdx := int(math.Floor(h * 6)) % 6
	f := h*6 - math.Floor(h*6)
	p := 0.5 * (1 - 0.8)
	q := 0.5 * (1 - f*0.8)
	t := 0.5 * (1 - (1-f)*0.8)

	var r, g, b float64
	switch hIdx {
	case 0:
		r, g, b = 0.5, t, p
	case 1:
		r, g, b = q, 0.5, p
	case 2:
		r, g, b = p, 0.5, t
	case 3:
		r, g, b = p, q, 0.5
	case 4:
		r, g, b = t, p, 0.5
	case 5:
		r, g, b = 0.5, p, q
	}

	rIdx := int(math.Round(r * 5))
	gIdx := int(math.Round(g * 5))
	bIdx := int(math.Round(b * 5))

	if rIdx < 0 {
		rIdx = 0
	}
	if rIdx > 5 {
		rIdx = 5
	}
	if gIdx < 0 {
		gIdx = 0
	}
	if gIdx > 5 {
		gIdx = 5
	}
	if bIdx < 0 {
		bIdx = 0
	}
	if bIdx > 5 {
		bIdx = 5
	}

	idx := 16 + 36*rIdx + 6*gIdx + bIdx
	return fmt.Sprintf("\033[38;5;%dm", idx)
}

func resetColor() string {
	return "\033[39m"
}

func intMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clampF(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// ── View ─────────────────────────────────────────────────────────────────────

func (m model) View() string {
	if m.height == 0 {
		return "Resizing..."
	}

	width := m.width
	if width < 20 {
		width = 20
	}

	rainbowH := math.Mod(m.rainbow, 1.0)
	color1 := hslToAnsi256(rainbowH)
	color2 := hslToAnsi256(math.Mod(rainbowH+0.33, 1.0))
	color3 := hslToAnsi256(math.Mod(rainbowH+0.66, 1.0))

	// Rainbow "Hello, World!"
	centerY := m.height / 2
	hello := m.message
	runes := []rune(hello)
	var coloredHello strings.Builder
	for i, r := range runes {
		var c string
		switch i % 3 {
		case 0:
			c = color1
		case 1:
			c = color2
		case 2:
			c = color3
		}
		coloredHello.WriteString(c + string(r))
	}
	coloredHello.WriteString(resetColor())

	bgStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(fmt.Sprintf("%d", 232+int(rainbowH*23)))).
		Width(width)
	bgLine := bgStyle.Render(strings.Repeat(" ", width))

	var sb strings.Builder

	blank := strings.Repeat(" ", width)

	// Upper area
	for i := 0; i < intMax(0, centerY-4); i++ {
		sb.WriteString(bgLine + "\n")
	}

	// Ball trail particles
	for i := 0; i < 3; i++ {
		trailY := centerY - 2 - i
		trailX := int(m.x) - 3 - i*2
		if trailY >= 0 && trailY < m.height && trailX >= 0 && trailX < width {
			line := blank[:intMin(trailX, width-1)] + "°"
			sb.WriteString(line + strings.Repeat(" ", intMax(0, width-trailX-1)) + "\n")
		}
	}

	// Ball with glow
	ballLine := blank
	ballX := int(m.x)
	ballY := int(m.y)
	if ballX > 0 && ballX+2 < width {
		ballLine = blank[:ballX] + "🫧 " + blank[ballX+3:]
	}
	if ballX > 1 {
		ballLine = blank[:ballX-1] + "·" + ballLine[ballX:]
	}
	if ballX+2 < width {
		ballLine = ballLine[:ballX+2] + "·" + ballLine[ballX+2:]
	}
	if ballY > 1 && ballY < m.height {
		sb.WriteString(ballLine + "\n")
	}

	// Particles
	for _, p := range m.particles {
		px := int(p.x)
		py := int(p.y)
		if py >= 0 && py < m.height && px >= 0 && px < width {
			line := blank[:px] + p.emoji
			sb.WriteString(line + strings.Repeat(" ", intMax(0, width-px-len(p.emoji))) + "\n")
		}
	}

	// Greeting text - centered horizontally
	helloLen := len(hello)
	padLeft := intMax(0, (width-helloLen)/2)
	sb.WriteString(blank[:padLeft] + coloredHello.String() + "\n")

	// Bottom area
	for i := 0; i < intMax(0, m.height-centerY-3); i++ {
		sb.WriteString(bgLine + "\n")
	}

	// Footer
	spinner := []string{"⠋", "⠙", "⠸", "⠴", "⠦", "⠇", "⠋", "⠙"}
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Width(width).
		Render(fmt.Sprintf(
			" %s  [space] particles  [r] reset  [q] quit",
			spinner[m.spinner],
		))

	sb.WriteString(footer)

	return sb.String()
}

// ── Main ─────────────────────────────────────────────────────────────────────

func main() {
	rand.Seed(time.Now().UnixNano())

	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}
