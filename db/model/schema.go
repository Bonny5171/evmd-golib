package model

import "encoding/json"

//Schema type
type Schema struct {
	ID          int             `db:"id"`
	TenantID    int             `db:"tenant_id"`
	Name        string          `db:"name"`
	Type        string          `db:"type"`
	Description string          `db:"description"`
	DocMetaData json.RawMessage `db:"doc_meta_data"` // In DB is JSONB
}

//Schemas type
type Schemas []Schema
