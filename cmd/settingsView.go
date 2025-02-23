package cmd

import (
	"Lazyface/internal/cli"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingsModel struct {
	username      string
	fullName      string
	organizations []string
	tokenName     string
	permissions   []string
	logoutBtn     bool
	cursor        int
	pageSize      int
}

func InitialSettingsModel() (SettingsModel, error) {
	userData, err := cli.LoadUserData()
	if err != nil {
		return SettingsModel{}, fmt.Errorf("failed to load user data: %w", err)
	}

	return SettingsModel{
		username:      userData.Name,
		fullName:      userData.FullName,
		organizations: userData.Orgs,
		tokenName:     userData.TokenName,
		permissions:   userData.Permissions,
		logoutBtn:     false,
		cursor:        0,
		pageSize:      10,
	}, nil
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "f":
			fmt.Println("Logging out...")
			if err := cli.Logout(); err != nil {
				fmt.Println("Error logging out:", err)
			}
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.permissions)-m.pageSize {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m SettingsModel) View() string {
	styleHeader := lipgloss.NewStyle().Bold(true)
	styleText := lipgloss.NewStyle().Foreground(lipgloss.Color("#f88e64"))
	permText := lipgloss.NewStyle().Foreground(lipgloss.Color("#fe6375"))
	btnText := lipgloss.NewStyle().Background(lipgloss.Color("#fe6375"))
	styleDim := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	userInfo := "\n" + styleHeader.Width(50).Render("User Information") + "\n" +
		styleText.Render("Username: ") + m.username + "\n" +
		styleText.Render("Full Name: ") + m.fullName + "\n" +
		styleText.Render("Organizations: ") + fmt.Sprintf("%v", m.organizations) + "\n" +
		styleText.Render("Token Name: ") + m.tokenName + "\n"

	logoutBtn := "\n" + btnText.Render("Press f to logout")

	first := lipgloss.JoinHorizontal(lipgloss.Top, userInfo, logoutBtn)

	permissionsHeader := styleHeader.Render("Token Permissions") + "\n"
	var permissionsList string

	start := m.cursor
	end := start + m.pageSize
	if end > len(m.permissions) {
		end = len(m.permissions)
	}

	for _, perm := range m.permissions[start:end] {
		permissionsList += permText.Render(perm) + "\n"
	}

	totalPages := (len(m.permissions) + m.pageSize - 1) / m.pageSize
	currentPage := (m.cursor / m.pageSize) + 1
	pagination := styleDim.Render(fmt.Sprintf("Page %d/%d (Use ↑/↓ to navigate through permissions list)", currentPage, totalPages))

	return lipgloss.JoinVertical(lipgloss.Left,
		first,
		permissionsHeader,
		permissionsList,
		pagination,
	)
}
