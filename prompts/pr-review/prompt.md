# Code review: six pull requests against a Bubble Tea app

You are reviewing pull requests for `hello-bubbletea`, a small terminal
animation program written in Go with `charmbracelet/bubbletea` and
`lipgloss`. The full current source of `main.go` (the only source file) is
given below, followed by six pull requests. Each PR is an independent change
against that same base file — no PR builds on another, and every PR compiles
successfully. Review each PR on its own.

For every PR, judge whether the change does what its description claims,
whether it introduces or contains any defect, and whether there is a clearly
better way to achieve the same goal. Review the change itself: mention
pre-existing problems in the base program only where the PR touches or claims
to fix them.

Format your response in markdown, exactly one section per PR:

```markdown
## PR <n>: <title>

**Verdict:** one of APPROVE | APPROVE WITH SUGGESTIONS | REQUEST CHANGES

**Findings:**
1. **[blocking|suggestion|nit]** <location> — what the problem is and why.
   (Write `none` if you have no findings.)

**Summary:** one or two sentences justifying the verdict.
```

Verdict meanings: APPROVE = correct as-is; APPROVE WITH SUGGESTIONS = correct
and mergeable, but you have improvements worth making; REQUEST CHANGES = the
PR has a defect that must be fixed before merging. Every finding must point at
a specific line or hunk of the diff.

## Base source: `main.go`

