package dao

import (
	"github.com/jmoiron/sqlx"

	"github.com/CognyHub/evmd-golib/db"
	"github.com/CognyHub/evmd-golib/db/model"
)

// GetRecordType func
func GetRecordType(conn *sqlx.DB, tid int, devRef string) (s model.RecordType, err error) {
	const query = `
		SELECT id, tenant_id, is_active, created_at, updated_at, is_deleted, deleted_at, "name", developer_ref, is_system_type 
		  FROM public.record_type 
		 WHERE tenant_id = $1 AND is_active = TRUE AND is_deleted = FALSE AND developer_ref = $2 
		 LIMIT 1;`

	if err = conn.Get(&s, query, tid, devRef); err != nil {
		err = db.WrapError(err, "conn.Get()")
		return
	}

	return
}
