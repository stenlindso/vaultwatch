package watch

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Notifier formats and writes watch Events to an output writer.
type Notifier struct {
	out io.Writer
}

// NewNotifier creates a Notifier writing to out.
func NewNotifier(out io.Writer) *Notifier {
	return &Notifier{out: out}
}

// Handle writes a human-readable summary of the event.
func (n *Notifier) Handle(ev Event) {
	timestamp := ev.At.Format(time.RFC3339)
	fmt.Fprintf(n.out, "[%s] env=%s\n", timestamp, ev.Env)
	if len(ev.Added) > 0 {
		fmt.Fprintf(n.out, "  + added (%d):\n", len(ev.Added))
		for _, p := range ev.Added {
			fmt.Fprintf(n.out, "      %s\n", p)
		}
	}
	if len(ev.Removed) > 0 {
		fmt.Fprintf(n.out, "  - removed (%d):\n", len(ev.Removed))
		for _, p := range ev.Removed {
			fmt.Fprintf(n.out, "      %s\n", p)
		}
	}
}

// HandleAll drains the channel and handles each event.
func (n *Notifier) HandleAll(ch <-chan Event) {
	for ev := range ch {
		n.Handle(ev)
	}
}

// FormatEvent returns the event as a compact string (useful for logging).
func FormatEvent(ev Event) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "env=%s added=%d removed=%d", ev.Env, len(ev.Added), len(ev.Removed))
	return sb.String()
}
