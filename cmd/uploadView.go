package cmd

import (
	"Lazyface/internal/cli"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type uploadState int

const (
	inputUploadRepo uploadState = iota
	inputLocalPath
	inputRepoPath
	advancedOptions
	uploadConfirmation
	uploading
)

type uploadModel struct {
	state           uploadState
	repoInput       textinput.Model
	localPathInput  textinput.Model
	repoPathInput   textinput.Model
	repoTypeInput   textinput.Model
	includeInput    textinput.Model
	excludeInput    textinput.Model
	deleteInput     textinput.Model
	commitMsgInput  textinput.Model
	revisionInput   textinput.Model
	status          string
	cursorIndex     int
	advancedFields  []string
	advancedInputs  map[string]*textinput.Model
	showingAdvanced bool
}

func InitialUploadModel() uploadModel {
	repoInput := textinput.New()
	repoInput.Placeholder = "Enter Hugging Face repo name"
	repoInput.Focus()

	localPathInput := textinput.New()
	localPathInput.Placeholder = "Enter local path to upload"

	repoPathInput := textinput.New()
	repoPathInput.Placeholder = "Enter path in repo"

	repoTypeInput := textinput.New()
	repoTypeInput.Placeholder = "Enter repo type (optional)"

	includeInput := textinput.New()
	includeInput.Placeholder = "Include pattern (optional)"

	excludeInput := textinput.New()
	excludeInput.Placeholder = "Exclude pattern (optional)"

	deleteInput := textinput.New()
	deleteInput.Placeholder = "Delete pattern (optional)"

	commitMsgInput := textinput.New()
	commitMsgInput.Placeholder = "Commit message (optional)"

	revisionInput := textinput.New()
	revisionInput.Placeholder = "Revision (optional)"

	advancedFields := []string{
		"Repo Type",
		"Include Pattern",
		"Exclude Pattern",
		"Delete Pattern",
		"Commit Message",
		"Revision",
	}

	advancedInputs := map[string]*textinput.Model{
		"Repo Type":       &repoTypeInput,
		"Include Pattern": &includeInput,
		"Exclude Pattern": &excludeInput,
		"Delete Pattern":  &deleteInput,
		"Commit Message":  &commitMsgInput,
		"Revision":        &revisionInput,
	}

	return uploadModel{
		state:          inputUploadRepo,
		repoInput:      repoInput,
		localPathInput: localPathInput,
		repoPathInput:  repoPathInput,
		repoTypeInput:  repoTypeInput,
		includeInput:   includeInput,
		excludeInput:   excludeInput,
		deleteInput:    deleteInput,
		commitMsgInput: commitMsgInput,
		revisionInput:  revisionInput,
		status:         "Ready to upload.",
		advancedFields: advancedFields,
		advancedInputs: advancedInputs,
	}
}

func (m uploadModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m uploadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			switch m.state {
			case inputUploadRepo:
				if m.repoInput.Value() != "" {
					m.state = inputLocalPath
					m.localPathInput.Focus()
					m.repoInput.Blur()
				}
			case inputLocalPath:
				if m.localPathInput.Value() != "" {
					m.state = inputRepoPath
					m.repoPathInput.Focus()
					m.localPathInput.Blur()
				}
			case inputRepoPath:
				if m.repoPathInput.Value() != "" {
					m.state = advancedOptions
					m.repoPathInput.Blur()
				}
			case advancedOptions:
				if m.showingAdvanced {
					if m.cursorIndex < len(m.advancedFields) {
						field := m.advancedFields[m.cursorIndex]
						input := m.advancedInputs[field]
						input.Blur()
						m.cursorIndex++
						if m.cursorIndex < len(m.advancedFields) {
							m.advancedInputs[m.advancedFields[m.cursorIndex]].Focus()
						} else {
							m.state = uploadConfirmation
						}
					}
				} else {
					m.state = uploadConfirmation
				}
			case uploadConfirmation:
				m.state = uploading
				go func() {
					_, err := cli.UploadToHuggingFace(
						m.repoInput.Value(),
						m.localPathInput.Value(),
						m.repoPathInput.Value(),
						m.repoTypeInput.Value(),
						m.includeInput.Value(),
						m.excludeInput.Value(),
						m.deleteInput.Value(),
						m.commitMsgInput.Value(),
						m.revisionInput.Value(),
					)
					if err != nil {
						m.status = fmt.Sprintf("Error: %v", err)
					} else {
						m.status = "Upload complete!"
					}
				}()
			}

		case "tab":
			if m.state == advancedOptions {
				m.showingAdvanced = !m.showingAdvanced
				if m.showingAdvanced {
					m.cursorIndex = 0
					m.advancedInputs[m.advancedFields[0]].Focus()
				}
			}

		case "up":
			if m.state == advancedOptions && m.showingAdvanced {
				if m.cursorIndex > 0 {
					m.advancedInputs[m.advancedFields[m.cursorIndex]].Blur()
					m.cursorIndex--
					m.advancedInputs[m.advancedFields[m.cursorIndex]].Focus()
				}
			}

		case "down":
			if m.state == advancedOptions && m.showingAdvanced {
				if m.cursorIndex < len(m.advancedFields)-1 {
					m.advancedInputs[m.advancedFields[m.cursorIndex]].Blur()
					m.cursorIndex++
					m.advancedInputs[m.advancedFields[m.cursorIndex]].Focus()
				}
			}
		}
	}

	switch m.state {
	case inputUploadRepo:
		m.repoInput, cmd = m.repoInput.Update(msg)
	case inputLocalPath:
		m.localPathInput, cmd = m.localPathInput.Update(msg)
	case inputRepoPath:
		m.repoPathInput, cmd = m.repoPathInput.Update(msg)
	case advancedOptions:
		if m.showingAdvanced && m.cursorIndex < len(m.advancedFields) {
			field := m.advancedFields[m.cursorIndex]
			*m.advancedInputs[field], cmd = m.advancedInputs[field].Update(msg)
		}
	}

	return m, cmd
}