```go
// hello-bubbletea — a colorful, animated "Hello, World!" using Bubble Tea.
//
// Controls:
//   Space  cycle animation mode
//   r      reset to defaults
//   q      quit
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ──────────────────────────── types ────────────────────────────

const numModes = 4

type mode int

const (
	modeWave mode = iota // sine-wave bouncing letters
	modeRain             // letters fall from top like rain
	modeSpin             // text rotates in a circle
	modePulse            // color + size pulses
)

func (m mode) String() string {
	return [...]string{"Wave", "Rain", "Spin", "Pulse"}[m]
}

type model struct {
	width   int
	height  int
	mode    mode
	tick    int
	spun    float64 // radians for spin mode
	rainPos []rainChar
}

type rainChar struct {
	r   rune
	y   int
	v   int // velocity
	col lipgloss.Style
}

// ──────────────────────────── update ───────────────────────────

func (m *model) Init() tea.Cmd {
	return tickerCmd()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if len(m.rainPos) == 0 {
			for _, r := range "Hello, World!" {
				m.rainPos = append(m.rainPos, rainChar{
					r:   r,
					y:   -10 - randInt(20),
					v:   1 + randInt(3),
					col: rainbowStyle(0),
				})
			}
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case " ":
			m.mode = (m.mode + 1) % numModes
			return m, nil
		case "r":
			m.reset()
			return m, nil
		}
	}

	if _, ok := msg.(tickerMsg); ok {
		m.tick++
		m.spun += 0.05
		return m, tickerCmd()
	}

	return m, nil
}

func (m *model) reset() {
	m.tick = 0
	m.spun = 0
	m.mode = 0
	for i := range m.rainPos {
		m.rainPos[i].y = -10 - randInt(20)
	}
}

// ──────────────────────────── view ─────────────────────────────

func (m *model) View() string {
	var lines []string

	// Title bar
	info := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Width(m.width).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("Mode: %s  |  Space=cycle  |  R=reset  |  Q=quit", m.mode))
	lines = append(lines, info)
	lines = append(lines, "")

	// Render the greeting in the active mode
	switch m.mode {
	case modeWave:
		lines = append(lines, renderWave(m)...)
	case modeRain:
		// Update rain positions in-place so they persist across frames
		for i := range m.rainPos {
			m.rainPos[i].y += m.rainPos[i].v
			m.rainPos[i].col = rainbowStyle(float64(i) / float64(len(m.rainPos)))
			if m.rainPos[i].y > m.height+5 {
				m.rainPos[i].y = -5
			}
		}
		lines = append(lines, renderRainStatic(m)...)
	case modeSpin:
		lines = append(lines, renderSpin(m)...)
	case modePulse:
		lines = append(lines, renderPulse(m)...)
	}

	// Footer
	sep := strings.Repeat("─", minInt(m.width, 80))
	foot := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Width(m.width).
		Align(lipgloss.Center).
		Render(sep)
	lines = append(lines, foot)

	return strings.Join(lines, "\n")
}

// ── mode: wave ──────────────────────────────────────────────────

func renderWave(m *model) []string {
	text := "Hello, World!"
	centerY := m.height / 2

	// Build main line with rainbow colors
	var chars []string
	for i, r := range text {
		col := rainbowStyle(float64(i) / float64(len(text)))
		chars = append(chars, col.Render(string(r)))
	}
	mainRow := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(strings.Join(chars, ""))

	// Shadow/glow row (dimmer, offset down by 1)
	var shadow []string
	for _, r := range text {
		col := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Render(string(r))
		shadow = append(shadow, col)
	}
	shadowRow := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width).Render(strings.Join(shadow, ""))

	// Vertical wave motion
	avgOffset := math.Sin(float64(m.tick)*0.15) * 2
	lineY := centerY + int(avgOffset)

	result := make([]string, m.height)
	for y := 0; y < lineY && y < m.height; y++ {
		result[y] = ""
	}
	if lineY >= 0 && lineY < m.height {
		result[lineY] = mainRow
	}
	if lineY+1 >= 0 && lineY+1 < m.height {
		result[lineY+1] = shadowRow
	}
	for y := lineY + 2; y < m.height; y++ {
		result[y] = ""
	}

	return result
}

// ── mode: rain ──────────────────────────────────────────────────

func renderRainStatic(m *model) []string {
	// Build each row based on current rain positions
	result := make([]string, m.height)
	for y := 0; y < m.height; y++ {
		var chars []string
		spacing := m.width / len(m.rainPos)
		for i := range m.rainPos {
			if m.rainPos[i].y == y {
				chars = append(chars, m.rainPos[i].col.Render(string(m.rainPos[i].r)))
			} else if m.rainPos[i].y == y-1 && m.rainPos[i].y >= 0 {
				// Trail character (dimmer)
				chars = append(chars, m.rainPos[i].col.Foreground(lipgloss.Color("8")).Render(string(m.rainPos[i].r)))
			} else {
				chars = append(chars, strings.Repeat(" ", spacing))
			}
		}
		result[y] = strings.Join(chars, "")
	}

	return result
}

// ── mode: spin ──────────────────────────────────────────────────

func renderSpin(m *model) []string {
	text := "Hello, World!"
	size := len(text)
	cx := m.width / 2
	radius := int(math.Min(float64(minInt(m.width, 60)/2), float64(size)))
	if radius < 2 {
		radius = 2
	}

	// Build a grid with text arranged in a rotating circle
	rows := 2*radius + 1
	result := make([]string, rows)
	for dy := -radius; dy <= radius; dy++ {
		var row []rune
		for dx := -radius; dx <= radius; dx++ {
			dist := math.Sqrt(float64(dx*dx + dy*dy))
			if dist >= float64(radius)-1.2 && dist <= float64(radius)+1.2 {
				angle := math.Atan2(float64(dy), float64(dx))
				charIdx := int(math.Mod((angle+m.spun)/(2*math.Pi)*float64(size)+float64(size), float64(size)))
				col := rainbowStyle(float64(charIdx) / float64(size))
				row = append(row, []rune(col.Render(string(text[charIdx%size])))...)
			} else {
				row = append(row, ' ')
			}
		}
		// Center each row horizontally
		line := string(row)
		padLeft := strings.Repeat(" ", maxInt(0, cx-radius))
		result[dy+radius] = padLeft + line
	}

	return result
}

// ── mode: pulse ─────────────────────────────────────────────────

func renderPulse(m *model) []string {
	text := "Hello, World!"
	t := float64(m.tick) * 0.1
	scale := 1.0 + 0.3*math.Sin(t)
	alpha := 0.5 + 0.5*math.Sin(t)

	// Build colored chars
	var rendered []string
	for i, r := range text {
		col := rainbowStyle(float64(i) / float64(len(text)))
		rendered = append(rendered, col.Render(string(r)))
	}

	mainLine := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Bold(math.Abs(math.Sin(t)) > 0.5).
		Italic(math.Cos(t) > 0).
		Render(strings.Join(rendered, ""))

	// Glow row
	glowVal := 9 + 15*int(15*math.Max(0, math.Sin(t)))/15
	glowCol := lipgloss.NewStyle().
		Foreground(lipgloss.Color(fmt.Sprintf("%d", glowVal))).
		Align(lipgloss.Center).
		Width(m.width).
		Render(strings.Join(rendered, ""))

	statusLine := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Foreground(lipgloss.Color("10")).
		Render(fmt.Sprintf("pulse: %.2f  scale: %.2f", alpha, scale))

	result := make([]string, m.height)
	cy := m.height / 2
	for y := 0; y < cy-2 && y >= 0; y++ {
		result[y] = ""
	}
	if cy-2 >= 0 && cy-2 < m.height {
		result[cy-2] = ""
	}
	if cy-1 >= 0 && cy-1 < m.height {
		result[cy-1] = mainLine
	}
	if cy >= 0 && cy < m.height {
		result[cy] = glowCol
	}
	if cy+1 >= 0 && cy+1 < m.height {
		result[cy+1] = statusLine
	}
	for y := cy + 2; y < m.height; y++ {
		result[y] = ""
	}

	return result
}

// ──────────────────────────── helpers ────────────────────────────

func rainbowStyle(t float64) lipgloss.Style {
	idx := int(math.Floor(t*15)) % 16
	colors := []string{
		"229", "220", "208", "203", "213", "219",
		"141", "46", "50", "39", "57", "147",
		"199", "197", "206", "217",
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(colors[idx]))
}

// tickerMsg fires each frame to animate.
type tickerMsg struct{}

func tickerCmd() tea.Cmd {
	return func() tea.Msg {
		return tickerMsg{}
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// randInt returns a pseudo-random int in [0, max).
func randInt(max int) int {
	return (int(time.Now().UnixNano()) / 1e6) % max
}

// ──────────────────────────── main ─────────────────────────────

func main() {
	m := &model{}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "oops: %v\n", err)
	}
}
```

