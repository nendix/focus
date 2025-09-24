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
	timer      *timer.Timer
	notifier   *notification.Notifier
	width      int
	height     int
	editMode   bool
	helpMode   bool
	editPhase  timer.Phase
	blinkState bool
	lastBlink  time.Time
}

type tickMsg struct{}

func NewModel() *Model {
	// Create timer
	t := timer.New()

	// Create notifier
	notifier := notification.New()

	m := &Model{
		timer:      t,
		notifier:   notifier,
		editMode:   false,
		helpMode:   false,
		editPhase:  timer.Work,
		blinkState: true,
		lastBlink:  time.Now(),
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
		case "?":
			m.toggleHelpMode()
		case "esc":
			if m.helpMode {
				m.helpMode = false
			} else if m.editMode {
				m.editMode = false
			}
		case " ":
			switch m.timer.Status {
			case timer.Running:
				m.timer.Stop()
			case timer.Stopped:
				m.timer.Start()
			}
		case "r":
			m.timer.Reset()
		case "tab":
			if m.editMode {
				m.nextEditPhase()
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
	if m.helpMode {
		return m.renderHelpScreen()
	}

	if m.editMode {
		return m.renderEditMode()
	}

	return m.renderMainTimer()
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

func (m *Model) toggleHelpMode() {
	if !m.helpMode {
		m.helpMode = true
	} else if m.helpMode {
		m.helpMode = false
	}
}

func (m *Model) toggleEditMode() {
	if !m.editMode {
		m.editMode = true
		m.editPhase = m.timer.Phase
		m.blinkState = true
		m.lastBlink = time.Now()
	} else if m.editMode {
		m.editMode = false
	}
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
}

func (m *Model) incrementDigit() {
	currentMinutes := m.timer.GetDurationForPhase(m.editPhase)
	newMinutes := m.adjustDigit(currentMinutes, 5)
	if newMinutes >= 1 && newMinutes <= 60 {
		m.timer.SetDurationForPhase(m.editPhase, newMinutes)
	}
}

func (m *Model) decrementDigit() {
	currentMinutes := m.timer.GetDurationForPhase(m.editPhase)
	newMinutes := m.adjustDigit(currentMinutes, -5)
	if newMinutes >= 1 && newMinutes <= 60 {
		m.timer.SetDurationForPhase(m.editPhase, newMinutes)
	}
}

func (m *Model) adjustDigit(minutes, delta int) int {
	// Simple ±5 minute adjustment
	newValue := minutes + delta

	// Handle wraparound for 1-60 range
	if newValue < 1 {
		// Wrap from bottom: calculate how far below 1, then wrap from 60
		underflow := 1 - newValue
		return 60 - (underflow - 1)
	} else if newValue > 60 {
		// Wrap from top: calculate how far above 60, then wrap from 1
		overflow := newValue - 60
		return overflow
	}

	return newValue
}

func (m *Model) formatEditTime(minutes int) string {
	tens := minutes / 10
	units := minutes % 10

	// Create the time string with potential blinking
	timeStr := fmt.Sprintf("%d%d:00", tens, units)

	// Make entire time blink (replace with spaces to maintain width)
	if !m.blinkState {
		timeStr = " " // 5 spaces to match "XX:00" width
	}

	return timeStr
}

func (m *Model) getPhaseString(phase timer.Phase) string {
	switch phase {
	case timer.Work:
		return "work"
	case timer.ShortBreak:
		return "break"
	case timer.LongBreak:
		return "long break"
	default:
		return "work"
	}
}

func (m *Model) renderHelpScreen() string {
	helpContent := `focus - A minimalist pomodoro timer

Controls:
  ?               toggle help
  Space           start/stop timer
  R               reset current session  
  E               toggle edit mode 
  Esc             exit from any mode
  Q               quit

Edit mode:
  Tab             switch phase (Work/Break/Long Break)
  J/K or ↓/↑      adjust time ±5 minutes
`

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Align(lipgloss.Left).
		Render(helpContent)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("4")).
			Padding(1, 2).
			Render(help),
	)
}

func (m *Model) renderEditMode() string {
	phaseColor := lipgloss.Color("2")
	phaseText := fmt.Sprintf("edit: %s duration", m.getPhaseString(m.editPhase))

	phaseDisplay := lipgloss.NewStyle().
		Foreground(phaseColor).
		Render(phaseText)

	// Time display (ASCII art) - show edit time
	minutes := m.timer.GetDurationForPhase(m.editPhase)
	timeStr := m.formatEditTime(minutes)
	asciiTime := ascii.ToASCII(timeStr)

	timeText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Align(lipgloss.Center).
		Padding(1, 2).
		Render(asciiTime)

	// Layout
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		phaseDisplay,
		"",
		timeText,
	)

	// Center the content
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("2")).
			Padding(1, 2).
			Render(content),
	)
}

func (m *Model) renderMainTimer() string {
	// Phase info
	var phaseColor lipgloss.Color
	var phaseText string

	// Check if timer is paused - use gray colors
	if m.timer.Status == timer.Stopped {
		phaseColor = lipgloss.Color("8") // gray
		if m.timer.Phase == timer.Work {
			phaseText = fmt.Sprintf("%s (%d/4)", m.getPhaseString(m.timer.Phase), m.timer.SessionCount)
		} else {
			phaseText = m.getPhaseString(m.timer.Phase)
		}
	} else {
		// Normal colors when running
		if m.timer.Phase == timer.Work {
			phaseColor = lipgloss.Color("4")
			phaseText = fmt.Sprintf("%s (%d/4)", m.getPhaseString(m.timer.Phase), m.timer.SessionCount)
		} else {
			phaseColor = lipgloss.Color("3")
			phaseText = m.getPhaseString(m.timer.Phase)
		}
	}

	phaseDisplay := lipgloss.NewStyle().
		Foreground(phaseColor).
		Render(phaseText)

	// Time display (ASCII art)
	timeStr := m.timer.FormatTime()
	asciiTime := ascii.ToASCII(timeStr)

	// Timer color based on status
	timerColor := lipgloss.Color("7") // white
	switch m.timer.Status {
	case timer.Running:
		timerColor = lipgloss.Color("7") // white
	case timer.Stopped:
		timerColor = lipgloss.Color("8") // Gray
	}

	timeText := lipgloss.NewStyle().
		Foreground(timerColor).
		Align(lipgloss.Center).
		Padding(1, 2).
		Render(asciiTime)

	// Layout
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		phaseDisplay,
		"",
		timeText,
	)

	// Center the content
	// Border color - gray when paused, normal when running
	borderColor := lipgloss.Color("7") // default white
	if m.timer.Status == timer.Stopped {
		borderColor = lipgloss.Color("8") // gray when paused
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Render(content),
	)
}
