package model

import (
	"database/sql"
	"time"

	m "github.com/CognyHub/evmd-golib/modelbase"
)

// Device type
type Device struct {
	ID        string       `db:"device_id"`
	GroupID   m.NullString `db:"group_id"`
	TableName string       `db:"table_name"`
	Qty       int          `db:"qty"`
}

// DeviceTableField type
type DeviceTableField struct {
	ObjectID   int            `db:"sf_object_id"`
	ObjectName string         `db:"sf_object_name"`
	TableName  string         `db:"sfa_table_name"`
	Fields     m.SliceMapJSON `db:"from_to_fields"`
	PrimaryKey m.NullString   `db:"primary_key"`
	ExternalID m.NullString   `db:"external_id"`
	SfaPks     m.JSONB        `db:"sfa_pks"`
}

// DeviceData type
type DeviceData struct {
	ID              string         `db:"id"`
	TenantID        int            `db:"tenant_id"`
	SchemaName      string         `db:"schema_name"`
	TableName       string         `db:"table_name"`
	ObjectID        int            `db:"sf_object_id"`
	ObjectName      string         `db:"sf_object_name"`
	UserID          sql.NullString `db:"user_id"`
	PK              sql.NullString `db:"pk"`
	ExternalID      sql.NullString `db:"external_id"`
	SfID            sql.NullString `db:"sf_id"`
	ActionType      sql.NullString `db:"action_type"`
	JSONData        m.JSONB        `db:"json_data"`
	BrewedJSONData  m.JSONB
	AppID           sql.NullString `db:"app_id"`
	DeviceID        sql.NullString `db:"device_id"`
	DeviceCreatedAt time.Time      `db:"device_created_at"`
	GroupID         sql.NullString `db:"group_id"`
	Sequential      sql.NullInt64  `db:"sequential"`
	Try             int            `db:"try"`
	IsActive        bool           `db:"is_active"`
	IsDeleted       bool           `db:"is_deleted"`
}
