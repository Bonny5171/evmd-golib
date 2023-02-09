package dao

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/CognyHub/evmd-golib/db"
)

// GetCompleteOrgID func
func GetCompleteOrgID(conn *sqlx.DB, orgID string) (cOrgID string, err error) {
	log.Printf("GetCompleteOrgID: %v", orgID)
	const query = `
		SELECT org_id
		FROM public.tenant
		WHERE LEFT(org_id, 15) = LEFT($1, 15)
		AND is_active = TRUE
		AND is_deleted = FALSE
		LIMIT 1;`

	row := conn.QueryRow(query, orgID)

	if e := row.Scan(&cOrgID); e != nil {
		err = db.WrapError(e, "row.Scan()")
		return
	}

	return
}
