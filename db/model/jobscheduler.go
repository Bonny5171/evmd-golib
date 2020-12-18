package model

import (
	"database/sql"
)

//JobScheduler type
type JobScheduler struct {
	ID                int64          `db:"id"`
	OrgID             string         `db:"org_id"`
	TenantID          int            `db:"tenant_id"`
	TenantName        string         `db:"tenant_name"`
	StackID           int            `db:"stack_id"`
	JobName           string         `db:"job_name"`
	FunctionName      string         `db:"function_name"`
	Queue             string         `db:"queue"`
	Cron              sql.NullString `db:"cron"`
	Parameters        sql.NullString `db:"parameters"`
	Retry             int16          `db:"retry"`
	AllowsConcurrency bool           `db:"allows_concurrency"`
	AllowsSchedule    bool           `db:"allows_schedule"`
	ScheduleTime      int16          `db:"schedule_time"`
	Description       sql.NullString `db:"description"`
	IsActive          bool           `db:"is_active"`
	IsDeleted         bool           `db:"is_deleted"`
	AppEngineName     string         `db:"appengine_name"`
}
