package model

import (
	"time"

	m "github.com/CognyHub/evmd-golib/modelbase"
)

// DeviceDataLogReport type
type DeviceDataLogReport struct {
	ID               int          `db:"id"`
	TenantID         int          `db:"tenant_id"`
	CreatedAt        time.Time    `db:"created_at"`
	UpdatedAt        time.Time    `db:"updated_at"`
	IsActive         bool         `db:"is_active"`
	IsDeleted        bool         `db:"is_deleted"`
	LogStatusName    string       `db:"log_status_name"`
	LogError         string       `db:"log_error"`
	DeviceDataID     string       `db:"device_data_id"`
	DeviceCreatedAt  time.Time    `db:"device_created_at"`
	TableName        string       `db:"table_name"`
	PK               string       `db:"pk"`
	SFID             m.NullString `db:"sf_id"`
	ActionType       string       `db:"action_type"`
	ExternalID       string       `db:"external_id"`
	DeviceID         string       `db:"device_id"`
	UserID           string       `db:"user_id"`
	GroupID          string       `db:"group_id"`
	OriginalJSONData m.JSONB      `db:"original_json_data"`
	BrewedJSONData   m.JSONB      `db:"brewed_json_data"`
	ExecID           int          `db:"execution_id"`
	ExecStatusName   string       `db:"execution_status_name"`
	FaktoryID        string       `db:"execution_job_faktory_id"`
	Reported         bool         `db:"reported"`
}
