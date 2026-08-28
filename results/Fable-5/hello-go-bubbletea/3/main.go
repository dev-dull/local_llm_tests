// Hello, World! — a small terminal celebration built with Bubble Tea.
//
// The greeting bounces on a sine wave, colored by a moving rainbow, over a
// field of twinkling stars. Space cycles through world languages, enter fires
// confetti, and q always quits.
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
	"github.com/mattn/go-runewidth"
)

const fps = 30

var greetings = []struct {
	text string
	lang string
}{
	{"Hello, World!", "English"},
	{"¡Hola, Mundo!", "Spanish"},
	{"Bonjour, le Monde !", "French"},
	{"Hallo, Welt!", "German"},
	{"Ciao, Mondo!", "Italian"},
	{"Olá, Mundo!", "Portuguese"},
	{"Hej, Världen!", "Swedish"},
	{"Witaj, Świecie!", "Polish"},
	{"Merhaba, Dünya!", "Turkish"},
	{"Привет, мир!", "Russian"},
	{"Γεια σου, Κόσμε!", "Greek"},
	{"שלום, עולם!", "Hebrew"},
	{"مرحبا بالعالم!", "Arabic"},
	{"नमस्ते, दुनिया!", "Hindi"},
	{"สวัสดีชาวโลก!", "Thai"},
	{"こんにちは、世界！", "Japanese"},
	{"안녕, 세상아!", "Korean"},
	{"你好，世界！", "Chinese"},
}

type star struct {
	x, y  int
	phase float64 // twinkle offset
	speed float64
}

type confetto struct {
	x, y   float64
	vx, vy float64
	life   int
	glyph  rune
	color  lipgloss.Color
}

type model struct {
	width, height int
	frame         int
	greeting      int
	stars         []star
	confetti      []confetto
	rng           *rand.Rand
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func newModel() model {
	return model{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m *model) seedStars() {
	m.stars = m.stars[:0]
	if m.width <= 0 || m.height <= 0 {
		return
	}
	n := m.width * m.height / 40
	for i := 0; i < n; i++ {
		m.stars = append(m.stars, star{
			x:     m.rng.Intn(m.width),
			y:     m.rng.Intn(m.height),
			phase: m.rng.Float64() * 2 * math.Pi,
			speed: 0.5 + m.rng.Float64()*1.5,
		})
	}
}

func (m *model) burstConfetti() {
	glyphs := []rune{'*', '•', '+', '❋', '✦', '·', 'o'}
	cx, cy := float64(m.width)/2, float64(m.height)/2
	for i := 0; i < 60; i++ {
		angle := m.rng.Float64() * 2 * math.Pi
		speed := 0.3 + m.rng.Float64()*1.2
		m.confetti = append(m.confetti, confetto{
			x:     cx,
			y:     cy,
			vx:    math.Cos(angle) * speed * 2, // terminal cells are tall; widen x
			vy:    math.Sin(angle) * speed,
			life:  20 + m.rng.Intn(25),
			glyph: glyphs[m.rng.Intn(len(glyphs))],
			color: lipgloss.Color(hsvToHex(m.rng.Float64()*360, 0.9, 1)),
		})
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.seedStars()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case " ", "right", "l":
			m.greeting = (m.greeting + 1) % len(greetings)
			m.burstConfetti()
			return m, nil
		case "left", "h":
			m.greeting = (m.greeting + len(greetings) - 1) % len(greetings)
			m.burstConfetti()
			return m, nil
		case "enter":
			m.burstConfetti()
			return m, nil
		}

	case tickMsg:
		m.frame++
		alive := m.confetti[:0]
		for _, c := range m.confetti {
			c.x += c.vx
			c.y += c.vy
			c.vy += 0.05 // gravity
			c.life--
			if c.life > 0 && c.y < float64(m.height) {
				alive = append(alive, c)
			}
		}
		m.confetti = alive
		return m, tick()
	}
	return m, nil
}

// hsvToHex converts HSV (h in degrees, s and v in [0,1]) to a #rrggbb string.
func hsvToHex(h, s, v float64) string {
	h = math.Mod(math.Mod(h, 360)+360, 360)
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
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
	off := v - c
	return fmt.Sprintf("#%02x%02x%02x",
		int((r+off)*255), int((g+off)*255), int((b+off)*255))
}

type cell struct {
	ch    rune
	style lipgloss.Style
	set   bool
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "warming up..."
	}

	canvas := make([][]cell, m.height)
	for y := range canvas {
		canvas[y] = make([]cell, m.width)
	}
	// put places ch at (x, y). A double-width rune claims the next cell too
	// (marked set with ch == 0 so the renderer emits nothing for it), keeping
	// every row exactly m.width columns wide.
	put := func(x, y int, ch rune, st lipgloss.Style) {
		w := runewidth.RuneWidth(ch)
		if x < 0 || x+w > m.width || y < 0 || y >= m.height {
			return
		}
		canvas[y][x] = cell{ch: ch, style: st, set: true}
		if w == 2 {
			canvas[y][x+1] = cell{set: true}
		}
	}

	t := float64(m.frame) / fps

	// Twinkling stars.
	for _, s := range m.stars {
		bright := (math.Sin(t*s.speed*2+s.phase) + 1) / 2
		var ch rune
		var col string
		switch {
		case bright > 0.8:
			ch, col = '✦', "#ffffff"
		case bright > 0.5:
			ch, col = '*', "#8888aa"
		case bright > 0.25:
			ch, col = '·', "#555577"
		default:
			continue
		}
		put(s.x, s.y, ch, lipgloss.NewStyle().Foreground(lipgloss.Color(col)))
	}

	// Confetti.
	for _, c := range m.confetti {
		put(int(c.x), int(c.y), c.glyph, lipgloss.NewStyle().Foreground(c.color))
	}

	// The greeting: bouncing chars with a rainbow wave.
	g := greetings[m.greeting]
	runes := []rune(g.text)
	x := (m.width - runewidth.StringWidth(g.text)) / 2
	midY := m.height / 2
	for i, r := range runes {
		bounce := int(math.Round(math.Sin(t*4+float64(i)*0.45) * 1.6))
		hue := math.Mod(t*80+float64(i)*(360/float64(len(runes))), 360)
		st := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(hsvToHex(hue, 0.85, 1)))
		put(x, midY+bounce, r, st)
		x += runewidth.RuneWidth(r)
	}

	// Language label under the greeting.
	label := "— " + g.lang + " —"
	labelX := (m.width - runewidth.StringWidth(label)) / 2
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#7d7d9c")).Italic(true)
	for i, r := range []rune(label) {
		put(labelX+i, midY+3, r, labelStyle)
	}

	// Help line.
	help := "space/←/→ language · enter confetti · q quit"
	helpX := (m.width - len([]rune(help))) / 2
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#565663"))
	for i, r := range []rune(help) {
		put(helpX+i, m.height-1, r, helpStyle)
	}

	var b strings.Builder
	for y, row := range canvas {
		for _, c := range row {
			switch {
			case c.set && c.ch != 0:
				b.WriteString(c.style.Render(string(c.ch)))
			case !c.set:
				b.WriteByte(' ')
			}
		}
		if y < m.height-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
