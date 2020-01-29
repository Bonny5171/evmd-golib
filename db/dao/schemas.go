package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/v2/db"
	"bitbucket.org/everymind/evmd-golib/v2/db/model"
)

// GetSchemas
func GetSchemas(conn *sqlx.DB, tid, sid int) (s model.Schemas, err error) {
	const query = "SELECT name, description, type FROM itgr.schema WHERE tenant_id = $1 AND id = $2 LIMIT 1;"

	err = conn.QueryRowx(query, tid, sid).StructScan(&s)
	if err != nil {
		return nil, db.WrapError(err, "conn.QueryRowx()")
	}

	return s, nil
}
