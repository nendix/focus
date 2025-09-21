# Agent Guidelines for Focus Pomodoro Timer

## Commands
- Build: `go build -o out/focus ./cmd/focus` or `make build`
- Run: `./out/focus` or `make run`
- Test: `go test ./...` or `make test`
- Test single package: `go test ./pkg/timer`
- Format: `go fmt ./...` or `make fmt`
- Lint: `go vet ./...` or `make vet`
- Install: `go install ./cmd/focus` or `make install`

## Code Style
- Use Go standard formatting (`go fmt`)
- Package imports: stdlib first, then external, then local (`focus/pkg/...`)
- Naming: CamelCase for exported, camelCase for unexposed
- Types: Define enums using `iota`, group constants
- Error handling: Check all errors explicitly, return early on errors
- Structs: Group related fields, use pointer receivers for methods that modify state
- Functions: Keep functions small and focused, use descriptive names
- Comments: Document exported functions and types, avoid obvious comments
- UI Colors: Use ANSI color codes only for better terminal compatibility
  - Green (success/break): lipgloss.Color("2")
  - Red (work/alerts): lipgloss.Color("1") 
  - Yellow (time/warnings): lipgloss.Color("3")
  - Gray (inactive): lipgloss.Color("8")
  - Magenta (finished/borders): lipgloss.Color("5")
  - White (text): lipgloss.Color("7")

## Architecture
- Main entry: `cmd/focus/main.go`
- Packages: `pkg/{timer,ui,notification,ascii}` with clear separation of concerns
- UI: Uses Bubble Tea framework for terminal UI
- Notifications: Cross-platform using beeep library (macOS, Linux, Windows)
- Dependencies: Minimal external dependencies (Charm libraries for UI, beeep for notifications)