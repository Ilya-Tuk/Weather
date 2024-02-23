package worker

import (
	"context"
	"time"
)

type Worker struct {
	delay time.Duration
	svc   Service
}

func New(svc Service, delay time.Duration) *Worker {
	return &Worker{
		svc:   svc,
		delay: delay,
	}
}

func (w *Worker) RunNotify(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(w.delay):
		}
		w.svc.SetAlerts(ctx)
	}
}
