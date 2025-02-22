package cmd

import (
	"github.com/charmbracelet/lipgloss"
)

var footerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("245")).
	Padding(1, 0).
	Width(80).
	Align(lipgloss.Center)

type FooterModel struct{}

func (f FooterModel) View() string {
	return footerStyle.Render(
		"[Tab] Switch View | [Q / Esc] Quit",
	)
}
