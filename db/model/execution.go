package model

import (
	"database/sql"
	"time"

	"github.com/lib/pq"

	m "bitbucket.org/everymind/evmd-golib/v2/modelbase"
)

type Execution struct {
	ID               int64          `db:"id"`
	DocMetaData      m.JSONB        `db:"doc_meta_data"` // In DB is JSONB
	JobFaktoryID     sql.NullString `db:"job_faktory_id"`
	JobSchedulerID   int64          `db:"job_scheduler_id"`
	JobSchedulerName string         `db:"job_scheduler_name"`
	SchemaID         sql.NullInt64  `db:"schema_id"`
	StatusID         int16          `db:"status_id"`
	Status           sql.NullString `db:"status_name"`
	TenantID         int            `db:"tenant_id"`
	IsActive         bool           `db:"is_active"`
	IsDeleted        bool           `db:"is_deleted"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
	DeletedAt        pq.NullTime    `db:"deleted_at"`
}
