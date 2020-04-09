package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
)

//CleanGhostJobs func
func CleanGhostJobs(conn *sqlx.DB) (err error) {
	query := `SELECT itgr.fn_set_jobs_time_out();`

	if _, err = conn.Query(query); err != nil {
		return db.WrapError(err, "conn.Query()")
	}
	return
}
