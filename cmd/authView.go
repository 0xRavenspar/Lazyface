// auth_view.go
package cmd

import (
	"Lazyface/internal/cli"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AuthView struct {
	width            int
	height           int
	selected         int
	token            string
	focusIndex       int
	errorMsg         string
	isAuthenticating bool
}

// Message types
type loginSuccessMsg struct{}
type loginErrMsg struct{ err error }
type pasteMsg struct{ text string }

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF69B4"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700"))

	unselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))
)

// pasteFromClipboard is a command that reads from the clipboard
func pasteFromClipboard() tea.Msg {
	str, err := clipboard.ReadAll()
	if err != nil {
		return loginErrMsg{err}
	}
	return pasteMsg{text: str}
}

func NewAuthView() AuthView {
	return AuthView{
		selected:   0,
		width:      80,
		height:     24,
		token:      "",
		focusIndex: 0,
	}
}

func (a AuthView) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func performLogin(token string, addGitCredential bool) tea.Cmd {
	return func() tea.Msg {
		err := cli.Login(token, addGitCredential)
		if err != nil {
			return loginErrMsg{err}
		}
		return loginSuccessMsg{}
	}
}

func (a AuthView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case loginSuccessMsg:
		a.isAuthenticating = false
		return a, nil

	case loginErrMsg:
		a.isAuthenticating = false
		a.errorMsg = msg.err.Error()
		return a, nil

	case pasteMsg:
		if a.focusIndex == 1 {
			a.token = msg.text
			a.errorMsg = ""
		}
		return a, nil

	case tea.KeyMsg:
		if a.isAuthenticating {
			return a, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "tab":
			a.focusIndex = (a.focusIndex + 1) % 2
		case "shift+tab":
			a.focusIndex = (a.focusIndex - 1)
			if a.focusIndex < 0 {
				a.focusIndex = 1
			}
		case "up", "k":
			if a.focusIndex == 0 {
				if a.selected > 0 {
					a.selected--
				}
			}
		case "down", "j":
			if a.focusIndex == 0 {
				if a.selected < 1 {
					a.selected++
				}
			}
		case "ctrl+v", "ctrl+shift+v":
			if a.focusIndex == 1 {
				return a, tea.Cmd(pasteFromClipboard)
			}
		case "backspace":
			if a.focusIndex == 1 && len(a.token) > 0 {
				a.token = a.token[:len(a.token)-1]
				a.errorMsg = ""
			}
		case "enter":
			if len(a.token) > 0 {
				a.isAuthenticating = true
				return a, performLogin(a.token, a.selected == 0)
			} else if a.focusIndex == 0 && a.selected == 0 {
				a.focusIndex = 1
			}
		case "esc":
			a.focusIndex = 0
			a.errorMsg = ""
		default:
			if a.focusIndex == 1 {
				if len(msg.String()) == 1 {
					a.token += msg.String()
					a.errorMsg = ""
				}
			}
		}
	}
	return a, nil
}

func (a AuthView) View() string {
	if a.width == 0 {
		return "Initializing..."
	}

	var s strings.Builder

	// Navigation help
	helpText := "TAB cycle fields    ↑/k move up      ENTER select    CTRL+V paste"
	helpText += "\nESC back to menu   ↓/j move down    q quit"
	help := helpStyle.Render(helpText)

	// Title with emoji
	title := titleStyle.Render("😊 LazyFace")

	// Main content
	content := "\n\nTo log in, 'huggingface_hub' requires a token\ngenerated from:\n"
	content += linkStyle.Render("https://huggingface.co/settings/tokens")
	content += "\n\nWould you like to add token as git credential??\n\n"

	// Options
	yes := "Yes"
	no := "No"
	if a.selected == 0 {
		yes = selectedStyle.Render("> " + yes)
		no = unselectedStyle.Render("  " + no)
	} else {
		yes = unselectedStyle.Render("  " + yes)
		no = selectedStyle.Render("> " + no)
	}

	content += yes + "\n" + no + "\n\n"

	// Token input
	content += "Enter your token\n"
	if a.focusIndex == 1 {
		content += selectedStyle.Render("> " + a.token + "█")
	} else {
		content += unselectedStyle.Render("> " + a.token)
	}

	// Error message
	if a.errorMsg != "" {
		content += "\n\n" + errorStyle.Render(a.errorMsg)
	}

	// Authentication status
	if a.isAuthenticating {
		content += "\n\nAuthenticating..."
	}

	// Calculate available space

	// Combine all elements
	mainContent := fmt.Sprintf("%s\n%s\n%s", help, title, content)

	s.WriteString(mainContent)

	// Footer with safe padding calculation
	footer := "Esc/Back Home"
	s.WriteString("\n" + footer)

	return s.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
