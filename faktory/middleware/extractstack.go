package middleware

import (
	"context"

	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
)

// SetStackNameOnCtx add Queue name on context
func SetStackNameOnCtx(perform worker.Handler) worker.Handler {
	return func(ctx worker.Context, job *faktory.Job) (err error) {
		var stack string
		if s, ok := job.Custom["stack"]; ok {
			stack = s.(string)
		} else {
			stack = "default"
		}

		ctx = &worker.DefaultContext{
			Context: context.WithValue(ctx, "Stack", stack),
			JID:     job.Jid,
			Type:    job.Type,
		}

		return perform(ctx, job)
	}
}
