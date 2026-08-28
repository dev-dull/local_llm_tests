package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type greetingMode int

const (
	modeIdle greetingMode = iota
	modeLoading
	modeDone
)

type tickMsg struct{}

type model struct {
	greeting   string
	mode       greetingMode
	spinner    spinner.Model
	tilt       float64
	bounce     float64
	width      int
	height     int
	colorIndex int
	ticker     *time.Ticker
	ctx        context.Context
	cancel     context.CancelFunc
}

func initialModel() *model {
	ctx, cancel := context.WithCancel(context.Background())
	return &model{
		greeting:   "Hello, World!",
		mode:       modeLoading,
		spinner:    spinner.New(spinner.WithSpinner(spinner.Points)),
		width:      80,
		height:     24,
		colorIndex: 0,
		ticker:     time.NewTicker(50 * time.Millisecond),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.animate(),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q", "ctrl+c":
			m.cancel()
			return m, tea.Quit
		case "space", "enter":
			if m.mode == modeLoading {
				m.mode = modeDone
			} else {
				m.mode = modeLoading
				return m, m.startLoading()
			}
		case "a", "A":
			// Toggle animation
			if m.ticker != nil {
				m.ticker.Stop()
				m.ticker = nil
			} else {
				m.ticker = time.NewTicker(50 * time.Millisecond)
			}
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		if m.ticker != nil {
			m.bounce = math.Sin(float64(time.Now().UnixMilli())/200.0) * 2
			m.tilt = math.Sin(float64(time.Now().UnixMilli())/300.0) * 0.3
			m.colorIndex++
			return m, m.animate()
		}
	}
	return m, nil
}

func (m *model) View() string {
	s := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("69")).
		Padding(1)

	switch m.mode {
	case modeLoading:
		return s.Render(m.renderLoading())
	default:
		return s.Render(m.renderGreeting())
	}
}

func (m *model) startLoading() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		m.mode = modeDone
		return tickMsg{}
	}
}

func (m *model) renderLoading() string {
	msg := lipgloss.NewStyle().
		Foreground(lipgloss.Color("38")).
		Align(lipgloss.Center).
		Render("Loading your greeting...")

	spinnerView := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Render(m.spinner.View())

	return fmt.Sprintf("\n\n%s %s", spinnerView, msg)
}

func (m *model) renderGreeting() string {
	// Create a colorful gradient effect
	colors := []lipgloss.Color{"196", "202", "208", "214", "220", "226", "190", "154", "118", "82", "46", "45", "44", "43", "42", "39", "63", "99", "135", "171"}

	// Build a colorful greeting with gradient
	var result string
	for i, char := range m.greeting {
		color := colors[(i+m.colorIndex)%len(colors)]
		result += lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true).
			Render(string(char))
	}

	// Add some decorative elements
	frameStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("69")).
		Align(lipgloss.Center)

	bottomText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center).
		Render("Press 'q' to quit | Space to reload | 'a' to toggle animation")

	return fmt.Sprintf("%s\n\n%s\n\n%s", frameStyle.Render("✨"), result, bottomText)
}

func (m *model) animate() tea.Cmd {
	return func() tea.Msg {
		select {
		case <-m.ctx.Done():
			return nil
		case <-m.ticker.C:
			return tickMsg{}
		}
	}
}

func main() {
	if _, err := tea.NewProgram(
		initialModel(),
		tea.WithMouseCellMotion(),
	).Run(); err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		os.Exit(1)
	}
}
