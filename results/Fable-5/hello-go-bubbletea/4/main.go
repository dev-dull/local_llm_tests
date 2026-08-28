// hello-wave: a small "Hello, World!" toy built on Bubble Tea.
//
// The greeting rides a rainbow sine wave over a twinkling starfield.
// Press space for confetti, ←/→ to change language, p to pause, q to quit.
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

const (
	fps       = 30
	waveAmp   = 3.0  // vertical amplitude of the wave, in rows
	waveFreq  = 0.35 // phase difference between adjacent letters
	waveSpeed = 0.15 // phase advance per frame
	gravity   = 0.06 // confetti downward acceleration, cells/frame²
)

var greetings = []struct{ lang, text string }{
	{"English", "Hello, World!"},
	{"Spanish", "¡Hola, Mundo!"},
	{"French", "Bonjour, le Monde !"},
	{"German", "Hallo, Welt!"},
	{"Italian", "Ciao, Mondo!"},
	{"Portuguese", "Olá, Mundo!"},
	{"Swedish", "Hej, Världen!"},
	{"Polish", "Witaj, Świecie!"},
}

var confettiChars = []rune("•◦*+✦✺°")

type particle struct {
	x, y   float64
	vx, vy float64
	life   int
	char   rune
	color  string
}

type star struct {
	x, y  int
	phase float64
}

type tickMsg time.Time

type model struct {
	width, height int
	frame         int
	paused        bool
	greetIdx      int
	particles     []particle
	stars         []star
	rng           *rand.Rand
}

func newModel() model {
	return model{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.scatterStars()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case " ":
			m.burst()
		case "left", "h":
			m.greetIdx = (m.greetIdx + len(greetings) - 1) % len(greetings)
		case "right", "l", "enter":
			m.greetIdx = (m.greetIdx + 1) % len(greetings)
		case "p":
			m.paused = !m.paused
		}
		return m, nil

	case tickMsg:
		if !m.paused {
			m.frame++
			m.stepParticles()
		}
		return m, tick()
	}
	return m, nil
}

// scatterStars sprinkles a fresh starfield sized to the terminal.
func (m *model) scatterStars() {
	n := m.width * m.height / 40
	m.stars = make([]star, 0, n)
	for i := 0; i < n; i++ {
		m.stars = append(m.stars, star{
			x:     m.rng.Intn(max(m.width, 1)),
			y:     m.rng.Intn(max(m.height, 1)),
			phase: m.rng.Float64() * 2 * math.Pi,
		})
	}
}

// burst launches a ring of confetti from the middle of the greeting.
func (m *model) burst() {
	cx, cy := float64(m.width)/2, float64(m.height)/2
	for i := 0; i < 45; i++ {
		angle := m.rng.Float64() * 2 * math.Pi
		speed := 0.4 + m.rng.Float64()*1.1
		m.particles = append(m.particles, particle{
			x:     cx,
			y:     cy,
			vx:    math.Cos(angle) * speed * 2, // cells are ~2x taller than wide
			vy:    math.Sin(angle) * speed,
			life:  25 + m.rng.Intn(30),
			char:  confettiChars[m.rng.Intn(len(confettiChars))],
			color: hslToHex(m.rng.Float64()*360, 1, 0.65),
		})
	}
}

func (m *model) stepParticles() {
	alive := m.particles[:0]
	for _, p := range m.particles {
		p.x += p.vx
		p.y += p.vy
		p.vy += gravity
		p.life--
		if p.life > 0 && p.y < float64(m.height) {
			alive = append(alive, p)
		}
	}
	m.particles = alive
}

type cell struct {
	r     rune
	color string
	bold  bool
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "warming up..."
	}

	grid := make([][]cell, m.height)
	for y := range grid {
		grid[y] = make([]cell, m.width)
		for x := range grid[y] {
			grid[y][x] = cell{r: ' '}
		}
	}

	m.drawStars(grid)
	m.drawParticles(grid)
	m.drawGreeting(grid)
	m.drawFooter(grid)

	return renderGrid(grid)
}

func (m model) drawStars(grid [][]cell) {
	for _, s := range m.stars {
		if s.y >= m.height || s.x >= m.width {
			continue
		}
		twinkle := (math.Sin(float64(m.frame)*0.1+s.phase) + 1) / 2
		ch, shade := '·', "#3a3a5c"
		if twinkle > 0.85 {
			ch, shade = '✦', "#9a9ac8"
		} else if twinkle > 0.5 {
			shade = "#5c5c8a"
		}
		grid[s.y][s.x] = cell{r: ch, color: shade}
	}
}

func (m model) drawParticles(grid [][]cell) {
	for _, p := range m.particles {
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < m.width && y >= 0 && y < m.height {
			grid[y][x] = cell{r: p.char, color: p.color, bold: true}
		}
	}
}

func (m model) drawGreeting(grid [][]cell) {
	text := []rune(greetings[m.greetIdx].text)
	startX := (m.width - len(text)) / 2
	baseY := m.height / 2
	for i, r := range text {
		x := startX + i
		if x < 0 || x >= m.width {
			continue
		}
		offset := math.Sin(float64(m.frame)*waveSpeed+float64(i)*waveFreq) * waveAmp
		y := baseY - int(math.Round(offset))
		if y < 0 || y >= m.height {
			continue
		}
		hue := math.Mod(float64(m.frame)*3+float64(i)*12, 360)
		grid[y][x] = cell{r: r, color: hslToHex(hue, 1, 0.62), bold: true}
	}
}

func (m model) drawFooter(grid [][]cell) {
	state := ""
	if m.paused {
		state = "  ⏸ paused"
	}
	help := fmt.Sprintf("[%s]%s   space: confetti  ←/→: language  p: pause  q: quit",
		greetings[m.greetIdx].lang, state)
	y := m.height - 1
	startX := max((m.width-len([]rune(help)))/2, 0)
	for i, r := range help {
		x := startX + i
		if x >= m.width {
			break
		}
		grid[y][x] = cell{r: r, color: "#6c6c94"}
	}
}

// renderGrid turns the cell grid into styled terminal output, batching runs
// of identically-styled cells so lipgloss isn't invoked once per character.
func renderGrid(grid [][]cell) string {
	var b strings.Builder
	rows := make([]string, len(grid))
	for y, row := range grid {
		b.Reset()
		var run []rune
		var cur cell
		flush := func() {
			if len(run) == 0 {
				return
			}
			s := string(run)
			if cur.color != "" {
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(cur.color)).Bold(cur.bold)
				s = style.Render(s)
			}
			b.WriteString(s)
			run = run[:0]
		}
		for _, c := range row {
			if c.color != cur.color || c.bold != cur.bold {
				flush()
				cur = c
			}
			run = append(run, c.r)
		}
		flush()
		rows[y] = b.String()
	}
	return strings.Join(rows, "\n")
}

// hslToHex converts an HSL color (h in degrees, s and l in [0,1]) to "#rrggbb".
func hslToHex(h, s, l float64) string {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	mm := l - c/2
	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return fmt.Sprintf("#%02x%02x%02x",
		int((r+mm)*255), int((g+mm)*255), int((b+mm)*255))
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
