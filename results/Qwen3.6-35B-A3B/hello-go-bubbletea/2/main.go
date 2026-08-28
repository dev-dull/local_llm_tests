// Package main implements a fun, animated "Hello, World!" using bubbletea.
// A rainbow-colored "Hello, World!" bounces around the terminal.
// Arrow keys steer the direction, space toggles pause, q quits.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const text = "Hello, World!"

// rainbowColor returns a lipgloss foreground color for character at position i.
// Characters are colored like a rainbow: R, O, Y, G, C, I, V cycling.
func rainbowColor(i int) lipgloss.Color {
	hues := []float64{0.0, 0.1, 0.17, 0.33, 0.5, 0.65, 0.75}
	h := hues[i%len(hues)]
	r, g, b := colorfulColorFromHSL(h, 1.0, 0.55)
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", int(r), int(g), int(b)))
}

// --- Minimal HSL→RGB (no external colorful dependency needed) ---

func colorfulColorFromHSL(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		return l, l, l
	}
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
	return r * 255, g * 255, b * 255
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2 {
		return q
	}
	if t < 2.0/3 {
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}

// --- Model ---

// model holds the entire application state.
type model struct {
	width    int
	height   int
	x, y     int       // position of the bouncing box
	dx, dy   int       // direction (+1 or -1)
	paused   bool      // is animation paused?
	phase    float64   // animation phase for border color cycling
	tick     int       // tick counter
	progress progress.Model
}

// initModel creates the initial model with defaults.
func initModel() model {
	return model{
		x:        10,
		y:        5,
		dx:       1,
		dy:       1,
		phase:    0,
		tick:     0,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

// Init returns the initial tea command.
func (m model) Init() tea.Cmd {
	return tickCmd()
}

// Update handles user input and game logic.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Clamp position inside viewport (margin 1)
		boxW := len(text) + 2
		boxH := 3
		if m.x < 1 {
			m.x = 1
		}
		if m.x > m.width-boxW-1 {
			m.x = m.width - boxW - 1
		}
		if m.y < 1 {
			m.y = 1
		}
		if m.y > m.height-boxH-1 {
			m.y = m.height - boxH - 1
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			m.dy = -1
			return m, tickCmd()
		case "down":
			m.dy = 1
			return m, tickCmd()
		case "left":
			m.dx = -1
			return m, tickCmd()
		case "right":
			m.dx = 1
			return m, tickCmd()
		case " ":
			m.paused = !m.paused
			if m.paused {
				return m, nil
			}
			return m, tickCmd()
		}

	case tickMsg:
		if m.paused {
			return m, tickCmd()
		}
		m.tick++
		m.phase += 0.15

		// Bounce the box
		m.x += m.dx
		m.y += m.dy

		boxW := len(text) + 2
		boxH := 3

		// Horizontal bounce
		if m.x <= 1 || m.x >= m.width-boxW-1 {
			m.dx = -m.dx
			m.x = clamp(m.x, 1, m.width-boxW-1)
		}
		// Vertical bounce
		if m.y <= 1 || m.y >= m.height-boxH-1 {
			m.dy = -m.dy
			m.y = clamp(m.y, 1, m.height-boxH-1)
		}

		return m, tickCmd()
	}

	return m, nil
}

// View renders the UI.
func (m model) View() string {
	var b strings.Builder

	// Header
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("231")).
		Render("🎈 Bouncing Hello, World! 🎈")
	centered := strings.Repeat(" ", max(0, (m.width-len(title))/2)) + title
	b.WriteString(centered + "\n\n")

	// Build the rainbow-colored text
	var rainbowText string
	for i, ch := range text {
		col := rainbowColor(i)
		rainbowText += lipgloss.NewStyle().Foreground(col).Render(string(ch))
	}

	// Animated border color (cycles through rainbow)
	r := int(128 + 127*math.Sin(m.phase))
	g := int(128 + 127*math.Sin(m.phase+2.094)) // +120°
	bc := int(128 + 127*math.Sin(m.phase+4.189)) // +240°
	borderCol := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, bc))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(borderCol).
		Padding(0, 1)

	box := boxStyle.Render(rainbowText)

	// Render with positioning
	// Top margin
	for i := 0; i < m.y; i++ {
		b.WriteString("\n")
	}
	b.WriteString(strings.Repeat(" ", m.x))
	b.WriteString(box)

	// Spacer
	b.WriteString("\n\n")

	// Status indicator
	status := "running"
	if m.paused {
		status = "paused"
	}
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(fmt.Sprintf("● %s", status))
	if m.paused {
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render(fmt.Sprintf("● %s", status))
	}

	// Progress bar for visual flair
	progressStr := m.progress.ViewAs(float64(m.tick%100) / 100.0)

	// Footer
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Render("↑↓←→ steer  •  space pause  •  q quit")

	b.WriteString(strings.Repeat(" ", max(0, (m.width-len(statusStyle))/2)))
	b.WriteString(statusStyle)
	b.WriteString("\n\n")
	b.WriteString(progressStr)
	b.WriteString("\n\n")
	b.WriteString(strings.Repeat(" ", max(0, (m.width-len(footer))/2)))
	b.WriteString(footer)

	return b.String()
}

// clamp keeps v between lo and hi.
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// --- Tick ---

// tickMsg fires every frame to drive animation.
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(50 * time.Millisecond)
		return tickMsg{}
	}
}

func main() {
	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}
