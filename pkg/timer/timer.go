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
	Phase        Phase
	Status       Status
	Duration     time.Duration
	Remaining    time.Duration
	SessionCount int
	MaxSessions  int
	ticker       *time.Ticker
	done         chan bool
	OnPhaseEnd   func(Phase)
}

func New() *Timer {
	return &Timer{
		Phase:        Work,
		Status:       Paused,
		Duration:     3 * time.Second, // 25 minutes for work
		Remaining:    3 * time.Second,
		SessionCount: 1,
		MaxSessions:  4,
		done:         make(chan bool),
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
		t.SessionCount++
		if t.SessionCount >= t.MaxSessions {
			t.Phase = LongBreak
			t.SessionCount = 0
		} else {
			t.Phase = ShortBreak
		}
	case ShortBreak, LongBreak:
		t.Phase = Work
	}

	t.Duration = t.getDurationForPhase(t.Phase)
	t.Remaining = t.Duration
	t.Status = Paused

	// Auto-start the next phase
	t.Start()
}

func (t *Timer) getDurationForPhase(phase Phase) time.Duration {
	switch phase {
	case Work:
		return 25 * time.Minute
	case ShortBreak:
		return 5 * time.Minute
	case LongBreak:
		return 15 * time.Minute
	default:
		return 25 * time.Minute
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
