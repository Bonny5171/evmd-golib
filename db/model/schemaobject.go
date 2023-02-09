package model

import (
	"database/sql"

	"github.com/lib/pq"

	m "github.com/CognyHub/evmd-golib/modelbase"
)

// SchemaObject type
type SchemaObject struct {
	ID          int            `db:"id"`
	TenantID    int            `db:"tenant_id"`
	SchemaID    int            `db:"schema_id"`
	ObjectID    sql.NullInt64  `db:"sf_object_id"`
	ObjectName  sql.NullString `db:"sf_object_name"`
	Sequence    int16          `db:"sequence"`
	DocFields   m.JSONB        `db:"doc_fields"`
	Filter      m.NullString   `db:"filter"`
	RawCommand  sql.NullString `db:"raw_command"`
	DocMetaData m.JSONB        `db:"doc_meta_data"` // In DB is JSONB
}

// SchemaObjects type
type SchemaObjects []SchemaObject

// SchemaObjectToProcess type
type SchemaObjectToProcess struct {
	ID                int            `db:"id"`
	TenantID          int            `db:"tenant_id"`
	TenantName        string         `db:"tenant_name"`
	SchemaID          int            `db:"schema_id"`
	SchemaName        string         `db:"schema_name"`
	Type              string         `db:"type"`
	APIType           string         `db:"api_type"`
	ObjectID          sql.NullInt64  `db:"sf_object_id"`
	ObjectName        sql.NullString `db:"sf_object_name"`
	Sequence          int16          `db:"sequence"`
	DocFields         m.JSONB        `db:"doc_fields"`
	Filter            m.NullString   `db:"filter"`
	RawCommand        sql.NullString `db:"raw_command"`
	LastModifiedDate  pq.NullTime    `db:"sf_last_modified_date"`
	Layoutable        bool           `db:"layoutable"`
	CompactLayoutable bool           `db:"compactlayoutable"`
	Listviewable      bool           `db:"listviewable"`
	SfPkName          string         `db:"sf_pk_name"`
	SfaPks            m.JSONB        `db:"sfa_pks"`
}

// SchemaObjectToProcesses type
type SchemaObjectToProcesses []SchemaObjectToProcess
