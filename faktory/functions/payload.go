package functions

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Payload struct {
	JobID             int64
	JobName           string
	TenantID          int
	TenantName        string
	StackName         string
	AllowsConcurrency bool
	AllowsSchedule    bool
	ScheduleTime      int
	Parameters        map[string]interface{}
}

func ParsePayload(args ...interface{}) (p Payload, err error) {
	if len(args) <= 8 {
		return p, errors.Wrap(err, "wrong number of args")
	}

	var ok bool

	// parameter is int64
	fID, ok := args[0].(float64)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a number", 1)
	}
	p.JobID = int64(fID)

	// parameter is string
	p.JobName, ok = args[1].(string)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a string", 2)
	}

	// parameter is int
	fTeID, ok := args[2].(float64)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a number", 3)
	}
	p.TenantID = int(fTeID)

	// parameter is string
	p.TenantName, ok = args[3].(string)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a string", 4)
	}

	// parameter is string
	p.StackName, ok = args[4].(string)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a string", 5)
	}

	// parameter is bool
	p.AllowsConcurrency, ok = args[5].(bool)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a boolean", 6)
	}

	// parameter is bool
	p.AllowsSchedule, ok = args[6].(bool)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a boolean", 7)
	}

	// parameter is int
	fSchTime, ok := args[7].(float64)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a number", 8)
	}
	p.ScheduleTime = int(fSchTime)

	if len(args) > 8 {
		// parameter is int
		fParams, ok := args[8].(string)
		if !ok {
			return p, errors.Wrapf(err, "parameter %d of job payload isn't a json", 9)
		}

		if e := json.Unmarshal([]byte(fParams), &p.Parameters); e != nil {
			return p, errors.Wrap(e, "json.Unmarshal()")
		}
	}

	return p, nil
}
