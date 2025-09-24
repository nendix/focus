package main

import (
	"fmt"
	"os"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nendix/focus/pkg/ui"
)

func getVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "(devel)" && info.Main.Version != "" {
			return info.Main.Version
		}
	}

	return "dev"
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("focus %s\n", getVersion())
			return
		case "--help", "-h":
			fmt.Println("focus - A minimalist pomodoro timer")
			fmt.Println("\nUsage: focus [command]")
			fmt.Println("\nCommands:")
			fmt.Println("  -v, --version      show version information")
			fmt.Println("  -h, --help         show this help message")
			fmt.Println("\nControls:")
			fmt.Println("  ?                  toggle help")
			fmt.Println("  Space              start/stop timer")
			fmt.Println("  R                  restart timer")
			fmt.Println("  E                  toggle edit mode")
			fmt.Println("  Esc                exit from any mode")
			fmt.Println("  Q                  quit app")
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
