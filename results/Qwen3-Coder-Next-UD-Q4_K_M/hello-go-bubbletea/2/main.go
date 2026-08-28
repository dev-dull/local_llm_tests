package main

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width         int
	height        int
	tick          int
	spinner       spinner.Model
	isProcessing  bool
	showWelcome   bool
	quoteIndex    int
	quotes        []string
}

var quotes = []string{
	"Hello, World!",
	"Programming is fun!",
	"Bubbletea makes it beautiful!",
	"Go is awesome!",
	"Interactive apps rock!",
}

func initialModel() Model {
	s := spinner.New(spinner.WithSpinner(spinner.Line))
	return Model{
		tick:         0,
		spinner:      s,
		quotes:       quotes,
		showWelcome:  false,
		quoteIndex:   0,
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q", "esc", "ctrl+c":
			return m, tea.Quit
		case " ", "enter":
			if !m.isProcessing {
				return m, m.startProcessing()
			}
		case "right", "l":
			m.quoteIndex = (m.quoteIndex + 1) % len(m.quotes)
		case "left", "h":
			m.quoteIndex = (m.quoteIndex - 1 + len(m.quotes)) % len(m.quotes)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) startProcessing() tea.Cmd {
	m.isProcessing = true
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		m.isProcessing = false
		m.showWelcome = true
		return tea.Tick(0)
	}
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	s := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(1, 2)

	// Animated background
	bg := m.animatedBackground()

	// Title with rainbow effect
	title := m.rainbowTitle("Hello, World!")

	// Quote display
	quoteStyle := lipgloss.NewStyle().
		FontSize(16).
		Bold(true).
		Align(lipgloss.Center)

	// Status indicator
	var status string
	if m.isProcessing {
		status = fmt.Sprintf(" %s Processing...", m.spinner.View())
	} else if m.showWelcome {
		status = " Press q to quit | Use left/right to change quote"
	} else {
		status = " Press space/enter to start | Use left/right to change quote"
	}

	statusStyle := lipgloss.NewStyle().
		FontSize(12).
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		title,
		"",
		quoteStyle.Render(m.quotes[m.quoteIndex]),
		"",
		lipgloss.NewStyle().Height(3).Render(status),
	)

	return s.Render(bg + "\n" + content)
}

func (m Model) animatedBackground() string {
	lines := ""
	rows := m.height - 8
	cols := m.width - 4

	for i := 0; i < rows; i++ {
		line := ""
		for j := 0; j < cols; j++ {
			char := "."
			// Create a subtle moving pattern
			t := float64(m.tick) / 20.0
			x := float64(j)
			y := float64(i)
			val := math.Sin(x/10+t) * math.Cos(y/10+t)
			if val > 0.7 {
				char = "o"
			} else if val > 0.3 {
				char = "*"
			} else if val > -0.3 {
				char = "."
			}
			line += char
		}
		lines += line + "\n"
	}
	return lines
}

func (m Model) rainbowTitle(text string) string {
	style := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center)

	colors := []string{"196", "202", "208", "214", "220", "226", "118", "82", "39", "99"}

	result := ""
	for i, r := range text {
		color := colors[i%len(colors)]
		result += lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Render(string(r))
	}

	return style.Render(result)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
