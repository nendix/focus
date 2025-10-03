package timer

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
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
	Stopped
	Finished
)

type Timer struct {
	Phase              Phase
	Status             Status
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
	godotenv.Load()

	// Check for development mode
	if os.Getenv("ENV") == "dev" {
		return &Timer{
			Phase:              Work,
			Remaining:          5 * time.Second,
			WorkDuration:       5 * time.Second,
			ShortBreakDuration: 3 * time.Second,
			LongBreakDuration:  5 * time.Second,
			Status:             Running,
			SessionCount:       1,
			MaxSessions:        4,
			done:               make(chan bool),
		}
	}

	// Default production durations
	return &Timer{
		Phase:              Work,
		Remaining:          25 * time.Minute,
		WorkDuration:       25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
		Status:             Running,
		SessionCount:       1,
		MaxSessions:        4,
		done:               make(chan bool),
	}
}

func (t *Timer) Start() {
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

func (t *Timer) Stop() {
	t.Status = Stopped
	if t.ticker != nil {
		t.ticker.Stop()
	}
}

func (t *Timer) Reset() {
	t.Status = Stopped
	if t.ticker != nil {
		t.ticker.Stop()
	}
	t.Remaining = t.getDurationForPhase(t.Phase)
}

func (t *Timer) Quit() {
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

	t.Remaining = t.getDurationForPhase(t.Phase)

	// Auto-start the next phase
	t.Start()
}

func (t *Timer) GetWorkDuration() int {
	return int(t.WorkDuration.Minutes())
}

func (t *Timer) SetWorkDuration(minutes int) {
	t.WorkDuration = time.Duration(minutes) * time.Minute
	// Update remaining time if we're in a work phase
	if t.Phase == Work {
		t.Remaining = t.WorkDuration
	}
}

func (t *Timer) GetShortBreakDuration() int {
	return int(t.ShortBreakDuration.Minutes())
}

func (t *Timer) SetShortBreakDuration(minutes int) {
	t.ShortBreakDuration = time.Duration(minutes) * time.Minute
	// Update remaining time if we're in a short break phase
	if t.Phase == ShortBreak {
		t.Remaining = t.ShortBreakDuration
	}
}

func (t *Timer) GetLongBreakDuration() int {
	return int(t.LongBreakDuration.Minutes())
}

func (t *Timer) SetLongBreakDuration(minutes int) {
	t.LongBreakDuration = time.Duration(minutes) * time.Minute
	// Update remaining time if we're in a long break phase
	if t.Phase == LongBreak {
		t.Remaining = t.LongBreakDuration
	}
}

func (t *Timer) GetDurationForPhase(phase Phase) int {
	return int(t.getDurationForPhase(phase).Minutes())
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

func (t *Timer) FormatTime() string {
	minutes := int(t.Remaining.Minutes())
	seconds := int(t.Remaining.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
