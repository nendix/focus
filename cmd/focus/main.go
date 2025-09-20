package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"focus/pkg/ui"
)

func main() {
	model := ui.NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
