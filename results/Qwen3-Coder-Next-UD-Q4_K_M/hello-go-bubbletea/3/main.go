package main

import (
	"fmt"
	"math"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width    int
	height   int
	frame    int
	spinner  spinner.Model
	quitting bool
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
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

func (m model) View() string {
	if m.quitting {
		return ""
	}

	// Create a rotating gradient effect
	hue := (m.frame * 5) % 360
	color1 := lipgloss.Color(fmt.Sprintf("#%06x", hueToRGB(hue)))
	color2 := lipgloss.Color(fmt.Sprintf("#%06x", hueToRGB((hue+120)%360)))

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(color1).
		Background(color2).
		Padding(1, 4).
		MarginTop(m.height/3-2).
		Render("Hello, World!")

	// Animated subtitle
	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Italic(true).
		MarginTop(1).
		Render("Welcome to BubbleTea!")

	// Spinner with interactive text
	spinnerView := m.spinner.View() + " Press q to quit"

	spinnerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		MarginTop(2).
		MarginLeft(m.width/2 - len(spinnerView)/2)

	// Border decoration
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("69")).
		Width(m.width - 4).
		Height(8).
		MarginTop(m.height/3-4).
		MarginLeft(2)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		border.Render(title+"\n"+subtitle+"\n\n"),
		spinnerStyle.Render(spinnerView),
	)
}

func hueToRGB(hue int) uint32 {
	// Convert hue to a nice visible color
	saturation := 0.8
	lightness := 0.5

	h := float64(hue) / 360
	c := (1 - math.Abs(2*lightness-1)) * saturation
	x := c * (1 - math.Abs(math.Mod(h*6, 2) - 1))
	m := lightness - c/2

	r := int((c + m) * 255)
	g := int((x + m) * 255)
	b := int((m) * 255)

	return uint32((r << 16) | (g << 8) | b)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
