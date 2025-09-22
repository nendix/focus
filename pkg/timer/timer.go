package timer

import (
	"fmt"
	"time"
)

type Phase int

const (
	Work Phase = iota
	ShortBreak
	LongBreak
)

type Status int

const (
	Running Status = iota
	Paused
	Finished
)

type Timer struct {
	Phase              Phase
	Status             Status
	Duration           time.Duration
	Remaining          time.Duration
	SessionCount       int
	MaxSessions        int
	WorkDuration       time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
	ticker             *time.Ticker
	done               chan bool
	OnPhaseEnd         func(Phase)
}

func New() *Timer {
	return &Timer{
		Phase:              Work,
		Status:             Paused,
		Duration:           25 * time.Minute, // 25 minutes for work
		Remaining:          25 * time.Minute,
		SessionCount:       1,
		MaxSessions:        4,
		WorkDuration:       25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
		done:               make(chan bool),
		// for testing
		// Duration:           5 * time.Second,
		// Remaining:          5 * time.Second,
		// WorkDuration:       5 * time.Second,
		// ShortBreakDuration: 5 * time.Second,
		// LongBreakDuration:  5 * time.Second,
	}
}

func (t *Timer) Start() {
	if t.Status == Running {
		return
	}

	t.Status = Running
	t.ticker = time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-t.ticker.C:
				t.Remaining -= time.Second
				if t.Remaining <= 0 {
					t.finish()
					return
				}
			case <-t.done:
				return
			}
		}
	}()
}

func (t *Timer) Pause() {
	if t.Status != Running {
		return
	}

	t.Status = Paused
	if t.ticker != nil {
		t.ticker.Stop()
	}
}

func (t *Timer) Reset() {
	t.Status = Paused
	if t.ticker != nil {
		t.ticker.Stop()
	}
	t.Remaining = t.getDurationForPhase(t.Phase)
}

func (t *Timer) Stop() {
	t.Status = Paused
	if t.ticker != nil {
		t.ticker.Stop()
	}
	close(t.done)
}

func (t *Timer) finish() {
	t.Status = Finished
	if t.ticker != nil {
		t.ticker.Stop()
	}

	if t.OnPhaseEnd != nil {
		t.OnPhaseEnd(t.Phase)
	}

	t.nextPhase()
}

func (t *Timer) nextPhase() {
	switch t.Phase {
	case Work:
		if t.SessionCount >= t.MaxSessions {
			t.Phase = LongBreak
		} else {
			t.Phase = ShortBreak
		}
	case ShortBreak:
		t.Phase = Work
		t.SessionCount++
	case LongBreak:
		t.Phase = Work
		t.SessionCount = 1
	}

	t.Duration = t.getDurationForPhase(t.Phase)
	t.Remaining = t.Duration
	t.Status = Paused

	// Auto-start the next phase
	t.Start()
}

func (t *Timer) SetWorkDuration(minutes int) {
	t.WorkDuration = time.Duration(minutes) * time.Minute
	// Update current duration and remaining time if we're in a work phase
	if t.Phase == Work {
		t.Duration = t.WorkDuration
		t.Remaining = t.WorkDuration
	}
}

func (t *Timer) SetShortBreakDuration(minutes int) {
	t.ShortBreakDuration = time.Duration(minutes) * time.Minute
	// Update current duration and remaining time if we're in a short break phase
	if t.Phase == ShortBreak {
		t.Duration = t.ShortBreakDuration
		t.Remaining = t.ShortBreakDuration
	}
}

func (t *Timer) SetLongBreakDuration(minutes int) {
	t.LongBreakDuration = time.Duration(minutes) * time.Minute
	// Update current duration and remaining time if we're in a long break phase
	if t.Phase == LongBreak {
		t.Duration = t.LongBreakDuration
		t.Remaining = t.LongBreakDuration
	}
}

func (t *Timer) GetShortBreakDuration() int {
	return int(t.ShortBreakDuration.Minutes())
}

func (t *Timer) GetLongBreakDuration() int {
	return int(t.LongBreakDuration.Minutes())
}

func (t *Timer) GetWorkDuration() int {
	return int(t.WorkDuration.Minutes())
}

func (t *Timer) GetBreakDuration() int {
	return int(t.ShortBreakDuration.Minutes())
}

func (t *Timer) GetDurationForPhase(phase Phase) int {
	switch phase {
	case Work:
		return int(t.WorkDuration.Minutes())
	case ShortBreak:
		return int(t.ShortBreakDuration.Minutes())
	case LongBreak:
		return int(t.LongBreakDuration.Minutes())
	default:
		return int(t.WorkDuration.Minutes())
	}
}

func (t *Timer) SetDurationForPhase(phase Phase, minutes int) {
	switch phase {
	case Work:
		t.SetWorkDuration(minutes)
	case ShortBreak:
		t.SetShortBreakDuration(minutes)
	case LongBreak:
		t.SetLongBreakDuration(minutes)
	}
}

func (t *Timer) getDurationForPhase(phase Phase) time.Duration {
	switch phase {
	case Work:
		return t.WorkDuration
	case ShortBreak:
		return t.ShortBreakDuration
	case LongBreak:
		return t.LongBreakDuration
	default:
		return t.WorkDuration
	}
}

func (t *Timer) GetProgress() float64 {
	if t.Duration == 0 {
		return 0
	}
	elapsed := t.Duration - t.Remaining
	return float64(elapsed) / float64(t.Duration)
}

func (t *Timer) GetPhaseString() string {
	switch t.Phase {
	case Work:
		return "Work Session"
	case ShortBreak:
		return "Short Break"
	case LongBreak:
		return "Long Break"
	default:
		return "Work Session"
	}
}

func (t *Timer) GetStatusString() string {
	switch t.Status {
	case Running:
		return "Running"
	case Paused:
		return "Paused"
	case Finished:
		return "Finished"
	default:
		return "Paused"
	}
}

func (t *Timer) FormatTime() string {
	minutes := int(t.Remaining.Minutes())
	seconds := int(t.Remaining.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
