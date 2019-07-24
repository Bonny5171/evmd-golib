package model

import (
	"database/sql"
	"time"

	m "bitbucket.org/everymind/gopkgs/modelbase"
	"github.com/lib/pq"
)

type Execution struct {
	ID                 int64          `db:"id"`
	JobSchedulerID     int64          `db:"job_scheduler_id"`
	TenantID           int            `db:"tenant_id"`
	SchemaID           sql.NullInt64  `db:"schema_id"`
	StatusID           int16          `db:"status_id"`
	JobFaktoryID       sql.NullString `db:"job_faktory_id"`
	SfLastModifiedDate pq.NullTime    `db:"sf_last_modified_date"`
	DocMetaData        m.JSONB        `db:"doc_meta_data"` // In DB is JSONB
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
	IsActive           bool           `db:"is_active"`
	IsDeleted          bool           `db:"is_deleted"`
}
