package model

import (
	"database/sql"
)

type JobScheduler struct {
	ID                int64          `db:"id"`
	Name              string         `db:"name"`
	TenantID          int            `db:"tenant_id"`
	TenantName        string         `db:"tenant_name"`
	OrgID             string         `db:"org_id"`
	Queue             string         `db:"queue"`
	JobName           string         `db:"job_name"`
	Description       sql.NullString `db:"description"`
	Parameters        sql.NullString `db:"parameters"`
	Cron              string         `db:"cron"`
	Retry             int16          `db:"retry"`
	AllowsConcurrency bool           `db:"allows_concurrency"`
	AllowsSchedule    bool           `db:"allows_schedule"`
	ScheduleTime      int16          `db:"schedule_time"`
	DocMetaData       sql.NullString `db:"doc_meta_data"` // In DB is JSONB
	IsActive          bool           `db:"is_active"`
	IsDeleted         bool           `db:"is_deleted"`
}
