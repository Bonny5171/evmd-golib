package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
)

//GetCompleteOrgID func
func GetCompleteOrgID(conn *sqlx.DB, orgID string) (org string, err error) {
	const query = `
		SELECT org_id
		FROM public.tenant
		WHERE LEFT(org_id, 15) = LEFT($1, 15)
		AND is_active = TRUE
		AND is_deleted = FALSE
		LIMIT 1;`

	row := conn.QueryRow(query, orgID)

	if e := row.Scan(&org); e != nil {
		err = db.WrapError(e, "row.Scan()")
		return
	}

	return
}
