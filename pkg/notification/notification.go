package notification

import (
	"github.com/gen2brain/beeep"

	"focus/pkg/timer"
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

	return beeep.Alert(title, message, "../../assets/icon.png")
}

func (n *Notifier) NotifyPhaseEnd(phase timer.Phase) {
	if !n.enabled {
		return
	}

	var title, message string
	switch phase {
	case timer.Work:
		title = "Focus"
		message = "Work session complete!"
	case timer.ShortBreak:
		title = "Focus"
		message = "Short break over!"
	case timer.LongBreak:
		title = "Focus"
		message = "Long break over!"
	default:
		return
	}

	go func() {
		if err := n.Show(title, message); err != nil {
			// Silently fail if notification can't be shown
		}
	}()
}
