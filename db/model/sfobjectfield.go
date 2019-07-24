package model

import (
	"database/sql"

	m "bitbucket.org/everymind/evmd-golib/modelbase"
)

type SFObjectField struct {
	ID          int            `db:"id"`
	TenantID    int            `db:"tenant_id"`
	SfObjectID  int            `db:"sf_object_id"`
	SfFieldName string         `db:"sf_field_name"`
	SfType      string         `db:"sf_type"`
	SfLength    sql.NullInt64  `db:"sf_length"`
	SfPrecision sql.NullInt64  `db:"sf_precision"`
	SfScale     sql.NullInt64  `db:"sf_scale"`
	SfaName     sql.NullString `db:"sfa_name"`
	SfaType     sql.NullString `db:"sfa_type"`
	RawCommand  sql.NullString `db:"raw_command"`
	DocMetaData m.JSONB        `db:"doc_meta_data"` // In DB is JSONB
}
