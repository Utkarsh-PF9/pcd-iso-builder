package main

import (
	"fmt"
	"os"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(loadingModel{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}
