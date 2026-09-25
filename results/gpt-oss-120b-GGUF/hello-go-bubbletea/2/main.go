package main

import (
    "fmt"
    "os"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// tickMsg is a custom message type for the spinner tick.
type tickMsg struct{}

// model holds the UI state.
type model struct {
    quit   bool
    frames []string
    idx    int
    // No need to store time; we just trigger ticks.
}

// Init is called when the program starts. It returns a command that sends the first tick.
func (m model) Init() tea.Cmd {
    // Send a tick after 120ms, then continue ticking.
    return tea.Tick(time.Millisecond*120, func(t time.Time) tea.Msg { return tickMsg{} })
}

func initialModel() model {
    frames := []string{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"}
    return model{frames: frames, idx: 0}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Requirement: always quit on 'q' or 'Q'.
        if msg.String() == "q" || msg.String() == "Q" {
            m.quit = true
            return m, tea.Quit
        }
    case tickMsg:
        // Advance spinner frame and schedule the next tick.
        m.idx = (m.idx + 1) % len(m.frames)
        return m, tea.Tick(time.Millisecond*120, func(t time.Time) tea.Msg { return tickMsg{} })
    }
    return m, nil
}

func (m model) View() string {
    if m.quit {
        return ""
    }
    spinner := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6")).Render(m.frames[m.idx])
    greeting := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8BE9FD")).Render("Hello, World!")
    hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render("q")
    return fmt.Sprintf("%s %s\n\nPress %s to quit.", spinner, greeting, hint)
}

func main() {
    p := tea.NewProgram(initialModel())
    if err := p.Start(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
