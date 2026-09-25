// hello.go
package main

import (
    "fmt"
    "os"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type model struct {
    tickCount int
    quitting  bool
}

type tickMsg time.Time

func (m model) Init() tea.Cmd {
    return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            m.quitting = true
            return m, tea.Quit
        }
    case tickMsg:
        m.tickCount++
        return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
            return tickMsg(t)
        })
    }
    return m, nil
}

var viewStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#0077CC")).
    Padding(1, 2).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#00FFAA"))

func (m model) View() string {
    if m.quitting {
        return ""
    }
    colors := []string{"#FF6B6B", "#F7B267", "#6EE7B7", "#4CC9F0", "#E0AFA0"}
    color := colors[m.tickCount%len(colors)]
    greeting := lipgloss.NewStyle().
        Foreground(lipgloss.Color(color)).
        Render("Hello, World!")
    return viewStyle.Render(fmt.Sprintf("%s\n\nPress q to quit.", greeting))
}

func main() {
    p := tea.NewProgram(model{})
    if err := p.Start(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
