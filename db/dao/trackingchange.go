package dao

import (
	"fmt"

	"bitbucket.org/everymind/evmd-golib/db"
	"github.com/jmoiron/sqlx"
)

func GetCountParameter(conn *sqlx.DB, tid int) (int, error) {
	var count int
	query := `SELECT value FROM itgr.parameter WHERE tenant_id=$1 AND name='REBUILD_TRACKING_CHANGE_COUNT'`

	if err := conn.Get(&count, query); err != nil {
		return 0, db.WrapError(err, "conn.Get()")
	}

	return count, nil
}

func SelectRebuildTables(conn *sqlx.DB, tid int) ([]string, error) {
	var tableName []string
	query := fmt.Sprintf("SELECT table_name FROM itgr.vw_tenant_clone WHERE table_name LIKE 'sfa_%' AND table_schema = tn_%03d", tid)

	if err := conn.Get(&tableName, query); err != nil {
		return nil, db.WrapError(err, "conn.Get()")
	}

	return tableName, nil
}

func CountTableRows(conn *sqlx.DB, tid int, tableName string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT count(*) FROM tn_%03d.%s", tid, tableName)

	if err := conn.Get(&count, query); err != nil {
		return 0, db.WrapError(err, "conn.Get()")
	}

	return count, nil
}

//RebuildTrackingChange func
func RebuildTrackingChange(conn *sqlx.DB, tid int, targetTable string) error {
	query := `SELECT sync.fn_rebuild_tracking_change($1, $2);`

	if _, err := conn.Exec(query, tid, targetTable); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
