package notification

import (
	"fmt"
	"os/exec"
	"runtime"

	"focus/pkg/timer"
)

type Notifier struct {
	enabled bool
}

func New() *Notifier {
	return &Notifier{
		enabled: runtime.GOOS == "darwin",
	}
}

func (n *Notifier) Show(title, message string) error {
	if !n.enabled {
		return nil
	}

	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "default"`, message, title)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
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
