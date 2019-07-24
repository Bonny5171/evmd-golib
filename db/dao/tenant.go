package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func GetTenant(conn *sqlx.DB, orgID string) (tid int, err error) {
	const query = `
		SELECT id
		  FROM public.tenant
		 WHERE org_id = $1
		 LIMIT 1;`

	row := conn.QueryRow(query, orgID)

	if e := row.Scan(&tid); e != nil {
		err = errors.Wrap(e, "row.Scan()")
		return
	}

	return
}
