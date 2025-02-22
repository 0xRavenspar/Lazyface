package main

import (
	"Lazyface/cmd"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	mainContentStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).Width(80)
)

type model struct {
	views        []tea.Model
	activeView   int
	navigationUI cmd.NavigationModel
	footerUI     cmd.FooterModel
}

func initialModel() model {
	return model{
		views: []tea.Model{
			cmd.InitialModel(),       // Search
			cmd.InitialUploadModel(), // Upload
		},
		activeView:   0,
		navigationUI: cmd.NavigationModel{},
		footerUI:     cmd.FooterModel{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+l", "left":
			if m.activeView > 0 {
				m.activeView--
			}
		case "ctrl+r", "right":
			if m.activeView < len(m.views)-1 {
				m.activeView++
			}
		case "esc", "q":
			return m, tea.Quit
		}
	}

	updatedView, cmd := m.views[m.activeView].Update(msg)
	m.views[m.activeView] = updatedView
	m.navigationUI.ActiveView = m.activeView

	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.navigationUI.View(),
		mainContentStyle.Render(m.views[m.activeView].View()),
		m.footerUI.View(),
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen()) // Full-screen mode
	if err := p.Start(); err != nil {
		fmt.Println("Error starting the TUI:", err)
		os.Exit(1)
	}
}
