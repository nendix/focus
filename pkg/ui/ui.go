package ui

import (
	"fmt"
	"time"

	"github.com/nendix/focus/pkg/timer"

	"github.com/nendix/focus/pkg/notification"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/nendix/focus/pkg/ascii"
)

type Model struct {
	timer        *timer.Timer
	notifier     *notification.Notifier
	width        int
	height       int
	editMode     bool
	editPhase    timer.Phase
	editDigitPos int
	blinkState   bool
	lastBlink    time.Time
}

type tickMsg struct{}

func NewModel() *Model {
	// Create timer
	t := timer.New()

	// Create notifier
	notifier := notification.New()

	m := &Model{
		timer:        t,
		notifier:     notifier,
		editMode:     false,
		editPhase:    timer.Work,
		editDigitPos: 0,
		blinkState:   true,
		lastBlink:    time.Now(),
	}

	// Set up phase end callback
	t.OnPhaseEnd = m.onPhaseEnd

	// Auto-start the timer when the application launches
	t.Start()

	return m
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.tickCmd(),
		tea.EnterAltScreen,
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.timer.Quit()
			return m, tea.Quit
		case "e":
			m.toggleEditMode()
		case "escape":
			if m.editMode {
				m.exitEditMode()
			}
		case " ":
			if !m.editMode {
				if m.timer.Status == timer.Running {
					m.timer.Stop()
				} else if m.timer.Status == timer.Stopped {
					m.timer.Start()
				}
			}
		case "r":
			if !m.editMode {
				m.timer.Reset()
			}
		case "tab":
			if m.editMode {
				m.nextEditPhase()
			}
		case "h", "left":
			if m.editMode {
				m.editDigitPos = 0 // tens digit
			}
		case "l", "right":
			if m.editMode {
				m.editDigitPos = 1 // units digit
			}
		case "j", "down":
			if m.editMode {
				m.decrementDigit()
			}
		case "k", "up":
			if m.editMode {
				m.incrementDigit()
			}
		}

	case tickMsg:
		// Handle blinking in edit mode
		if m.editMode && time.Since(m.lastBlink) > 500*time.Millisecond {
			m.blinkState = !m.blinkState
			m.lastBlink = time.Now()
		}
		return m, m.tickCmd()
	}

	return m, nil
}

func (m *Model) View() string {
	// Phase info - show edit phase in edit mode
	phaseColor := lipgloss.Color("2")
	var phaseText string
	if m.editMode {
		if m.editPhase == timer.Work {
			phaseColor = lipgloss.Color("4")
		}
		phaseText = fmt.Sprintf("EDIT: %s Duration", m.getPhaseString(m.editPhase))
	} else {
		if m.timer.Phase == timer.Work {
			phaseColor = lipgloss.Color("4")
			phaseText = fmt.Sprintf("%s (%d/%d)", m.getPhaseString(m.timer.Phase), m.timer.SessionCount, m.timer.MaxSessions)
		} else {
			phaseText = fmt.Sprintf("%s", m.getPhaseString(m.timer.Phase))
		}
	}

	phaseDisplay := lipgloss.NewStyle().
		Foreground(phaseColor).
		Bold(true).
		Render(phaseText)

	// Time display (ASCII art) - show edit time in edit mode or regular time
	var timeStr string
	if m.editMode {
		minutes := m.timer.GetDurationForPhase(m.editPhase)
		timeStr = m.formatEditTime(minutes)
	} else {
		timeStr = m.timer.FormatTime()
	}
	asciiTime := ascii.ToASCII(timeStr)

	// Timer color based on mode: edit=green, running=yellow, paused=gray
	timerColor := lipgloss.Color("3") // Default yellow
	if m.editMode {
		timerColor = lipgloss.Color("2") // Green in edit mode
	} else {
		switch m.timer.Status {
		case timer.Running:
			timerColor = lipgloss.Color("3") // Yellow
		case timer.Stopped:
			timerColor = lipgloss.Color("8") // Gray
		}
	}

	timeText := lipgloss.NewStyle().
		Bold(true).
		Foreground(timerColor).
		Align(lipgloss.Center).
		Padding(1, 2).
		Render(asciiTime)

	// Controls - different text for edit mode vs normal mode
	var controlsText string
	if m.editMode {
		controlsText = "[Tab] Switch Phase  [H/L] Select Digit  [J/K] Adjust Value  [E] Exit [Q] Quit"
	} else {
		controlsText = "[Space] Start/Pause  [R] Reset  [E] Edit [Q] Quit"
	}

	controls := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Render(controlsText)

	// Layout
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		phaseDisplay,
		"",
		timeText,
		"",
		controls,
	)

	// Center the content
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("7")).
			Padding(1, 2).
			Render(content),
	)
}

