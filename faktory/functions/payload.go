package functions

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Payload struct {
	ID                int64
	TenantID          int
	TenantName        string
	StackName         string
	AllowsConcurrency bool
	AllowsSchedule    bool
	ScheduleTime      int
	Parameters        map[string]interface{}
}

func ParsePayload(args ...interface{}) (p Payload, err error) {
	if len(args) < 4 {
		return p, errors.Wrap(err, "number of args is invalid to this job")
	}

	var ok bool

	// parameter is int64
	fID, ok := args[0].(float64)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a number", 1)
	}
	p.ID = int64(fID)

	// parameter is int
	fTeID, ok := args[1].(float64)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a number", 2)
	}
	p.TenantID = int(fTeID)

	// parameter is string
	p.TenantName, ok = args[2].(string)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a string", 3)
	}

	// parameter is string
	p.StackName, ok = args[3].(string)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a string", 4)
	}

	// parameter is bool
	p.AllowsConcurrency, ok = args[4].(bool)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a boolean", 5)
	}

	// parameter is bool
	p.AllowsSchedule, ok = args[5].(bool)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a boolean", 6)
	}

	// parameter is int
	fSchTime, ok := args[6].(float64)
	if !ok {
		return p, errors.Wrapf(err, "parameter %d of job payload isn't a number", 7)
	}
	p.ScheduleTime = int(fSchTime)

	if len(args) > 7 {
		// parameter is int
		fParams, ok := args[7].(string)
		if !ok {
			return p, errors.Wrapf(err, "parameter %d of job payload isn't a json", 8)
		}

		if e := json.Unmarshal([]byte(fParams), &p.Parameters); e != nil {
			return p, errors.Wrap(e, "json.Unmarshal()")
		}
	}

	return p, nil
}
