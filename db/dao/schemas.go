package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/gopkgs/db/model"
)

// GetSchemas
func GetSchemas(conn *sqlx.DB, tid, sid int) (s model.Schemas, err error) {
	const query = "SELECT name, description, type FROM itgr.schema WHERE tenant_id = ? AND id = ? LIMIT 1;"

	err = conn.QueryRowx(query, tid, sid).StructScan(&s)
	if err != nil {
		return nil, errors.Wrap(err, "dbq.(*sqlx.DB).QueryRowx()")
	}

	return s, nil
}
