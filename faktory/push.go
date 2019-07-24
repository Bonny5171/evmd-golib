package faktory

import (
	"time"

	faktory "github.com/contribsys/faktory/client"
	"github.com/pkg/errors"
)

func Push(jobName, queue, stack, dsn string, retry int, at time.Time, params []interface{}) error {
	custom := map[string]interface{}{
		"dsn":   dsn,
		"stack": stack,
	}

	if err := PushCustom(jobName, queue, retry, at, params, custom); err != nil {
		return errors.Wrap(err, "PushCustom()")
	}

	return nil
}

func PushCustom(jobName, queue string, retry int, at time.Time, params []interface{}, custom map[string]interface{}) error {
	cl, err := faktory.Open()
	if err != nil {
		time.Sleep(5 * time.Second)
		cl, err = faktory.Open()
		if err != nil {
			return errors.Wrap(err, "faktory.Open()")
		}
	}

	job := faktory.NewJob(jobName, params...)
	job.Queue = queue
	job.Retry = int(retry)

	zeroTime := time.Time{}
	if at != zeroTime {
		job.At = at.Format(time.RFC3339Nano)
	}

	job.Custom = custom

	if err = cl.Push(job); err != nil {
		return errors.Wrap(err, "cl.Push()")
	}

	return nil
}
