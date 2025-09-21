# Focus

Strongly vibe-coded-minimalist TUI for the Pomodoro timer technique built with Go and Bubble Tea.

## Features

- **Pomodoro Technique**: 25-minute work sessions, 5-minute short breaks, 15-minute long breaks
- **Customizable Durations**: Edit work and break durations with vim-like controls (1-60 minutes)
- **Audio Notifications**: Different sounds for work phase ending and break phase ending
- **Simple Controls**: Keyboard commands for all operations
- **Clean TUI**: Minimal terminal interface with colors and styling

## Usage

### Build and Run
```bash
git clone https://github.com/nendix/focus.git
cd focus
make run
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

## Default Timer Durations

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

## Requirements

- Go 1.19+
- Terminal with color support
- Audio output (for notifications)
