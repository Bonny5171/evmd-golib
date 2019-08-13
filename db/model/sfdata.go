package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

type SFData struct {
	ID           int64           `db:"id"`
	TenantID     int             `db:"tenant_id"`
	ExecutionID  int64           `db:"execution_id"`
	RecordTypeID int             `db:"record_type_id"`
	StatusID     int16           `db:"status_id"`
	ObjectName   string          `db:"sf_object_name"`
	ObjectID     int64           `db:"sf_object_id"`
	DocID        sql.NullString  `db:"doc_id"`
	DocRecord    json.RawMessage `db:"doc_record"`
	DocMetaData  json.RawMessage `db:"doc_meta_data"` // In DB is JSONB
	IsActive     bool            `db:"is_active"`
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at"`
	IsDeleted    bool            `db:"is_deleted"`
	DeletedAt    pq.NullTime     `db:"deleted_at"`
}
