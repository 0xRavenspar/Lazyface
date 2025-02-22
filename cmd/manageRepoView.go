package cmd

import (
	"Lazyface/internal/cli"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type manageState int

const (
	selectOperation manageState = iota
	createRepoState
	deleteRepoState
	updateVisibilityState
	moveRepoState
	confirmOperation
	processingOperation
)

type manageModel struct {
	state          manageState
	selectedOp     string
	cursorPosition int
	tokenInput     textinput.Model
	repoTypeInput  textinput.Model
	repoNameInput  textinput.Model
	orgInput       textinput.Model
	isPrivate      bool
	sdkInput       textinput.Model
	fromRepoInput  textinput.Model
	toRepoInput    textinput.Model
	currentField   int
	status         string
	error          string
}

var operations = []string{
	"Create Repository",
	"Delete Repository",
	"Update Repository Visibility",
	"Move Repository",
}

func (m *manageModel) blurAllInputs() {
	m.tokenInput.Blur()
	m.repoTypeInput.Blur()
	m.repoNameInput.Blur()
	m.orgInput.Blur()
	m.sdkInput.Blur()
	m.fromRepoInput.Blur()
	m.toRepoInput.Blur()
}

func InitialManageModel() manageModel {
	tokenInput := textinput.New()
	tokenInput.Placeholder = "Enter Hugging Face token"

	repoTypeInput := textinput.New()
	repoTypeInput.Placeholder = "Enter repo type (model/dataset/space)"

	repoNameInput := textinput.New()
	repoNameInput.Placeholder = "Enter repository name"

	orgInput := textinput.New()
	orgInput.Placeholder = "Enter organization (optional)"

	sdkInput := textinput.New()
	sdkInput.Placeholder = "Enter SDK"

	fromRepoInput := textinput.New()
	fromRepoInput.Placeholder = "Enter source repository"

	toRepoInput := textinput.New()
	toRepoInput.Placeholder = "Enter destination repository"

	return manageModel{
		state:         selectOperation,
		tokenInput:    tokenInput,
		repoTypeInput: repoTypeInput,
		repoNameInput: repoNameInput,
		orgInput:      orgInput,
		sdkInput:      sdkInput,
		fromRepoInput: fromRepoInput,
		toRepoInput:   toRepoInput,
		isPrivate:     false,
	}
}

func (m manageModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m manageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle space key press first, before text input processing
		if msg.String() == " " {
			switch m.state {
			case createRepoState, updateVisibilityState:
				m.isPrivate = !m.isPrivate
				return m, cmd
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			switch m.state {
			case selectOperation:
				if m.cursorPosition > 0 {
					m.cursorPosition--
				}
			case createRepoState:
				if m.currentField > 0 {
					m.currentField--
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.repoNameInput.Focus()
					case 3:
						m.orgInput.Focus()
					case 4:
						m.sdkInput.Focus()
					}
				}
			case deleteRepoState:
				if m.currentField > 0 {
					m.currentField--
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.repoNameInput.Focus()
					case 3:
						m.orgInput.Focus()
					}
				}
			case updateVisibilityState:
				if m.currentField > 0 {
					m.currentField--
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.repoNameInput.Focus()
					}
				}
			case moveRepoState:
				if m.currentField > 0 {
					m.currentField--
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.fromRepoInput.Focus()
					case 3:
						m.toRepoInput.Focus()
					}
				}
			}

		case "down", "j":
			switch m.state {
			case selectOperation:
				if m.cursorPosition < len(operations)-1 {
					m.cursorPosition++
				}
			case createRepoState:
				if m.currentField < 4 {
					m.currentField++
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.repoNameInput.Focus()
					case 3:
						m.orgInput.Focus()
					case 4:
						m.sdkInput.Focus()
					}
				}
			case deleteRepoState:
				if m.currentField < 3 {
					m.currentField++
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.repoNameInput.Focus()
					case 3:
						m.orgInput.Focus()
					}
				}
			case updateVisibilityState:
				if m.currentField < 2 {
					m.currentField++
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.repoNameInput.Focus()
					}
				}
			case moveRepoState:
				if m.currentField < 3 {
					m.currentField++
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.fromRepoInput.Focus()
					case 3:
						m.toRepoInput.Focus()
					}
				}
			}

		case "enter":
			switch m.state {
			case selectOperation:
				m.selectedOp = operations[m.cursorPosition]
				m.currentField = 0
				m.blurAllInputs()
				m.tokenInput.Focus()

				switch operations[m.cursorPosition] {
				case "Create Repository":
					m.state = createRepoState
				case "Delete Repository":
					m.state = deleteRepoState
				case "Update Repository Visibility":
					m.state = updateVisibilityState
				case "Move Repository":
					m.state = moveRepoState
				}
				return m, cmd

			case createRepoState:
				if m.currentField < 4 {
					m.currentField++
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.repoNameInput.Focus()
					case 3:
						m.orgInput.Focus()
					case 4:
						m.sdkInput.Focus()
					}
				} else {
					m.state = confirmOperation
				}

			case deleteRepoState:
				if m.currentField < 3 {
					m.currentField++
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.repoNameInput.Focus()
					case 3:
						m.orgInput.Focus()
					}
				} else {
					m.state = confirmOperation
				}

			case updateVisibilityState:
				if m.currentField < 2 {
					m.currentField++
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.repoNameInput.Focus()
					}
				} else {
					m.state = confirmOperation
				}

			case moveRepoState:
				if m.currentField < 3 {
					m.currentField++
					m.blurAllInputs()
					switch m.currentField {
					case 0:
						m.tokenInput.Focus()
					case 1:
						m.repoTypeInput.Focus()
					case 2:
						m.fromRepoInput.Focus()
					case 3:
						m.toRepoInput.Focus()
					}
				} else {
					m.state = confirmOperation
				}

			case confirmOperation:
				m.state = processingOperation
				go func() {
					var err error
					switch m.selectedOp {
					case "Create Repository":
						err = cli.CreateRepo(
							m.tokenInput.Value(),
							m.repoTypeInput.Value(),
							m.repoNameInput.Value(),
							m.orgInput.Value(),
							m.isPrivate,
							m.sdkInput.Value(),
						)
					case "Delete Repository":
						err = cli.DeleteRepo(
							m.tokenInput.Value(),
							m.repoTypeInput.Value(),
							m.repoNameInput.Value(),
							m.orgInput.Value(),
						)
					case "Update Repository Visibility":
						err = cli.UpdateRepoVisibility(
							m.tokenInput.Value(),
							m.repoTypeInput.Value(),
							m.repoNameInput.Value(),
							m.isPrivate,
						)
					case "Move Repository":
						err = cli.MoveRepo(
							m.tokenInput.Value(),
							m.fromRepoInput.Value(),
							m.toRepoInput.Value(),
							m.repoTypeInput.Value(),
						)
					}
					if err != nil {
						m.error = err.Error()
					} else {
						m.status = "Operation completed successfully!"
					}
				}()
			}
		}
	}

	// Handle text input updates
	if msg, ok := msg.(tea.KeyMsg); ok {
		// Skip text input processing if we just handled a space key
		if msg.String() == " " {
			return m, cmd
		}
	}

	switch m.state {
	case createRepoState, deleteRepoState, updateVisibilityState, moveRepoState:
		var tmpCmd tea.Cmd
		switch m.currentField {
		case 0:
			m.tokenInput, tmpCmd = m.tokenInput.Update(msg)
			cmd = tmpCmd
		case 1:
			m.repoTypeInput, tmpCmd = m.repoTypeInput.Update(msg)
			cmd = tmpCmd
		case 2:
			switch m.state {
			case createRepoState, deleteRepoState, updateVisibilityState:
				m.repoNameInput, tmpCmd = m.repoNameInput.Update(msg)
			case moveRepoState:
				m.fromRepoInput, tmpCmd = m.fromRepoInput.Update(msg)
			}
			cmd = tmpCmd
		case 3:
			switch m.state {
			case createRepoState, deleteRepoState:
				m.orgInput, tmpCmd = m.orgInput.Update(msg)
			case moveRepoState:
				m.toRepoInput, tmpCmd = m.toRepoInput.Update(msg)
			}
			cmd = tmpCmd
		case 4:
			if m.state == createRepoState {
				m.sdkInput, tmpCmd = m.sdkInput.Update(msg)
				cmd = tmpCmd
			}
		}
	}

	return m, cmd
}

