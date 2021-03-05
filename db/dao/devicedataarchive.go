package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
)

//ArchiveDeviceData func
func ArchiveDeviceData(conn *sqlx.DB, tid int, days int) error {
	query := `SELECT itgr.fn_archive_device_data($1, $2);`

	if _, err := conn.Exec(query, tid, days); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
