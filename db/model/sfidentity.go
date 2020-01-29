package model

import (
	m "bitbucket.org/everymind/evmd-golib/v2/modelbase"
)

type SFIdentity struct {
	ID          int     `db:"id"`
	TenantID    int     `db:"tenant_id"`
	ExecutionID int64   `db:"execution_id"`
	Name        string  `db:"name"`
	DocOrg      m.JSONB `db:"doc_org"`
	DocObjects  m.JSONB `db:"doc_objects"`
	DocMetaData m.JSONB `db:"doc_meta_data"` // In DB is JSONB
}
