package main

import (
	"fmt"
	"lazyface-priv/cmd"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main(){
	token := ""
	addGitCredential := true

	err := cli.Login(token, addGitCredential)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Login completed successfully")
	}

	repoID := "deepseek-ai/DeepSeek-R1"

	fmt.Println("Listing files in repo: ", repoID)
	files, err := cli.ListRepoFiles(repoID)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Files in repo: ")
		for _, file := range files {
			fmt.Println("-", file)
		}
	}

	err = cli.Logout()
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Logout success")
	}
}
