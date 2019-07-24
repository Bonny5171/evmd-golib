package model

import (
	"database/sql"
	"time"

	m "bitbucket.org/everymind/evmd-golib/modelbase"
)

type Device struct {
	ID  string `db:"device_id"`
	Qty int    `db:"qty"`
}

type DeviceTableField struct {
	ObjectID   int            `db:"sf_object_id"`
	ObjectName string         `db:"sf_object_name"`
	TableName  string         `db:"sfa_table_name"`
	Fields     m.SliceMapJSON `db:"from_to_fields"`
	PrimaryKey m.NullString   `db:"primary_key"`
	ExternalID m.NullString   `db:"external_id"`
}

type DeviceData struct {
	ID              string         `db:"id"`
	TenantID        int            `db:"tenant_id"`
	SchemaName      string         `db:"schema_name"`
	TableName       string         `db:"table_name"`
	ObjectID        int            `db:"sf_object_id"`
	ObjectName      string         `db:"sf_object_name"`
	UserID          sql.NullString `db:"user_id"`
	PK              sql.NullString `db:"pk"`
	SfID            sql.NullString `db:"sf_id"`
	JSONData        m.JSONB        `db:"json_data"`
	BrewedJSONData  m.JSONB
	AppID           sql.NullString `db:"app_id"`
	DeviceID        sql.NullString `db:"device_id"`
	DeviceCreatedAt time.Time      `db:"device_created_at"`
}