func (m *Model) tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *Model) onPhaseEnd(phase timer.Phase) {
	// Show system notification
	if m.notifier != nil {
		m.notifier.NotifyPhaseEnd(phase)
	}
}

func (m *Model) toggleEditMode() {
	if m.editMode {
		m.exitEditMode()
	} else {
		m.enterEditMode()
	}
}

func (m *Model) enterEditMode() {
	m.editMode = true
	m.editPhase = m.timer.Phase
	m.editDigitPos = 0
	m.blinkState = true
	m.lastBlink = time.Now()
}

func (m *Model) exitEditMode() {
	m.editMode = false
}

func (m *Model) nextEditPhase() {
	switch m.editPhase {
	case timer.Work:
		m.editPhase = timer.ShortBreak
	case timer.ShortBreak:
		m.editPhase = timer.LongBreak
	case timer.LongBreak:
		m.editPhase = timer.Work
	}
	m.editDigitPos = 0
}

func (m *Model) incrementDigit() {
	currentMinutes := m.timer.GetDurationForPhase(m.editPhase)
	newMinutes := m.adjustDigit(currentMinutes, 1)
	if newMinutes >= 1 && newMinutes <= 60 {
		m.timer.SetDurationForPhase(m.editPhase, newMinutes)
	}
}

func (m *Model) decrementDigit() {
	currentMinutes := m.timer.GetDurationForPhase(m.editPhase)
	newMinutes := m.adjustDigit(currentMinutes, -1)
	if newMinutes >= 1 && newMinutes <= 60 {
		m.timer.SetDurationForPhase(m.editPhase, newMinutes)
	}
}

func (m *Model) adjustDigit(minutes, delta int) int {
	tens := minutes / 10
	units := minutes % 10

	if m.editDigitPos == 0 { // tens digit
		newTens := tens + delta
		if newTens < 0 {
			newTens = 6 // wrap to 60
			units = 0   // force units to 0 for 60
		} else if newTens > 6 {
			newTens = 0 // wrap to 0x
		} else if newTens == 6 {
			units = 0 // force units to 0 when tens = 6 (only 60 allowed)
		}
		return newTens*10 + units
	} else { // units digit
		newUnits := units + delta
		if tens == 6 {
			// Special case: when tens=6, units must stay 0 (only 60 allowed)
			return 60
		} else if newUnits < 0 {
			if tens == 0 {
				// At 0x, wrap to 60
				return 60
			} else {
				// Normal wrap to 9
				newUnits = 9
			}
		} else if newUnits > 9 {
			if tens == 5 && newUnits == 10 {
				// 59 -> 60 (special case)
				return 60
			} else {
				// Normal wrap to 0
				newUnits = 0
			}
		}
		return tens*10 + newUnits
	}
}

func (m *Model) formatEditTime(minutes int) string {
	tens := minutes / 10
	units := minutes % 10

	// Create the time string with potential blinking
	tensStr := fmt.Sprintf("%d", tens)
	unitsStr := fmt.Sprintf("%d", units)

	// Make selected digit blink (replace with space to maintain width)
	if !m.blinkState {
		if m.editDigitPos == 0 {
			tensStr = " " // Use space to maintain character width
		} else {
			unitsStr = " " // Use space to maintain character width
		}
	}

	return fmt.Sprintf("%s%s:00", tensStr, unitsStr)
}

func (m *Model) getPhaseString(phase timer.Phase) string {
	switch phase {
	case timer.Work:
		return "Work"
	case timer.ShortBreak:
		return "Break"
	case timer.LongBreak:
		return "Long Break"
	default:
		return "Work"
	}
}
