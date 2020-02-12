package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

// GetSchemas return all schemas from tenantID
func GetSchemas(conn *sqlx.DB, tid int) (s model.Schemas, err error) {
	const query = "SELECT name, description, type FROM itgr.schema WHERE tenantid = $1;"

	err = conn.QueryRowx(query, tid).StructScan(&s)
	if err != nil {
		return nil, db.WrapError(err, "conn.QueryRowx()")
	}

	return s, nil
}

// GetSchema return schema details searched by id
func GetSchema(conn *sqlx.DB, sid int) (s model.Schema, err error) {
	const query = "SELECT name, description, type FROM itgr.schema WHERE id = $1 LIMIT 1;"

	err = conn.QueryRowx(query, sid).StructScan(&s)
	if err != nil {
		return s, db.WrapError(err, "conn.QueryRowx()")
	}

	return s, nil
}

// GetSchemaByName return schema details searched by name
func GetSchemaByName(conn *sqlx.DB, tid int, name string) (s model.Schema, err error) {
	const query = "SELECT id, tenant_id, name, description, type FROM itgr.schema WHERE tenant_id = $1 AND name = $2 LIMIT 1;"

	err = conn.QueryRowx(query, tid, name).StructScan(&s)
	if err != nil {
		return s, db.WrapError(err, "conn.QueryRowx()")
	}

	return s, nil
}