func (m manageModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5f00")).Padding(1).Align(lipgloss.Center)
	bodyStyle := lipgloss.NewStyle().Padding(1, 2)

	switch m.state {
	case selectOperation:
		var s string
		s = "Select an operation:\n\n"

		for i, op := range operations {
			cursor := " "
			if m.cursorPosition == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, op)
		}

		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Manage Hugging Face Repositories"),
			bodyStyle.Render(s+"\nUse ↑/↓ to select, Enter to confirm, Q to quit"),
		)

	case createRepoState:
		fields := []string{
			fmt.Sprintf("HF Token: %s", m.tokenInput.View()),
			fmt.Sprintf("Repo Type (model/dataset/space): %s", m.repoTypeInput.View()),
			fmt.Sprintf("Repo Name: %s", m.repoNameInput.View()),
			fmt.Sprintf("Organization (optional): %s", m.orgInput.View()),
			fmt.Sprintf("SDK: %s", m.sdkInput.View()),
			fmt.Sprintf("Private: %v (space to toggle)", m.isPrivate),
		}

		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Create Repository"),
			bodyStyle.Render(fmt.Sprintf("%s\n\nPress Enter to continue, Q to quit",
				lipgloss.JoinVertical(lipgloss.Left, fields...))),
		)

	case deleteRepoState:
		fields := []string{
			fmt.Sprintf("HF Token: %s", m.tokenInput.View()),
			fmt.Sprintf("Repo Type (model/dataset/space): %s", m.repoTypeInput.View()),
			fmt.Sprintf("Repo Name: %s", m.repoNameInput.View()),
			fmt.Sprintf("Organization (optional): %s", m.orgInput.View()),
		}

		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Delete Repository"),
			bodyStyle.Render(fmt.Sprintf("%s\n\nPress Enter to continue, Q to quit",
				lipgloss.JoinVertical(lipgloss.Left, fields...))),
		)

	case updateVisibilityState:
		fields := []string{
			fmt.Sprintf("HF Token: %s", m.tokenInput.View()),
			fmt.Sprintf("Repo Type (model/dataset/space): %s", m.repoTypeInput.View()),
			fmt.Sprintf("Repo Name: %s", m.repoNameInput.View()),
			fmt.Sprintf("Private: %v (space to toggle)", m.isPrivate),
		}

		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Update Repository Visibility"),
			bodyStyle.Render(fmt.Sprintf("%s\n\nPress Enter to continue, Q to quit",
				lipgloss.JoinVertical(lipgloss.Left, fields...))),
		)

	case moveRepoState:
		fields := []string{
			fmt.Sprintf("HF Token: %s", m.tokenInput.View()),
			fmt.Sprintf("Repo Type (model/dataset/space): %s", m.repoTypeInput.View()),
			fmt.Sprintf("From Repo: %s", m.fromRepoInput.View()),
			fmt.Sprintf("To Repo: %s", m.toRepoInput.View()),
		}

		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Move Repository"),
			bodyStyle.Render(fmt.Sprintf("%s\n\nPress Enter to continue, Q to quit",
				lipgloss.JoinVertical(lipgloss.Left, fields...))),
		)

	case confirmOperation:
		var details string
		switch m.selectedOp {
		case "Create Repository":
			details = fmt.Sprintf("Create new %s repository: %s\nOrganization: %s\nPrivate: %v\nSDK: %s",
				m.repoTypeInput.Value(), m.repoNameInput.Value(), m.orgInput.Value(), m.isPrivate, m.sdkInput.Value())
		case "Delete Repository":
			details = fmt.Sprintf("Delete %s repository: %s\nOrganization: %s",
				m.repoTypeInput.Value(), m.repoNameInput.Value(), m.orgInput.Value())
		case "Update Repository Visibility":
			details = fmt.Sprintf("Update visibility of %s repository: %s\nNew visibility: %v",
				m.repoTypeInput.Value(), m.repoNameInput.Value(), m.isPrivate)
		case "Move Repository":
			details = fmt.Sprintf("Move %s repository\nFrom: %s\nTo: %s",
				m.repoTypeInput.Value(), m.fromRepoInput.Value(), m.toRepoInput.Value())
		}

		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Confirm Operation"),
			bodyStyle.Render(fmt.Sprintf("%s\n\nPress Enter to execute, Q to quit", details)),
		)

	case processingOperation:
		status := m.status
		if m.error != "" {
			status = fmt.Sprintf("Error: %s", m.error)
		}
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Processing"),
			bodyStyle.Render(status),
		)
	}

	return ""
}
