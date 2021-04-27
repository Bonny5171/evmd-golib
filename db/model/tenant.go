package model

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

//Tenant type
type Tenant struct {
	ID                 int            `db:"id"`
	CompanyID          string         `db:"company_id"`
	Name               string         `db:"name"`
	OrgID              string         `db:"org_id"`
	OrgType            string         `db:"organization_type"`
	CustomDomain       string         `db:"custom_domain"`
	IsSandbox          bool           `db:"is_sandbox"`
	IsActive           bool           `db:"is_active"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
	IsDeleted          bool           `db:"is_deleted"`
	DeletedAt          pq.NullTime    `db:"deleted_at"`
	LastModifiedByID   sql.NullString `db:"last_modified_by_id"`
	IsCloned           bool           `db:"is_cloned"`
	SfClientID         string         `db:"sf_client_id"`
	SfClientSecret     string         `db:"sf_client_secret"`
	SfCallbackTokenUrl string         `db:"sf_callback_token_url"`
}

type TenantCloneTable struct {
	TemplateTenantID string `db:"table_schema"`
	TableName        string `db:"table_name"`
}
