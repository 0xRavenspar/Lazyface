package cmd

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const logoArt = `
__                                      ______                                     
|  \                                    /      \                                    
| $$       ______   ________  __    __ |  $$$$$$\ ______    _______   ______        
| $$      |      \ |        \|  \  |  \| $$_  \$$|      \  /       \ /      \       
| $$       \$$$$$$\ \$$$$$$$$| $$  | $$| $$ \     \$$$$$$\|  $$$$$$$|  $$$$$$\      
| $$      /      $$  /    $$ | $$  | $$| $$$$    /      $$| $$      | $$    $$      
| $$_____|  $$$$$$$ /  $$$$_ | $$__/ $$| $$     |  $$$$$$$| $$_____ | $$$$$$$$      
| $$     \\$$    $$|  $$    \ \$$    $$| $$      \$$    $$ \$$     \ \$$     \      
 \$$$$$$$$ \$$$$$$$ \$$$$$$$$ _\$$$$$$$ \$$       \$$$$$$$  \$$$$$$$  \$$$$$$$      
                             |  \__| $$                                             
                              \$$    $$                                             
                               \$$$$$$                                              
`

const frames = "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏⠸⠼⠴⠦⠧⠇⠏⠋⠙⠹"

type AnimateModel struct {
	frames     []string
	frameIndex int
	done       bool
	logoStyle  lipgloss.Style
}

type tickMsg time.Time
type animateFinishedMsg struct{}

func NewAnimateModel() AnimateModel {
	return AnimateModel{
		frames:    strings.Split(frames, ""),
		logoStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#f88e64")).Bold(true),
	}
}

func (m AnimateModel) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Millisecond*90, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
		tea.Sequence(
			tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
				return animateFinishedMsg{}
			}),
		),
	)
}

func (m AnimateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case tickMsg:
		if !m.done {
			m.frameIndex = (m.frameIndex + 1) % len(m.frames)
			return m, tea.Tick(time.Millisecond*90, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}
	case animateFinishedMsg:
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

func (m AnimateModel) View() string {
	if m.done {
		return ""
	}

	frame := m.frames[m.frameIndex]
	return "\n" + m.logoStyle.Render(logoArt) + "\n\n" +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#f88e64")).Render("Loading... ") +
		frame + "\n"
}
