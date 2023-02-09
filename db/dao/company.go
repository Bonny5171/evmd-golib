package dao

import (
	"github.com/jmoiron/sqlx"

	"github.com/CognyHub/evmd-golib/db"
)

// CheckEnvironmentAvailability func
func CheckEnvironmentAvailability(conn *sqlx.DB, companyID string, isSandbox bool) (result bool, err error) {
	var prodQtd, sndbQtd int

	query := `SELECT production_remaining, sandbox_remaining FROM public.vw_environment_count WHERE company_id = $1;`

	row := conn.QueryRow(query, companyID)

	if err := row.Scan(&prodQtd, &sndbQtd); err != nil {
		return false, db.WrapError(err, "row.Scan()")
	}

	if isSandbox {
		if sndbQtd > 0 {
			result = true
		}
	} else {
		if prodQtd > 0 {
			result = true
		}
	}

	return result, nil
}
