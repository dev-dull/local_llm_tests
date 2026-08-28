package main

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	rotation  float64
	frame     int
	spinner   spinner.Model
	quitting  bool
	width     int
	height    int
}

func initialModel() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return m.Tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "Q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		return model{
			width:   msg.Width,
			height:  msg.Height,
			spinner: m.spinner,
		}, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m model) View() string {
	if m.quitting {
		return "\n  " + lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true).
			Render("Goodbye from Bubbletea! 🌈\n")
	}

	// Create a rotating animation effect
	rotatedX := int(math.Sin(m.rotation) * 10)
	rotatedY := int(math.Cos(m.rotation) * 5)

	welcomeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("69")).
		Background(lipgloss.Color("235")).
		Bold(true).
		Padding(1, 3).
		MarginTop(max(0, rotatedY+2)).
		MarginLeft(max(0, rotatedX))

	helloText := "Hello, World! 🚀"
	greeting := welcomeStyle.Render(helloText)

	// Spinner section
	spinnerView := m.spinner.View()
	spinnerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		MarginTop(1).
		AlignHorizontal(lipgloss.Center)

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		AlignHorizontal(lipgloss.Center).
		MarginTop(2).
		Render("Press 'q' to quit")

	// Rainbow gradient bar
	rainbow := lipgloss.NewStyle().
		Height(1).
		MarginTop(1).
		Render(lipgloss.NewStyle().
			Inline(true).
			Background(lipgloss.Color("1")).
			Render(" ") + lipgloss.NewStyle().Background(lipgloss.Color("3")).Render(" ") +
			lipgloss.NewStyle().Background(lipgloss.Color("4")).Render(" ") +
			lipgloss.NewStyle().Background(lipgloss.Color("5")).Render(" ") +
			lipgloss.NewStyle().Background(lipgloss.Color("6")).Render(" ") +
			lipgloss.NewStyle().Background(lipgloss.Color("201")).Render(" "))

	return "\n" + rainbow + "\n" + spinnerStyle.Render(spinnerView) + "\n" + greeting + "\n" + instructions + "\n"
}

func (m model) Tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		m.rotation += 0.1
		m.frame++
		return spinner.TickMsg{}
	})
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}
