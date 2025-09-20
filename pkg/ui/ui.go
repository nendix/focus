package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"pomodoro/pkg/ascii"
	"pomodoro/pkg/audio"
	"pomodoro/pkg/timer"
)

type Model struct {
	timer  *timer.Timer
	audio  *audio.Player
	width  int
	height int
}

type tickMsg struct{}

func NewModel() *Model {
	// Create timer
	t := timer.New()

	// Create audio player
	audioPlayer, err := audio.New()
	if err != nil {
		// Continue without audio if it fails
		audioPlayer = nil
	}

	m := &Model{
		timer: t,
		audio: audioPlayer,
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
			m.timer.Stop()
			if m.audio != nil {
				m.audio.Close()
			}
			return m, tea.Quit
		case " ":
			if m.timer.Status == timer.Running {
				m.timer.Pause()
			} else if m.timer.Status == timer.Paused {
				m.timer.Start()
			}
		case "r":
			m.timer.Reset()
		}

	case tickMsg:
		return m, m.tickCmd()
	}

	return m, nil
}

func (m *Model) View() string {
	// Phase info
	phaseColor := lipgloss.Color("2")
	if m.timer.Phase == timer.Work {
		phaseColor = lipgloss.Color("4")
	}

	phaseText := lipgloss.NewStyle().
		Foreground(phaseColor).
		Bold(true).
		Render(fmt.Sprintf("%s (%d/4)", m.timer.GetPhaseString(), m.timer.SessionCount))

	// Time display (ASCII art) - color based on status
	timeStr := m.timer.FormatTime()
	asciiTime := ascii.ToASCII(timeStr)

	// Timer color based on status: running=yellow, paused=gray
	timerColor := lipgloss.Color("8") // Default gray
	switch m.timer.Status {
	case timer.Running:
		timerColor = lipgloss.Color("3") // Yellow
	case timer.Paused:
		timerColor = lipgloss.Color("8") // Gray
	}

	timeText := lipgloss.NewStyle().
		Bold(true).
		Foreground(timerColor).
		Align(lipgloss.Center).
		Padding(1, 2).
		Render(asciiTime)

	// Controls
	controls := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Render("[Space] Start/Pause  [R] Reset  [Q] Quit")

	// Layout
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		phaseText,
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
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *Model) onPhaseEnd(phase timer.Phase) {
	if m.audio != nil {
		switch phase {
		case timer.Work:
			go m.audio.PlayWorkEndSound()
		case timer.ShortBreak, timer.LongBreak:
			go m.audio.PlayBreakEndSound()
		}
	}
}
