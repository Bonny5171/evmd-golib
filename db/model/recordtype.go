package model

import (
	"time"

	"github.com/lib/pq"
)

//RecordType type
type RecordType struct {
	ID           int         `db:"id"`
	TenantID     int         `db:"tenant_id"`
	IsActive     bool        `db:"is_active"`
	CreatedAt    time.Time   `db:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at"`
	IsDeleted    bool        `db:"is_deleted"`
	DeletedAt    pq.NullTime `db:"deleted_at"`
	Name         string      `db:"name"`
	DeveloperRef string      `db:"developer_ref"`
	IsSystemType bool        `db:"is_system_type"`
}
