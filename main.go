package main

import (
	"Lazyface/cmd"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	p := tea.NewProgram(cmd.InitialModel())

	if err := p.Start(); err != nil {
		fmt.Println("Error starting tje TUI: ", err)
		os.Exit(1)
	}

}
