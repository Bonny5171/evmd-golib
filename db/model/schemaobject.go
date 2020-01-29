package model

import (
	"database/sql"

	"github.com/lib/pq"

	m "bitbucket.org/everymind/evmd-golib/v2/modelbase"
)

type SchemaObject struct {
	ID          int            `db:"id"`
	TenantID    int            `db:"tenant_id"`
	SchemaID    int            `db:"schema_id"`
	ObjectID    int            `db:"sf_object_id"`
	Sequence    int16          `db:"sequence"`
	DocFields   m.JSONB        `db:"doc_fields"`
	Filter      m.NullString   `db:"filter"`
	RawCommand  sql.NullString `db:"raw_command"`
	DocMetaData m.JSONB        `db:"doc_meta_data"` // In DB is JSONB
}

type SchemaObjects []SchemaObject

type SchemaObjectToProcess struct {
	ID                int            `db:"id"`
	TenantID          int            `db:"tenant_id"`
	TenantName        string         `db:"tenant_name"`
	SchemaID          int            `db:"schema_id"`
	SchemaName        string         `db:"schema_name"`
	Type              string         `db:"type"`
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
}

type SchemaObjectToProcesses []SchemaObjectToProcess
