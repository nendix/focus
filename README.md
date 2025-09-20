# Pomodoro TUI

A terminal-based Pomodoro Technique timer built with Go and Bubble Tea.

## Features

- 🍅 **Classic Pomodoro Technique**: 25-minute work sessions, 5-minute short breaks, 15-minute long breaks
- 📊 **Visual Progress Bar**: Real-time countdown with colorful progress indication
- 🔊 **Audio Notifications**: Different sounds for work-end and break-end
- ⌨️ **Simple Controls**: Keyboard shortcuts for all operations
- 🎨 **Clean TUI**: Beautiful terminal interface with colors and styling

## Usage

### Build and Run
```bash
go build ./cmd/pomodoro
./pomodoro
```

### Controls
- **Space**: Start/Pause timer
- **R**: Reset current phase
- **Q**: Quit application

## Timer Phases

1. **Work Session** (25 minutes) - Red progress bar
2. **Short Break** (5 minutes) - Green progress bar  
3. After 4 work sessions: **Long Break** (15 minutes) - Green progress bar

## Audio

The app generates beep tones for notifications:
- **Work End**: Double beep
- **Break End**: Single long beep

## Requirements

- Go 1.19+
- Terminal with color support
- Audio output (for notifications)