## PR 1: Fix 100% CPU usage during animation

> The app pegs a CPU core while running. Cap the frame rate at roughly 30 fps so the animation loop stops busy-spinning.

```diff
--- a/main.go
+++ b/main.go
@@ -330,6 +330,7 @@
 
 func tickerCmd() tea.Cmd {
 	return func() tea.Msg {
+		time.Sleep(33 * time.Millisecond)
 		return tickerMsg{}
 	}
 }
```

## PR 2: Use tea.Tick for frame scheduling

> The animation ticker returns a message immediately with no delay, so the event loop spins as fast as it can. Schedule frames with Bubble Tea's tick mechanism at ~30 fps instead.

```diff
--- a/main.go
+++ b/main.go
@@ -329,9 +329,9 @@
 type tickerMsg struct{}
 
 func tickerCmd() tea.Cmd {
-	return func() tea.Msg {
+	return tea.Tick(time.Second/30, func(time.Time) tea.Msg {
 		return tickerMsg{}
-	}
+	})
 }
 
 func maxInt(a, b int) int {
```

## PR 3: Add +/- keys to control spin speed

> Let the user speed up or slow down the Spin animation with + and -. The speed is clamped to a sensible minimum, reset restores the default, and the help bar documents the new keys.

```diff
--- a/main.go
+++ b/main.go
@@ -35,12 +35,13 @@
 }
 
 type model struct {
-	width   int
-	height  int
-	mode    mode
-	tick    int
-	spun    float64 // radians for spin mode
-	rainPos []rainChar
+	width     int
+	height    int
+	mode      mode
+	tick      int
+	spun      float64 // radians for spin mode
+	spinSpeed float64 // radians advanced per tick in spin mode
+	rainPos   []rainChar
 }
 
 type rainChar struct {
@@ -83,6 +84,12 @@
 		case "r":
 			m.reset()
 			return m, nil
+		case "+", "=":
+			m.spinSpeed += 0.02
+			return m, nil
+		case "-":
+			m.spinSpeed = math.Max(0.01, m.spinSpeed-0.02)
+			return m, nil
 		}
 	}
 
@@ -98,6 +105,7 @@
 func (m *model) reset() {
 	m.tick = 0
 	m.spun = 0
+	m.spinSpeed = 0.05
 	m.mode = 0
 	for i := range m.rainPos {
 		m.rainPos[i].y = -10 - randInt(20)
@@ -114,7 +122,7 @@
 		Foreground(lipgloss.Color("14")).
 		Width(m.width).
 		Align(lipgloss.Center).
-		Render(fmt.Sprintf("Mode: %s  |  Space=cycle  |  R=reset  |  Q=quit", m.mode))
+		Render(fmt.Sprintf("Mode: %s  |  Space=cycle  |  +/-=spin speed  |  R=reset  |  Q=quit", m.mode))
 	lines = append(lines, info)
 	lines = append(lines, "")
 
@@ -356,7 +364,7 @@
 // ──────────────────────────── main ─────────────────────────────
 
 func main() {
-	m := &model{}
+	m := &model{spinSpeed: 0.05}
 	p := tea.NewProgram(m, tea.WithAltScreen())
 	if _, err := p.Run(); err != nil {
 		fmt.Fprintf(os.Stderr, "oops: %v\n", err)
```

