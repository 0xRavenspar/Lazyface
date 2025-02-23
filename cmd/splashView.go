package cmd

import (
	"fmt"
	"image/color"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

type Splashmodel struct {
	cursor   int
	selected bool
}

func InitialSplashModel() *Splashmodel { // Return a pointer to Splashmodel
	return &Splashmodel{cursor: 0, selected: false} // Return an address
}

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}

func (m *Splashmodel) Init() tea.Cmd {
	return nil
}

func (m *Splashmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 1 { // This is correct since we have 2 choices (0 and 1)
				m.cursor++
			}
		case "enter":
			m.selected = true
		}
	}
	return m, nil
}

func (m *Splashmodel) View() string {
	// Gradient Header
	blends := gamut.Blends(lipgloss.Color("#FD5392"), lipgloss.Color("#F86F64"), 50)
	titleStyle := lipgloss.NewStyle().Align(lipgloss.Center).Bold(true)
	firstln := titleStyle.Render(rainbow(lipgloss.NewStyle(), "Easily browse, download, and manage HuggingFace", blends))
	secondln := titleStyle.Render(rainbow(lipgloss.NewStyle(), "models & datasets right from your terminal.", blends))

	// Column styles
	columnStyle := lipgloss.NewStyle().Width(45)
	check := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Render("✔")
	xmark := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4D4D")).Render("✖")

	textStyleCheck := lipgloss.NewStyle().Foreground(lipgloss.Color("#F88E64"))
	textStyleXMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#FE6375"))

	leftColumn := []string{
		fmt.Sprintf("%s %s", check, textStyleCheck.Render("Browse public models & datasets")),
		fmt.Sprintf("%s %s", check, textStyleCheck.Render("Download public models")),
		fmt.Sprintf("%s %s", check, textStyleCheck.Render("Run inference on local models")),
	}

	rightColumn := []string{
		fmt.Sprintf("%s %s", xmark, textStyleXMark.Render("Access private or gated models/datasets")),
		fmt.Sprintf("%s %s", xmark, textStyleXMark.Render("Upload models/datasets to the Hub")),
		fmt.Sprintf("%s %s", xmark, textStyleXMark.Render("Use Inference API without rate limits")),
		fmt.Sprintf("%s %s", xmark, textStyleXMark.Render("Manage organizations & settings")),
	}

	var leftView, rightView string
	for i := 0; i < len(leftColumn) || i < len(rightColumn); i++ {
		if i < len(leftColumn) {
			leftView += leftColumn[i] + "\n"
		} else {
			leftView += "\n"
		}

		if i < len(rightColumn) {
			rightView += rightColumn[i] + "\n"
		} else {
			rightView += "\n"
		}
	}

	columns := lipgloss.JoinHorizontal(lipgloss.Center, columnStyle.Render(leftView), columnStyle.Render(rightView))

	// Yes/No selection styling
	choices := []string{"Yes", "No"}
	selection := "Would you like to login?\n\n"
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)

	for i, choice := range choices {
		cursor := "  " // Default spacing
		displayChoice := lipgloss.NewStyle().Render(choice)

		if m.cursor == i {
			cursor = "👉"                                 // Visually indicate selection
			displayChoice = selectedStyle.Render(choice) // Highlight the selected choice
		}

		selection += fmt.Sprintf("%s %s\n", cursor, displayChoice)
	}

	return fmt.Sprintf("\n%s\n%s\n\n%s\n\n%s\n\n%s", firstln, secondln, lipgloss.NewStyle().Render("Some Actions Require Login"), columns, selection)
}
