package model

import (
	"database/sql"
)

type Parameter struct {
	ID          int            `db:"id"`
	TenantID    int            `db:"tenant_id"`
	OrgID       string         `db:"org_id"`
	Name        string         `db:"name"`
	Value       string         `db:"value"`
	Type        string         `db:"type"`
	Description sql.NullString `db:"description"`
}

type Parameters []Parameter

func (p *Parameters) ByName(name string) (result string) {
	for _, item := range *(p) {
		if item.Name == name {
			result = item.Value
			break
		}
	}

	return
}
