package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"focus/pkg/ui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("focus %s (commit: %s, built: %s)\n", version, commit, date)
			return
		case "--help", "-h":
			fmt.Println("Focus - A minimalist pomodoro timer")
			fmt.Println("\nUsage: focus [command]")
			fmt.Println("\nCommands:")
			fmt.Println("  version, -v, --version    Show version information")
			fmt.Println("  help, -h, --help         Show this help message")
			fmt.Println("\nControls:")
			fmt.Println("  Space    Start/pause timer")
			fmt.Println("  R        Reset current session")
			fmt.Println("  Q        Quit")
			fmt.Println("  M        Edit mode (modify durations)")
			fmt.Println("  Tab      Cycle edit phases (in edit mode)")
			fmt.Println("  ←/→, H/L Select digit (in edit mode)")
			fmt.Println("  ↑/↓, J/K Adjust time (in edit mode)")
			fmt.Println("  Esc      Exit edit mode")
			return
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			os.Exit(1)
			return
		}
	}

	model := ui.NewModel()

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
