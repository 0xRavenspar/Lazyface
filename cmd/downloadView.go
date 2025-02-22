package cmd

import (
	"Lazyface/internal/cli"
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	inputRepo state = iota
	selectFiles
	selectDownloadPath
	confirmation
	downloading
)

const (
	filesPerColumn    = 10 //Number of files per column
	maxVisibleColumns = 3
)

type statusMsg string
type progressMsg float64

// Define the model structure
type model struct {
	state              state
	repoInput          textinput.Model
	files              []string
	selectedFiles      []string
	checked            map[int]bool
	cursorIndex        int
	scrollOffset       int
	downloadPathChoice string
	path               string
	customPathInput    textinput.Model
	status             string
	statusChan         chan string
	progressChan       chan float64
	progress           progress.Model
}

// Initialize the model
func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter Hugging Face repo name"
	ti.Focus()

	customPathInput := textinput.New()
	customPathInput.Placeholder = "Enter custom path"

	return model{
		state:           inputRepo,
		repoInput:       ti,
		checked:         make(map[int]bool),
		customPathInput: customPathInput,
		status:          "Ready to download.",
		statusChan:      make(chan string),
		progressChan:    make(chan float64),
		progress:        progress.New(progress.WithScaledGradient("#fd5392", "#f86f64"), progress.WithWidth(80)),
	}
}

// Init function (optional)
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func listenForStatus(ch chan string) tea.Cmd {
	return func() tea.Msg {
		return statusMsg(<-ch)
	}
}

func listenForProgress(ch chan float64) tea.Cmd {
	return func() tea.Msg {
		return progressMsg(<-ch)
	}
}

// Handle user input
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q": // Quit the program
			return m, tea.Quit

		case "enter": // Print input and quit
			if m.state == inputRepo {
				repoID := m.repoInput.Value()
				files, err := cli.ListRepoFiles(repoID)
				if err != nil {
					return m, nil
				}
				m.files = files
				m.state = selectFiles
				m.cursorIndex = 0
				m.scrollOffset = 0
				return m, nil
			} else if m.state == selectFiles {
				selectedFiles := []string{}
				for i, file := range m.files {
					if m.checked[i] {
						selectedFiles = append(selectedFiles, file)
					}
				}
				m.selectedFiles = selectedFiles
				if len(selectedFiles) > 0 {
					m.state = selectDownloadPath
				}
				return m, nil
			} else if m.state == selectDownloadPath {
				var path string
				var err error
				if m.downloadPathChoice == "custom" {
					if m.customPathInput.Value() == "" {
						return m, nil
					}
					path, err = cli.GetDownloadPath("custom", m.customPathInput.Value(), m.repoInput.Value())
				} else {
					path, err = cli.GetDownloadPath(m.downloadPathChoice, "", m.repoInput.Value())
				}

				if err != nil {
					fmt.Println("Error:", err)
					return m, nil
				}
				m.path = path
				m.state = confirmation
				return m, nil
			} else if m.state == confirmation {
				m.state = downloading
				go func() {
					err := cli.Download(
						m.repoInput.Value(),
						m.selectedFiles,
						m.path,
						func(status string) {
							m.statusChan <- status
						},
						func(progress float64) {
							m.progressChan <- progress
						})
					if err != nil {
						m.statusChan <- fmt.Sprintf("Error: %v", err)
					}
					m.statusChan <- "Download complete"
				}()
				return m, tea.Batch(
					listenForStatus(m.statusChan),
					listenForProgress(m.progressChan),
				)
			}

		case "a": //Toggle select all
			allSelected := len(m.checked) == len(m.files)
			m.checked = make(map[int]bool)
			if !allSelected {
				for i := range m.files {
					m.checked[i] = true
				}
			}
			return m, nil

		case " ":
			if m.state == selectFiles && m.cursorIndex < len(m.files) {
				m.checked[m.cursorIndex] = !m.checked[m.cursorIndex]
			}
			return m, nil

			// Left movement (move one column left)
		case "left":
			if (m.cursorIndex-filesPerColumn) >= 0 && (m.cursorIndex/filesPerColumn) > m.scrollOffset {
				m.cursorIndex -= filesPerColumn
			}
			return m, nil

		// Right movement (move one column right)
		case "right":
			if (m.cursorIndex+filesPerColumn) < len(m.files) &&
				(m.cursorIndex/filesPerColumn) < (m.scrollOffset+maxVisibleColumns-1) {
				m.cursorIndex += filesPerColumn
			}
			return m, nil

		// Move cursor up within a column
		case "up":
			if (m.cursorIndex % filesPerColumn) > 0 {
				m.cursorIndex--
			}
			return m, nil

		// Move cursor down within a column
		case "down":
			if (m.cursorIndex%filesPerColumn) < (filesPerColumn-1) && (m.cursorIndex+1) < len(m.files) {
				m.cursorIndex++
			}
			return m, nil

		// Scroll view left `[`
		case "[":
			if m.scrollOffset > 0 {
				m.scrollOffset--
				// Keep cursor within the visible columns
				if (m.cursorIndex / filesPerColumn) < m.scrollOffset {
					m.cursorIndex = m.scrollOffset * filesPerColumn
				}
			}
			return m, nil

		// Scroll view right `]`
		case "]":
			if m.scrollOffset+maxVisibleColumns < (len(m.files)+filesPerColumn-1)/filesPerColumn {
				m.scrollOffset++
				// Ensure cursor stays in the visible area
				if (m.cursorIndex / filesPerColumn) >= (m.scrollOffset + maxVisibleColumns) {
					m.cursorIndex = (m.scrollOffset + maxVisibleColumns - 1) * filesPerColumn
					if m.cursorIndex >= len(m.files) {
						m.cursorIndex = len(m.files) - 1
					}
				}
			}
			return m, nil

		case "1":
			m.downloadPathChoice = "Downloads"
			return m, nil

		case "2":
			m.downloadPathChoice = "Current Directory"
			return m, nil

		case "3":
			m.downloadPathChoice = "custom"
			m.customPathInput.Focus()
			return m, nil
		}

	case statusMsg:
		m.status = string(msg)
		return m, listenForStatus(m.statusChan)

	case progressMsg:
		cmd := m.progress.SetPercent(float64(msg))
		return m, tea.Batch(cmd, listenForProgress(m.progressChan))

	case progress.FrameMsg:
		var cmds []tea.Cmd
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	if m.downloadPathChoice == "custom" {
		m.customPathInput, cmd = m.customPathInput.Update(msg)
		return m, cmd
	}

	if m.state == inputRepo {
		m.repoInput, cmd = m.repoInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func generateColumns(files []string, checked map[int]bool, cursor int, scrollOffset int) string {
	numColumns := (len(files) + filesPerColumn - 1) / filesPerColumn
	visibleEnd := scrollOffset + maxVisibleColumns
	if visibleEnd > numColumns {
		visibleEnd = numColumns
	}

	columnStrings := make([]string, 0, maxVisibleColumns)

	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))

	for col := scrollOffset; col < visibleEnd; col++ {
		start := col * filesPerColumn
		end := start + filesPerColumn
		if end > len(files) {
			end = len(files)
		}

		column := ""
		for i := start; i < end; i++ {
			checkmark := "[ ]"
			if checked[i] {
				checkmark = "[✓]"
			}
			cursorMarker := "  "
			if i == cursor {
				cursorMarker = cursorStyle.Render("➤")
			}
			column += fmt.Sprintf("%s%s %s\n", cursorMarker, checkmark, files[i])
		}
		columnStrings = append(columnStrings, column)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, columnStrings...)
}

