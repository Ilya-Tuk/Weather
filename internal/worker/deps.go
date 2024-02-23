package worker

import "context"

type Service interface {
	SetAlerts(context.Context)
}