## PR 4: Fix correlated random values in rain initialization

> randInt derived its value from the current wall-clock millisecond, so consecutive calls returned strongly correlated values and the rain started in lockstep. Use a real PRNG.

```diff
--- a/main.go
+++ b/main.go
@@ -9,9 +9,9 @@
 import (
 	"fmt"
 	"math"
+	"math/rand"
 	"os"
 	"strings"
-	"time"
 
 	tea "github.com/charmbracelet/bubbletea"
 	"github.com/charmbracelet/lipgloss"
@@ -350,7 +350,7 @@
 
 // randInt returns a pseudo-random int in [0, max).
 func randInt(max int) int {
-	return (int(time.Now().UnixNano()) / 1e6) % max
+	return rand.Intn(max)
 }
 
 // ──────────────────────────── main ─────────────────────────────
```

## PR 5: Fix header being pushed off-screen by tall views

> The View can emit more lines than the terminal has rows, which scrolls the header (with the keybinding help) out of sight. Clamp the output to the window height so the header always stays visible.

```diff
--- a/main.go
+++ b/main.go
@@ -146,6 +146,10 @@
 		Align(lipgloss.Center).
 		Render(sep)
 	lines = append(lines, foot)
+
+	if m.height > 0 && len(lines) > m.height {
+		lines = lines[:m.height]
+	}
 
 	return strings.Join(lines, "\n")
 }
```

## PR 6: Move rain state updates out of View() and into Update()

> View() must be a pure function of the model, but the rain animation was advancing positions inside View(). Move the per-frame rain updates into Update()'s tick handling where they belong.

```diff
--- a/main.go
+++ b/main.go
@@ -89,6 +89,15 @@
 	if _, ok := msg.(tickerMsg); ok {
 		m.tick++
 		m.spun += 0.05
+		if m.mode == modeRain {
+			for i, c := range m.rainPos {
+				c.y += c.v
+				c.col = rainbowStyle(float64(i) / float64(len(m.rainPos)))
+				if c.y > m.height+5 {
+					c.y = -5
+				}
+			}
+		}
 		return m, tickerCmd()
 	}
 
@@ -123,14 +132,6 @@
 	case modeWave:
 		lines = append(lines, renderWave(m)...)
 	case modeRain:
-		// Update rain positions in-place so they persist across frames
-		for i := range m.rainPos {
-			m.rainPos[i].y += m.rainPos[i].v
-			m.rainPos[i].col = rainbowStyle(float64(i) / float64(len(m.rainPos)))
-			if m.rainPos[i].y > m.height+5 {
-				m.rainPos[i].y = -5
-			}
-		}
 		lines = append(lines, renderRainStatic(m)...)
 	case modeSpin:
 		lines = append(lines, renderSpin(m)...)
```
