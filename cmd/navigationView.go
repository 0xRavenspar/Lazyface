package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	tabStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 1)
	activeTabStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170")).Underline(true).Padding(0, 1)

	navigationStyle = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Padding(1, 2).Width(80).Align(lipgloss.Center)
)

type NavigationModel struct {
	ActiveView int
}

var viewNames = []string{"Search", "Download", "Upload", "Auth", "Settings"}

func (m NavigationModel) RenderTabs() string {
	var tabs []string
	for i, name := range viewNames {
		if i == m.ActiveView {
			tabs = append(tabs, activeTabStyle.Render(name))
		} else {
			tabs = append(tabs, tabStyle.Render(name))
		}
	}
	return navigationStyle.Render(strings.Join(tabs, " | "))
}

func (m NavigationModel) View() string {
	return fmt.Sprintf("%s", m.RenderTabs())
}
