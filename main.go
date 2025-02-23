package main

import (
    "Lazyface/cmd"
    "fmt"
    "time"

    "github.com/qeesung/image2ascii/convert"
    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    // Path to the PNG file
    pngFile := "huggingface.png"

    // Create a new converter
    converter := convert.NewImageConverter()

    // Convert the image to ASCII art
    options := convert.DefaultOptions
    options.FixedWidth = 100 // Adjust width as needed
    options.FixedHeight = 50 // Adjust height as needed

    asciiArt := converter.ImageFile2ASCIIString(pngFile, &options)
    fmt.Println(asciiArt)

    // Wait for 10 seconds
    fmt.Println("Waiting for 10 seconds...")
    time.Sleep(10 * time.Second)

	// token := ""
	// addGitCredential := true
	//
	// err := cli.Login(token, addGitCredential)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// } else {
	// 	fmt.Println("Login completed successfully")
	// }

	model, err := cmd.InitialSettingsModel()
	if err != nil {
		fmt.Println("Error initializing settings model:", err)
		return
	}

	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		fmt.Println("Error running settings view:", err)
	}
	//
	// p := tea.NewProgram(cmd.InitialModel())
	//
	// if err := p.Start(); err != nil {
	// 	fmt.Println("Error starting tje TUI: ", err)
	// 	os.Exit(1)
	// }
	//
	// u := tea.NewProgram(cmd.InitialSplashModel())
	// if _, err := u.Run(); err != nil {
	// 	fmt.Printf("Error: %v", err)
	// 	os.Exit(1)
	// }
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
