package notification

import (
	"github.com/gen2brain/beeep"

	"github.com/nendix/focus/pkg/timer"
)

type Notifier struct {
	enabled bool
}

func New() *Notifier {
	return &Notifier{
		enabled: true, // Works cross-platform
	}
}

func (n *Notifier) Show(title, message string) error {
	if !n.enabled {
		return nil
	}

	beeep.AppName = "Focus"
	return beeep.Alert(title, message, "")
}

func (n *Notifier) NotifyPhaseEnd(phase timer.Phase) {
	if !n.enabled {
		return
	}

	var title, message string
	switch phase {
	case timer.Work:
		title = "Break time!"
		message = "Your work is over."
	case timer.ShortBreak:
		title = "Work time!"
		message = "Your break is over."
	case timer.LongBreak:
		title = "Work time!"
		message = "Your long break is over."
	default:
		return
	}

	go func() {
		if err := n.Show(title, message); err != nil {
			// Silently fail if notification can't be shown
		}
	}()
}
