package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
)

func GetCompleteOrgID(conn *sqlx.DB, pOrgID string) (orgID string, err error) {
	const query = `
		SELECT org_id
		FROM public.tenant
		WHERE LOWER(LEFT(org_id, 15)) = LOWER(LEFT($1, 15)) AND is_deleted = FALSE
		ORDER BY id DESC
		LIMIT 1;`

	row := conn.QueryRowx(query, pOrgID)
	if e := row.StructScan(&orgID); e != nil {
		err = db.WrapError(e, "row.StructScan()")
		return
	}

	return
}
