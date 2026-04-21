package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Handler is a function that receives and processes an Alert.
type Handler func(a *Alert)

// Dispatcher routes alerts to one or more registered handlers.
type Dispatcher struct {
	handlers   []Handler
	threshold  Threshold
}

// NewDispatcher creates a Dispatcher with the given threshold and optional handlers.
func NewDispatcher(t Threshold, handlers ...Handler) *Dispatcher {
	return &Dispatcher{
		handlers:  handlers,
		threshold: t,
	}
}

// Dispatch evaluates the change count for an environment and notifies all
// registered handlers if a threshold is exceeded.
func (d *Dispatcher) Dispatch(env string, changedCount int) {
	a := Evaluate(env, changedCount, d.threshold)
	if a == nil {
		return
	}
	for _, h := range d.handlers {
		h(a)
	}
}

// StdoutHandler returns a Handler that writes alert summaries to stdout.
func StdoutHandler() Handler {
	return WriterHandler(os.Stdout)
}

// WriterHandler returns a Handler that writes formatted alerts to the given writer.
func WriterHandler(w io.Writer) Handler {
	return func(a *Alert) {
		fmt.Fprintf(w, "[%s] ALERT %s — %s (at %s)\n",
			a.Severity,
			a.Environment,
			a.Message,
			a.TriggeredAt.Format(time.RFC3339),
		)
	}
}