// Define the view
func (m model) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5f00")).Padding(1).Align(lipgloss.Center)
	bodyStyle := lipgloss.NewStyle().Padding(1, 2)
	switch m.state {
	case inputRepo:
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Download"),
			bodyStyle.Render(fmt.Sprintf("%s\n\n%s\n\n%s", "Enter Repo Name:", m.repoInput.View(), "Press Enter to confirm, Q to quit")),
		)

	case selectFiles:
		columns := generateColumns(m.files, m.checked, m.cursorIndex, m.scrollOffset)
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Select Files"),
			bodyStyle.BorderStyle(lipgloss.RoundedBorder()).Render(columns+fmt.Sprintf("\nPage %d/%d\n[SPACE] Toggle Select  [A] Toggle All  [←/→] Move  [↑/↓] Move  [ Scroll-Left  ] Scroll-Right  [ENTER] Confirm  [q] Quit", m.scrollOffset+1, (len(m.files)+filesPerColumn-1)/filesPerColumn)),
		)

	case selectDownloadPath:
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Select Download Path"),
			bodyStyle.BorderStyle(lipgloss.RoundedBorder()).Render(fmt.Sprintf("[1] Default (Downloads/hfmodels/%s)\n[2] Current-Directory (%s)", m.repoInput.Value(), m.repoInput.Value())),
			bodyStyle.Render(fmt.Sprintf("Current Choice: %s\nPress ENTER to confirm, Q to quit", m.downloadPathChoice)),
		)

	case confirmation:
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Confirmation"),
			bodyStyle.Render(fmt.Sprintf("Downloading %d files from %s to %s\n\nPress Enter to start, Q to quit", len(m.selectedFiles), m.repoInput.Value(), m.path)),
		)

	case downloading:
		return lipgloss.JoinVertical(lipgloss.Top,
			headerStyle.Render("Downloading..."),
			bodyStyle.Render(m.progress.View()),
			bodyStyle.Render(m.status),
		)
	}

	return ""
}
