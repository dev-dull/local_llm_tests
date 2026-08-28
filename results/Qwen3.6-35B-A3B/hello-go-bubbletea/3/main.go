// hello-bubbletea — a starfield "Hello, World!" with typewriter reveal,
// floating rainbow letters, and emoji particle bursts.
//
// Controls:
//   q        quit
//   space    burst emoji particles
//   ←→       tilt wind left / right
//   ↑↓       shift text up / down
//   g        toggle glow mode
//
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

// ── constants ────────────────────────────────────────────────────

const greeting = "Hello, World!"

const (
	typeInterval = 120 * time.Millisecond
	tickInterval = 50 * time.Millisecond
)

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// ── types ────────────────────────────────────────────────────────

type state int

const (
	stateTyping state = iota
	stateFloating
)

// sparkle is a static star in the background.
type sparkle struct {
	x, y   int
	ch     rune
	bright bool
}

// particle is a floating emoji after a burst.
type particle struct {
	x, y    int
	vx, vy  float64
	alpha   float64
	ch      rune
}

// floatChar holds per-character state during the floating phase.
type floatChar struct {
	phase     float64 // random offset
	bobSpeed  float64
	bobAmount float64
}

type model struct {
	width, height int
	state         state
	typingIndex   int
	tick          int
	wind          float64 // -1 to 1
	textOffsetY   int
	glow          bool
	sparkles      []sparkle
	particles     []particle
	floatChars    []floatChar
}

// ── init ─────────────────────────────────────────────────────────

func newModel() model {
	m := model{
		state:    stateTyping,
		sparkles: make([]sparkle, 0, 200),
		floatChars: make([]floatChar, len(greeting)),
	}
	// Seed sparkles
	for i := 0; i < 200; i++ {
		charset := []rune("·*:+xX$@")
		m.sparkles = append(m.sparkles, sparkle{
			x:      randRange(0, 80),
			y:      randRange(0, 24),
			ch:     charset[randRange(0, len(charset))],
			bright: randRange(0, 10) == 0,
		})
	}
	// Seed random float phases
	for i := range greeting {
		m.floatChars[i] = floatChar{
			phase:     math.Pi * 2 * float64(i) / float64(len(greeting)),
			bobSpeed:  0.8 + rng.Float64()*1.2,
			bobAmount: 1.0 + rng.Float64()*2.0,
		}
	}
	return m
}

// ── update ───────────────────────────────────────────────────────

func (m *model) Init() tea.Cmd {
	return tea.Batch(m.typeTick(), m.gameTick())
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case " ":
			m.burstParticles()
			return m, nil
		case "left":
			m.wind = math.Max(-1, m.wind-0.3)
			return m, nil
		case "right":
			m.wind = math.Min(1, m.wind+0.3)
			return m, nil
		case "up":
			if m.textOffsetY > -5 {
				m.textOffsetY--
			}
			return m, nil
		case "down":
			if m.textOffsetY < 5 {
				m.textOffsetY++
			}
			return m, nil
		case "g":
			m.glow = !m.glow
			return m, nil
		}

	case typeTickMsg:
		if m.state == stateTyping {
			m.typingIndex++
			if m.typingIndex >= len(greeting) {
				m.state = stateFloating
				return m, nil
			}
			return m, m.typeTick()
		}

	case gameTickMsg:
		m.tick++
		m.twinkleSparkles()
		m.updateParticles()
		return m, m.gameTick()
	}

	return m, nil
}

// ── view ─────────────────────────────────────────────────────────

func (m *model) View() string {
	var lines []string

	// ── title bar ──
	phase := ""
	switch m.state {
	case stateTyping:
		phase = "typing…"
	case stateFloating:
		phase = "floating"
	}
	glowTag := ""
	if m.glow {
		glowTag = " ✨glow"
	}
	info := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Width(m.width).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("[%s%s]  ←→wind ↑↓shift space=burst g=glow  |  q=quit", phase, glowTag))
	lines = append(lines, info)
	lines = append(lines, "")

	// ── starfield ──
	centerY := m.height / 2
	textCenterY := centerY + m.textOffsetY

	// Build grid
	grid := make([][]rune, m.height)
	for y := range grid {
		grid[y] = make([]rune, m.width)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}

	// Place sparkles
	for _, s := range m.sparkles {
		if s.x >= 0 && s.x < m.width && s.y >= 0 && s.y < m.height {
			if s.bright {
				grid[s.y][s.x] = []rune("·*:+xX$@")[randRange(0, 8)]
			} else {
				grid[s.y][s.x] = '.'
			}
		}
	}

	// Render starfield lines
	for y := 0; y < m.height; y++ {
		lines = append(lines, string(grid[y]))
	}

	// ── greeting ──
	if m.state == stateTyping {
		lines = append(lines, "")
		lines = append(lines, m.renderTypewriter())
	} else {
		lines = append(lines, "")
		lines = append(lines, m.renderFloating(textCenterY))
	}
	lines = append(lines, "")

	// ── footer ──
	sep := strings.Repeat("─", min(m.width-2, 80))
	foot := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Width(m.width).
		Align(lipgloss.Center).
		Render(sep)
	lines = append(lines, foot)

	return strings.Join(lines, "\n")
}