func (m uploadModel) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5f00")).Padding(1).Align(lipgloss.Center)
	bodyStyle := lipgloss.NewStyle().Padding(1, 2)

	switch m.state {
	case inputUploadRepo:
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Upload to Hugging Face"),
			bodyStyle.Render(fmt.Sprintf("%s\n\n%s\n\n%s",
				"Enter Repository Name:",
				m.repoInput.View(),
				"Press Enter to confirm, Q to quit")),
		)

	case inputLocalPath:
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Upload to Hugging Face"),
			bodyStyle.Render(fmt.Sprintf("%s\n\n%s\n\n%s",
				"Enter Local Path:",
				m.localPathInput.View(),
				"Press Enter to confirm, Q to quit")),
		)

	case inputRepoPath:
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Upload to Hugging Face"),
			bodyStyle.Render(fmt.Sprintf("%s\n\n%s\n\n%s",
				"Enter Path in Repository:",
				m.repoPathInput.View(),
				"Press Enter to confirm, Q to quit")),
		)

	case advancedOptions:
		var content string
		if m.showingAdvanced {
			var fields string
			for i, field := range m.advancedFields {
				cursor := " "
				if i == m.cursorIndex {
					cursor = ">"
				}
				fields += fmt.Sprintf("%s %s: %s\n", cursor, field, m.advancedInputs[field].View())
			}
			content = fmt.Sprintf("Advanced Options:\n\n%s\n\nUse ↑/↓ to navigate, Enter to confirm", fields)
		} else {
			content = "Press TAB to show/hide advanced options\nPress Enter to continue with default options"
		}
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Advanced Options"),
			bodyStyle.Render(content),
		)

	case uploadConfirmation:
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Confirmation"),
			bodyStyle.Render(fmt.Sprintf("Upload Details:\nRepository: %s\nLocal Path: %s\nRepo Path: %s\n\nPress Enter to start upload, Q to quit",
				m.repoInput.Value(),
				m.localPathInput.Value(),
				m.repoPathInput.Value())),
		)

	case uploading:
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Uploading..."),
			bodyStyle.Render(m.status),
		)
	}

	return ""
}
