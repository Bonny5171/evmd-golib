package functions

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cast"
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
	if len(args) < 8 {
		return p, errors.New("wrong number of args")
	}

	var e error

	// parameter is int64
	if p.JobID, e = cast.ToInt64E(args[0]); e != nil {
		return p, fmt.Errorf("parameter %d of job payload isn't a number: %w", 1, e)
	}

	// parameter is string
	if p.JobName, e = cast.ToStringE(args[1]); e != nil {
		return p, fmt.Errorf("parameter %d of job payload isn't a string: %w", 2, e)
	}

	// parameter is int
	if p.TenantID, e = cast.ToIntE(args[2]); e != nil {
		return p, fmt.Errorf("parameter %d of job payload isn't a number: %w", 3, e)
	}

	// parameter is string
	if p.TenantName, e = cast.ToStringE(args[3]); e != nil {
		return p, fmt.Errorf("parameter %d of job payload isn't a string: %w", 4, e)
	}

	// parameter is string
	if p.StackName, e = cast.ToStringE(args[4]); e != nil {
		return p, fmt.Errorf("parameter %d of job payload isn't a string: %w", 5, e)
	}

	// parameter is bool
	if p.AllowsConcurrency, e = cast.ToBoolE(args[5]); e != nil {
		return p, fmt.Errorf("parameter %d of job payload isn't a boolean: %w", 6, e)
	}

	// parameter is bool
	if p.AllowsSchedule, e = cast.ToBoolE(args[6]); e != nil {
		return p, fmt.Errorf("parameter %d of job payload isn't a boolean: %w", 7, e)
	}

	// parameter is int
	if p.ScheduleTime, e = cast.ToIntE(args[7]); e != nil {
		return p, fmt.Errorf("parameter %d of job payload isn't a number: %w", 8, e)
	}

	if len(args) > 8 {
		// parameter is int
		fParams, e := cast.ToStringE(args[8])
		if e != nil {
			return p, fmt.Errorf("parameter %d of job payload isn't a JSON: %w", 9, e)
		}

		if e := json.Unmarshal([]byte(fParams), &p.Parameters); e != nil {
			return p, fmt.Errorf("json.Unmarshal(): %w", e)
		}
	}

	return p, nil
}
