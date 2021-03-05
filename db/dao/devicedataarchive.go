package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
)

//ArchiveDeviceData func
func ArchiveDeviceData(conn *sqlx.DB, tid int) error {
	query := `SELECT itgr.fn_archive_device_data($1);`

	if _, err := conn.Exec(query, tid); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
