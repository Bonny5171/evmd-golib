package push

import (
	"fmt"
	"time"

	faktory "github.com/contribsys/faktory/client"
)

func Push(jobName, queue, stack, dsn string, retry int, at time.Time, params []interface{}) error {
	custom := map[string]interface{}{
		"dsn":   dsn,
		"stack": stack,
	}

	if err := PushCustom(jobName, queue, retry, at, params, custom); err != nil {
		return fmt.Errorf("PushCustom(): %w", err)
	}

	return nil
}

func PushCustom(jobName, queue string, retry int, at time.Time, params []interface{}, custom map[string]interface{}) error {
	cl, err := faktory.Open()
	if err != nil {
		time.Sleep(5 * time.Second)
		cl, err = faktory.Open()
		if err != nil {
			return fmt.Errorf("faktory.Open(): %w", err)
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
		return fmt.Errorf("cl.Push(): %w", err)
	}

	return nil
}

func RetryLater(jobName, queue, stack, dsn string, params []interface{}, after time.Duration) error {
	at := time.Now().Add(after)

	if err := Push(jobName, queue, stack, dsn, 1, at, params); err != nil {
		return fmt.Errorf("faktory.Push(): %w", err)
	}

	return nil
}
