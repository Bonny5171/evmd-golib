package dao

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/gopkgs/db/model"
)

type (
	Status     int
	StatusType int
)

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
)

const (
	EnumTypeStatusNil StatusType = iota
	EnumTypeStatusETL
	EnumTypeStatusExec
)

func (t Status) String() string {
	n := [...]string{"processing", "error", "success", "scheduled", "overrided", "processing", "error", "success", "pending", "warning"}
	if t < EnumStatusExecProcessing || t > EnumStatusEtlWarning {
		return ""
	}
	return n[t]
}

func (t StatusType) String() string {
	n := [...]string{"", "etl", "exec"}
	if t < EnumTypeStatusNil || t > EnumTypeStatusExec {
		return ""
	}
	return n[t]
}

// GetStatuses retorna a lista de status de processamento de uma determinada org (tenant_id)
func GetStatuses(conn *sqlx.DB, tenantId int, sType StatusType) (s model.Statuses, err error) {
	qb := strings.Builder{}
	var args []interface{}

	qb.WriteString("SELECT id, name, type FROM itgr.status WHERE tenant_id = $1 AND is_active = true AND is_deleted = false")
	args = append(args, tenantId)

	if sType != EnumTypeStatusNil {
		qb.WriteString(" AND type = $2")
		args = append(args, sType.String())
	}

	err = conn.Select(&s, qb.String(), args...)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Select()")
	}

	return s, nil
}
