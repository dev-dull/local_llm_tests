// hello-go-bubbletea: a small terminal toy that says hello with style.
//
// "Hello, World!" rides a sine wave in rainbow colors while sparkles
// twinkle around it. Press enter to greet the world in another language,
// tab to change the color theme, space to freeze time, and q to leave.
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
	waveAmp   = 3    // vertical amplitude of the wave, in rows
	waveFreq  = 0.35 // phase difference between adjacent letters
	waveSpeed = 0.18 // phase advance per frame
	maxSparks = 60
)

var greetings = []string{
	"Hello, World!",
	"Hola, Mundo!",
	"Bonjour, le Monde!",
	"Hallo, Welt!",
	"Ciao, Mondo!",
	"Olá, Mundo!",
	"Witaj, Świecie!",
	"Hej, Världen!",
	"Saluton, Mondo!",
}

// A theme maps a letter index and frame to a hue range.
type theme struct {
	name    string
	hueBase float64 // starting hue
	hueSpan float64 // how far the hue travels across the text
}

var themes = []theme{
	{"rainbow", 0, 360},
	{"ocean", 170, 90},
	{"sunset", 0, 70},
	{"neon", 270, 120},
}

type sparkle struct {
	x, y    int
	life    int
	maxLife int
}

type tickMsg time.Time

type model struct {
	width    int
	height   int
	frame    int
	paused   bool
	greeting int
	theme    int
	sparks   []sparkle
	rng      *rand.Rand
}

func newModel() model {
	return model{
		width:  80,
		height: 24,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case " ":
			m.paused = !m.paused
		case "enter":
			m.greeting = (m.greeting + 1) % len(greetings)
		case "tab":
			m.theme = (m.theme + 1) % len(themes)
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tickMsg:
		if !m.paused {
			m.frame++
			m.updateSparks()
		}
		return m, tick()
	}

	return m, nil
}

func (m *model) updateSparks() {
	alive := m.sparks[:0]
	for _, s := range m.sparks {
		s.life--
		if s.life > 0 {
			alive = append(alive, s)
		}
	}
	m.sparks = alive

	for i := 0; i < 2 && len(m.sparks) < maxSparks; i++ {
		if m.rng.Float64() < 0.6 && m.width > 2 && m.height > 3 {
			life := 10 + m.rng.Intn(25)
			m.sparks = append(m.sparks, sparkle{
				x:       m.rng.Intn(m.width),
				y:       m.rng.Intn(m.height - 2), // keep off the help line
				life:    life,
				maxLife: life,
			})
		}
	}
}

// cell is one styled character on the screen grid.
type cell struct {
	ch    rune
	color string
	bold  bool
	faint bool
}

func (m model) View() string {
	w, h := m.width, m.height
	if w < 4 || h < 4 {
		return "(terminal too small — q to quit)"
	}

	grid := make([][]cell, h)
	for y := range grid {
		grid[y] = make([]cell, w)
		for x := range grid[y] {
			grid[y][x] = cell{ch: ' '}
		}
	}

	m.paintSparks(grid)
	m.paintWave(grid)
	m.paintHelp(grid)

	return renderGrid(grid)
}

func (m model) paintSparks(grid [][]cell) {
	for _, s := range m.sparks {
		if s.y < 0 || s.y >= len(grid) || s.x < 0 || s.x >= len(grid[0]) {
			continue
		}
		frac := float64(s.life) / float64(s.maxLife)
		var ch rune
		var col string
		switch {
		case frac > 0.66:
			ch, col = '✦', "#ffffdd"
		case frac > 0.33:
			ch, col = '✧', "#aaaacc"
		default:
			ch, col = '·', "#555577"
		}
		grid[s.y][s.x] = cell{ch: ch, color: col}
	}
}

func (m model) paintWave(grid [][]cell) {
	text := []rune(greetings[m.greeting])
	w, h := len(grid[0]), len(grid)
	startCol := (w - len(text)) / 2
	centerRow := h / 2
	phase := float64(m.frame) * waveSpeed
	th := themes[m.theme]

	for i, r := range text {
		x := startCol + i
		if x < 0 || x >= w {
			continue
		}
		yOff := int(math.Round(waveAmp * math.Sin(phase+float64(i)*waveFreq)))
		y := centerRow - yOff
		if y < 0 || y >= h {
			continue
		}
		if r == ' ' {
			continue // leave background sparkles visible under spaces
		}
		hue := th.hueBase + th.hueSpan*float64(i)/float64(len(text)) + float64(m.frame)*2
		grid[y][x] = cell{ch: r, color: hslToHex(math.Mod(hue, 360), 0.9, 0.65), bold: true}
	}
}

func (m model) paintHelp(grid [][]cell) {
	status := ""
	if m.paused {
		status = "  ⏸ paused"
	}
	help := fmt.Sprintf("enter: greeting • tab: theme (%s) • space: pause • q: quit%s",
		themes[m.theme].name, status)
	row := len(grid) - 1
	startCol := (len(grid[0]) - len([]rune(help))) / 2
	if startCol < 0 {
		startCol = 0
	}
	for i, r := range help {
		x := startCol + i
		if x >= len(grid[0]) {
			break
		}
		grid[row][x] = cell{ch: r, color: "#888888", faint: true}
	}
}

// renderGrid turns the cell grid into a styled string, batching unstyled
// runs so we don't emit escape codes for every blank cell.
func renderGrid(grid [][]cell) string {
	var b strings.Builder
	for y, row := range grid {
		var plain strings.Builder
		flush := func() {
			b.WriteString(plain.String())
			plain.Reset()
		}
		for _, c := range row {
			if c.color == "" {
				plain.WriteRune(c.ch)
				continue
			}
			flush()
			st := lipgloss.NewStyle().Foreground(lipgloss.Color(c.color))
			if c.bold {
				st = st.Bold(true)
			}
			if c.faint {
				st = st.Faint(true)
			}
			b.WriteString(st.Render(string(c.ch)))
		}
		flush()
		if y < len(grid)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// hslToHex converts HSL (h in degrees, s and l in [0,1]) to a #rrggbb string.
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
	fmt.Println("Goodbye! 👋")
}