// ── rendering ────────────────────────────────────────────────────

func (m *model) renderTypewriter() string {
	var parts []string
	for i := 0; i < m.typingIndex && i < len(greeting); i++ {
		c := rainbowStyle(float64(i) / float64(len(greeting)))
		parts = append(parts, c.Render(string(rune(greeting[i]))))
	}
	// Cursor
	c := rainbowStyle(float64(m.typingIndex) / float64(len(greeting)))
	parts = append(parts, c.Render("█"))
	return lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(strings.Join(parts, ""))
}

func (m *model) renderFloating(centerY int) string {
	// For each row, check which characters are visible
	rows := make([]string, m.height)
	charWidth := m.width / len(greeting)

	// First pass: compute character positions
	charPositions := make([]struct{ x, y int }, len(greeting))
	for i := range greeting {
		bobY := int(m.sinOffset(i))
		charY := centerY + bobY + (i % 3)
		charX := charWidth*i + charWidth/2
		charPositions[i] = struct{ x, y int }{x: charX, y: charY}
	}

	// Second pass: build each row
	for y := 0; y < m.height; y++ {
		var chars []string
		for i := range greeting {
			pos := charPositions[i]
			if pos.y == y {
				c := rainbowStyle(float64(i) / float64(len(greeting)))
				if m.glow {
					c = c.Bold(true)
				}
				rendered := c.Render(string(rune(greeting[i])))
				// Pad to charWidth
				for len(rendered) < charWidth {
					rendered += " "
				}
				chars = append(chars, rendered)
			} else {
				chars = append(chars, strings.Repeat(" ", charWidth))
			}
		}
		rows[y] = strings.Join(chars, "")
	}
	return lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(strings.Join(rows, "\n"))
}

func (m *model) sinOffset(i int) float64 {
	fc := m.floatChars[i]
	return math.Sin(float64(m.tick)*0.05+fc.phase+math.Sin(float64(m.tick)*0.02)*m.wind*math.Pi) * fc.bobAmount
}

// ── effects ──────────────────────────────────────────────────────

func (m *model) burstParticles() {
	emojis := []rune("🎉✨🌟💫🎊🎆🎇⭐🔥")
	count := 30
	for i := 0; i < count; i++ {
		m.particles = append(m.particles, particle{
			x:   m.width / 2,
			y:   m.height / 2,
			vx:  (rng.Float64() - 0.5) * 4,
			vy:  (rng.Float64() - 0.5) * 4 - 1,
			alpha: 1.0,
			ch:    emojis[randRange(0, len(emojis))],
		})
	}
}

func (m *model) updateParticles() {
	for i := len(m.particles) - 1; i >= 0; i-- {
		p := &m.particles[i]
		p.x += int(p.vx)
		p.y += int(p.vy)
		p.vy += 0.05 // gravity
		p.alpha -= 0.015
		if p.alpha <= 0 {
			m.particles = append(m.particles[:i], m.particles[i+1:]...)
		}
	}
}

func (m *model) twinkleSparkles() {
	if randRange(0, 20) == 0 && len(m.sparkles) > 0 {
		idx := randRange(0, len(m.sparkles))
		m.sparkles[idx].bright = !m.sparkles[idx].bright
	}
}

// ── commands ─────────────────────────────────────────────────────

type typeTickMsg struct{}
type gameTickMsg struct{}

func (m *model) typeTick() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(typeInterval)
		return typeTickMsg{}
	}
}

func (m *model) gameTick() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(tickInterval)
		return gameTickMsg{}
	}
}

// ── helpers ──────────────────────────────────────────────────────

func rainbowStyle(t float64) lipgloss.Style {
	idx := int(math.Floor(t*16)) % 16
	colors := []string{
		"229", "220", "208", "203", "213", "219",
		"141", "46", "50", "39", "57", "147",
		"199", "197", "206", "217",
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colors[idx]))
}

func randRange(lo, hi int) int {
	return lo + rng.Intn(hi-lo)
}

// ── main ─────────────────────────────────────────────────────────

func main() {
	m := newModel()
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "oops: %v\n", err)
	}
}
