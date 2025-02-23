package cmd

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	tabStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 1)
	activeTabStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170")).Underline(true).Padding(0, 1)

	navigationStyle = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Padding(1, 2)
)

type NavigationModel struct {
	ActiveView int
	Width      int
	ViewNames  []string // Dynamically update based on available views
}

func (m NavigationModel) RenderTabs() string {
	var tabs []string
	for i, name := range m.ViewNames {
		if i == m.ActiveView {
			tabs = append(tabs, activeTabStyle.Render(name))
		} else {
			tabs = append(tabs, tabStyle.Render(name))
		}
	}

	navigationStyle = navigationStyle.Width(m.Width - 4) // Adjust width dynamically
	return navigationStyle.Render(strings.Join(tabs, " | "))
}

func (m NavigationModel) View() string {
	return m.RenderTabs()
}
