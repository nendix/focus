# Focus

Strongly vibe-coded minimal TUI for the Pomodoro timer technique built with Go and Bubble Tea.

## Features

- **Pomodoro Technique**: 25-minute work sessions, 5-minute short breaks, 15-minute long breaks
- **Customizable Durations**: Edit work and break durations (1-60 minutes) with ±5 minute adjustments
- **Audio Notifications**: Different sounds for work phase ending and break phase ending
- **Simple Controls**: Keyboard commands for all operations
- **Clean TUI**: Minimal terminal interface with colors and styling

## Installation
```bash
go install github.com/nendix/focus/cmd/focus@latest
```

## Controls

#### Normal Mode
- **?**: Toggle help screen
- **Space**: Start/Pause timer
- **R**: Reset current phase
- **E**: Toggle duration edit mode
- **Q**: Quit application

#### Edit Mode
- **Tab**: Switch between phases (Work → Short Break → Long Break)  
- **J/K** or **↑/↓**: Adjust time ±5 minutes
- **Esc**: Exit edit mode

## Default Timer Durations

1. **Work Session** (25 minutes) - Focus time
2. **Short Break** (5 minutes) - Quick rest  
3. **Long Break** (15 minutes) - Extended rest after 4 work sessions

## Requirements

- Go 1.19+
- Terminal with color support
- Audio output (for notifications)
