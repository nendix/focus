# Focus

A minimalist terminal-based Pomodoro timer built with Go and Bubble Tea.

## Features

- 🍅 **Classic Pomodoro Technique**: 25-minute work sessions, 5-minute short breaks, 15-minute long breaks
- ⏱️ **Customizable Durations**: Edit work and break durations with vim-like controls (1-60 minutes)
- 🎯 **Modal Duration Editor**: Press M to enter edit mode with intuitive keyboard navigation
- 🔊 **Audio Notifications**: Different sounds for work-end and break-end
- ⌨️ **Simple Controls**: Keyboard shortcuts for all operations
- 🎨 **Clean TUI**: Beautiful terminal interface with colors and styling
- 🔄 **Visual Status**: Color-coded timer (Yellow=running, Gray=paused, Green=editing)

## Usage

### Build and Run
```bash
go build -o focus ./cmd/focus
./focus
```

### Controls

#### Normal Mode
- **Space**: Start/Pause timer
- **R**: Reset current phase
- **M**: Enter duration edit mode
- **Q**: Quit application

#### Edit Mode
- **Tab**: Switch between phases (Work → Short Break → Long Break)
- **H/L** or **←/→**: Select digit position (tens/units)
- **J/K** or **↑/↓**: Increment/decrement selected digit
- **M/Esc**: Exit edit mode

## Timer Phases

1. **Work Session** (25 minutes) - Focus time
2. **Short Break** (5 minutes) - Quick rest  
3. After 4 work sessions: **Long Break** (15 minutes) - Extended rest

## Duration Editing

- Press **M** to enter edit mode (timer turns green)
- Use **Tab** to cycle through Work/Short Break/Long Break phases
- Use **H/L** to select tens or units digit
- Use **J/K** to adjust the selected digit (1-60 minute range)
- Selected digit blinks to show current position
- Changes apply immediately

## Audio

The app generates beep tones for notifications:
- **Work End**: Double beep
- **Break End**: Single long beep

## Requirements

- Go 1.19+
- Terminal with color support
- Audio output (for notifications)