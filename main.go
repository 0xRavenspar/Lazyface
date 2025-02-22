package main

import (
	"Lazyface/cmd"
	"Lazyface/internal/cli"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	mainContentStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).Width(80)
)

type model struct {
	views           []tea.Model
	activeView      int
	navigationUI    cmd.NavigationModel
	footerUI        cmd.FooterModel
	width           int
	height          int
	showSplash      bool
	isAuthenticated bool
	hasUserData     bool
}

func initialModel() model {
	var err error
	// Try initializing SettingsModel and check if user data exists
	hasUserData := err == nil // If no error, assume user data exists

	return model{
		views: []tea.Model{
			cmd.InitialSplashModel(), // Store as a pointer
		},
		activeView:   0,
		navigationUI: cmd.NavigationModel{},
		footerUI:     cmd.FooterModel{},
		showSplash:   true, // Start with splash screen
		hasUserData:  hasUserData,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y": // User chooses to log in
			if m.showSplash {
				m.isAuthenticated = true
				m.showSplash = false
				m.loadMainViews()
			}
		case "n": // User skips login
			if m.showSplash {
				m.isAuthenticated = false
				m.showSplash = false
				m.loadMainViews()
			}
		case "tab":
			if !m.showSplash {
				m.activeView = (m.activeView + 1) % len(m.views)
				m.navigationUI.ActiveView = m.activeView // Sync with navigation
			}
		case "esc", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		mainContentStyle = mainContentStyle.Width(m.width - 4)
		m.navigationUI.Width = m.width // Update navigation width dynamically
	}

	if !m.showSplash {
		updatedView, cmd := m.views[m.activeView].Update(msg)
		m.views[m.activeView] = updatedView
		m.navigationUI.ActiveView = m.activeView // Ensure navigation updates
		return m, cmd
	}

	return m, nil
}

func (m *model) loadMainViews() {
	var views []tea.Model
	var viewNames []string

	// Load user data
	_, err := cli.LoadUserData()
	m.hasUserData = err == nil

	if m.isAuthenticated {
		if m.hasUserData {
			// Authenticated & userData present -> No Auth View, include Settings
			views = append(views, cmd.InitialUploadModel(), cmd.InitialManageModel())

			settingsModel, _ := cmd.InitialSettingsModel()
			views = append(views, settingsModel)
			viewNames = append(viewNames, "Upload", "Manage", "Settings")
		} else {
			// Authenticated but no userData -> Show Auth View, no Settings
			views = append(views, cmd.NewAuthView(), cmd.InitialUploadModel(), cmd.InitialManageModel())
			viewNames = append(viewNames, "Auth", "Upload", "Manage")
		}
	} else {
		// Not authenticated -> Show only Download and Auth View
		views = append(views, cmd.InitialDownloadModel(), cmd.NewAuthView())
		viewNames = append(viewNames, "Download", "Auth")
	}

	m.views = views
	m.activeView = 0
	m.navigationUI.ViewNames = viewNames
}

func (m model) View() string {
	if m.showSplash {
		return mainContentStyle.Render(m.views[0].View()) // Render only splash screen
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.navigationUI.View(),
		mainContentStyle.Render(m.views[m.activeView].View()),
		m.footerUI.View(),
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error starting the TUI:", err)
		os.Exit(1)
	}
}
