package dao

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

type (
	//Status type
	Status int
	//StatusType types
	StatusType int
)

//EnumStatus types
const (
	EnumStatusExecProcessing Status = iota
	EnumStatusExecError
	EnumStatusExecSuccess
	EnumStatusExecScheduled
	EnumStatusExecOverrided
	EnumStatusEtlProcessing
	EnumStatusEtlError
	EnumStatusEtlSuccess
	EnumStatusEtlPending
	EnumStatusEtlWarning
	EnumStatusUpsertPending
	EnumStatusUpsertSuccess
	EnumStatusUpsertError
)

// EnumTyoe types
const (
	EnumTypeStatusNil StatusType = iota
	EnumTypeStatusETL
	EnumTypeStatusExec
	EnumTypeStatusUpsert
)

func (t Status) String() string {
	n := [...]string{"processing", "error", "success", "scheduled", "overrided", "processing", "error", "success", "pending", "warning", "pending", "success", "error"}
	if t < EnumStatusExecProcessing || t > EnumStatusUpsertError {
		return ""
	}
	return n[t]
}

func (t StatusType) String() string {
	n := [...]string{"", "etl", "exec", "upsert"}
	if t < EnumTypeStatusNil || t > EnumTypeStatusUpsert {
		return ""
	}
	return n[t]
}

// GetStatuses retorna a lista de status de processamento de uma determinada org (tenant_id)
func GetStatuses(conn *sqlx.DB, tenantID int, sType StatusType) (s model.Statuses, err error) {
	qb := strings.Builder{}
	var args []interface{}

	qb.WriteString("SELECT id, name, type FROM itgr.status WHERE tenant_id = $1 AND is_active = true AND is_deleted = false")
	args = append(args, tenantID)

	if sType != EnumTypeStatusNil {
		qb.WriteString(" AND type = $2")
		args = append(args, sType.String())
	}

	err = conn.Select(&s, qb.String(), args...)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return s, nil
}
