// hello-bubbletea: a colorful, animated "Hello, World!" TUI.
//
// Controls:
//
//	q / ctrl+c  quit (always)
//	space       toggle party mode
//	enter / →   next greeting language
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
	fps          = 30
	canvasHeight = 15
)

// A tiny 5-row block font covering just the glyphs in "HELLO, WORLD!".
var blockFont = map[rune][]string{
	'H': {"#  #", "#  #", "####", "#  #", "#  #"},
	'E': {"####", "#   ", "### ", "#   ", "####"},
	'L': {"#   ", "#   ", "#   ", "#   ", "####"},
	'O': {" ## ", "#  #", "#  #", "#  #", " ## "},
	'W': {"#   #", "#   #", "# # #", "# # #", " # # "},
	'R': {"### ", "#  #", "### ", "#  #", "#  #"},
	'D': {"### ", "#  #", "#  #", "#  #", "### "},
	',': {"  ", "  ", "  ", " #", "# "},
	'!': {"#", "#", "#", " ", "#"},
	' ': {"  ", "  ", "  ", "  ", "  "},
}

var greetings = []string{
	"Hello, World!",
	"¡Hola, Mundo!",
	"Bonjour, le Monde !",
	"Hallo, Welt!",
	"Ciao, Mondo!",
	"Olá, Mundo!",
	"Hej, Världen!",
	"Привет, мир!",
	"こんにちは、世界！",
	"안녕, 세상!",
	"Γειά σου, Κόσμε!",
	"Salve, Munde!",
}

var confettiGlyphs = []rune{'✦', '✧', '•', '○', '*', '·', '❋', '+'}

type confetto struct {
	x, y  float64
	vy    float64
	glyph rune
	hue   float64
}

type tickMsg time.Time

type model struct {
	frame    int
	width    int
	party    bool
	greeting int
	confetti []confetto
	rng      *rand.Rand
}

func newModel() model {
	return model{
		width: 80,
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case " ":
			m.party = !m.party
		case "enter", "right", "tab":
			m.greeting = (m.greeting + 1) % len(greetings)
		case "left":
			m.greeting = (m.greeting + len(greetings) - 1) % len(greetings)
		}
	case tea.WindowSizeMsg:
		if msg.Width > 0 {
			m.width = msg.Width
		}
	case tickMsg:
		m.frame++
		if m.frame%(fps*4) == 0 {
			m.greeting = (m.greeting + 1) % len(greetings)
		}
		m.stepConfetti()
		return m, tick()
	}
	return m, nil
}

func (m *model) stepConfetti() {
	spawnChance := 0.35
	if m.party {
		spawnChance = 1.6
	}
	for spawnChance > 0 {
		if spawnChance >= 1 || m.rng.Float64() < spawnChance {
			m.confetti = append(m.confetti, confetto{
				x:     m.rng.Float64() * float64(m.canvasWidth()),
				y:     0,
				vy:    0.15 + m.rng.Float64()*0.25,
				glyph: confettiGlyphs[m.rng.Intn(len(confettiGlyphs))],
				hue:   m.rng.Float64() * 360,
			})
		}
		spawnChance--
	}
	speed := 1.0
	if m.party {
		speed = 1.8
	}
	alive := m.confetti[:0]
	for _, c := range m.confetti {
		c.y += c.vy * speed
		c.x += math.Sin(c.y*1.3) * 0.3
		if c.y < canvasHeight {
			alive = append(alive, c)
		}
	}
	m.confetti = alive
}

func (m model) canvasWidth() int {
	w := m.width
	if w < 20 {
		w = 20
	}
	if w > 100 {
		w = 100
	}
	return w
}

type cell struct {
	ch rune
	fg string
}

func (m model) View() string {
	w := m.canvasWidth()

	canvas := make([][]cell, canvasHeight)
	for y := range canvas {
		canvas[y] = make([]cell, w)
		for x := range canvas[y] {
			canvas[y][x] = cell{ch: ' '}
		}
	}

	// Confetti behind the banner.
	for _, c := range m.confetti {
		x, y := int(c.x), int(c.y)
		if x >= 0 && x < w && y >= 0 && y < canvasHeight {
			canvas[y][x] = cell{ch: c.glyph, fg: hslHex(c.hue, 0.7, 0.6)}
		}
	}

	m.drawBanner(canvas, w)

	var b strings.Builder
	for _, row := range canvas {
		for _, c := range row {
			if c.ch == ' ' {
				b.WriteByte(' ')
			} else if c.fg != "" {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(c.fg)).Render(string(c.ch)))
			} else {
				b.WriteRune(c.ch)
			}
		}
		b.WriteByte('\n')
	}

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(hslHex(float64(m.frame*3%360), 0.8, 0.7))).
		Bold(true).
		Render(greetings[m.greeting])

	mode := "space: party mode"
	if m.party {
		mode = "space: calm down 🎉"
	}
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).
		Render(fmt.Sprintf("enter/←/→: language   %s   q: quit", mode))

	center := lipgloss.NewStyle().Width(w).Align(lipgloss.Center)
	b.WriteString(center.Render(subtitle))
	b.WriteString("\n\n")
	b.WriteString(center.Render(help))
	b.WriteByte('\n')
	return b.String()
}

// drawBanner renders "HELLO, WORLD!" in block letters, each column riding a
// sine wave and colored by a scrolling rainbow.
func (m model) drawBanner(canvas [][]cell, w int) {
	const text = "HELLO, WORLD!"

	bannerW := 0
	for _, r := range text {
		bannerW += len([]rune(blockFont[r][0])) + 1
	}
	bannerW--

	left := (w - bannerW) / 2
	if left < 0 {
		left = 0
	}
	top := (canvasHeight - 5) / 2

	speed := 1.0
	if m.party {
		speed = 2.5
	}
	phase := float64(m.frame) * speed

	x := left
	for _, r := range text {
		glyph := blockFont[r]
		glyphW := len([]rune(glyph[0]))
		for col := 0; col < glyphW; col++ {
			gx := x + col
			if gx < 0 || gx >= w {
				continue
			}
			wave := int(math.Round(math.Sin(phase/8+float64(gx)/4) * 1.5))
			hue := math.Mod(float64(gx)*6+phase*2, 360)
			for row := 0; row < 5; row++ {
				if []rune(glyph[row])[col] != '#' {
					continue
				}
				gy := top + row + wave
				if gy >= 0 && gy < canvasHeight {
					canvas[gy][gx] = cell{ch: '█', fg: hslHex(hue, 0.9, 0.65)}
				}
			}
		}
		x += glyphW + 1
	}
}

// hslHex converts HSL (h in degrees, s/l in [0,1]) to a #rrggbb hex string.
func hslHex(h, s, l float64) string {
	h = math.Mod(math.Mod(h, 360)+360, 360) / 360
	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		r = hueToRGB(p, q, h+1.0/3)
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h-1.0/3)
	}
	return fmt.Sprintf("#%02x%02x%02x", int(r*255), int(g*255), int(b*255))
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	switch {
	case t < 1.0/6:
		return p + (q-p)*6*t
	case t < 1.0/2:
		return q
	case t < 2.0/3:
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("Goodbye, World! 👋")
}
