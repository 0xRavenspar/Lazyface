package main

import (
	"Lazyface/cmd"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	//
	// p := tea.NewProgram(cmd.InitialModel())
	//
	// if err := p.Start(); err != nil {
	// 	fmt.Println("Error starting tje TUI: ", err)
	// 	os.Exit(1)
	// }
	//
	u := tea.NewProgram(cmd.InitialSplashModel())
	if _, err := u.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	//
	// manageModel := cmd.InitialManageModel()
	// q := tea.NewProgram(manageModel)
	// if _, err := q.Run(); err != nil {
	// 	fmt.Printf("Error: %v", err)
	// 	os.Exit(1)
	// }
	//
	// r := tea.NewProgram(
	// 	cmd.NewAuthView(),
	// 	// tea.WithAltScreen(),       // Use alternate screen buffer
	// 	tea.WithMouseCellMotion(), // Turn on mouse support
	// )
	//
	// if _, err := r.Run(); err != nil {
	// 	fmt.Printf("Error running program: %v", err)
	// 	os.Exit(1)
	// }
	//
}
