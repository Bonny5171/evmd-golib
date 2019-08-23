package dao

import (
	"bitbucket.org/everymind/evmd-golib/db"
	"github.com/jmoiron/sqlx"
)

func GetTenant(conn *sqlx.DB, orgID string) (tid int, err error) {
	const query = `
		SELECT id
		  FROM public.tenant
		 WHERE org_id = $1
		 LIMIT 1;`

	row := conn.QueryRow(query, orgID)

	if e := row.Scan(&tid); e != nil {
		err = db.WrapError(e, "row.Scan()")
		return
	}

	return
}
