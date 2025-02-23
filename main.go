package main

import (
	"Lazyface/cmd"
	"Lazyface/internal/cli"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	mainContentStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).Width(80)
)

type tickMsg time.Time

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
	animationModel  cmd.AnimateModel
	showAnimation   bool
}

func initialModel() model {
	var err error
	_, err = cli.LoadUserData()
	hasUserData := err == nil // Check if user data exists

	m := model{
		navigationUI:    cmd.NavigationModel{},
		footerUI:        cmd.FooterModel{},
		hasUserData:     hasUserData,
		isAuthenticated: hasUserData, // Assume authenticated if user data exists
		animationModel:  cmd.NewAnimateModel(),
	}

	if hasUserData {
		// Skip splash and animation, go directly to the main views
		m.showSplash = false
		m.showAnimation = false
		m.loadMainViews()
	} else {
		// Show splash and animation
		m.showSplash = true
		m.showAnimation = true
		m.views = []tea.Model{cmd.InitialSplashModel()}
	}

	return m
}

func (m model) Init() tea.Cmd {
	if m.showAnimation {
		// Show the animation for 3 seconds (adjust as needed)
		return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
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
	case tickMsg:
		if m.showAnimation {
			m.showAnimation = false
			m.showSplash = true
		}
	}

	if !m.showSplash {
		updatedView, cmd := m.views[m.activeView].Update(msg)
		m.views[m.activeView] = updatedView
		m.navigationUI.ActiveView = m.activeView
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
			views = append(views, cmd.InitialUploadModel(), cmd.InitialManageModel())

			settingsModel, _ := cmd.InitialSettingsModel()
			views = append(views, settingsModel)
			viewNames = append(viewNames, "Upload", "Manage", "Settings")
		} else {
			views = append(views, cmd.NewAuthView(), cmd.InitialUploadModel(), cmd.InitialManageModel())
			viewNames = append(viewNames, "Auth", "Upload", "Manage")
		}
	} else {
		views = append(views, cmd.InitialDownloadModel(), cmd.NewAuthView())
		viewNames = append(viewNames, "Download", "Auth")
	}

	m.views = views
	m.activeView = 0
	m.navigationUI.ViewNames = viewNames
}

func (m model) View() string {
	if m.showAnimation {
		return mainContentStyle.Render(m.animationModel.View())
	}
	if m.showSplash {
		return mainContentStyle.Render(m.views[0].View())
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
