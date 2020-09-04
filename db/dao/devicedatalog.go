package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
)

//ArchiveDeviceDataLog func
func ArchiveDeviceDataLog(conn *sqlx.DB, tid int) error {
	query := `SELECT itgr.fn_archive_device_data_logs($1);`

	if _, err := conn.Exec(query, tid); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
