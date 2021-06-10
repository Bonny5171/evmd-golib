package dao

import (
	"fmt"
	"strings"

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

func SelectRebuildTables(conn *sqlx.DB, tid int, excludedTables []string) ([]string, error) {
	var tableName []string
	// query := "SELECT table_name FROM itgr.vw_tenant_clone WHERE table_name LIKE 'sfa_%' AND table_schema = ?"

	var query strings.Builder

	query.WriteString(`SELECT distinct t.table_name AS table_name
						FROM sync.vw_tab t 
						WHERE t.column_name = ('tracking_change_id') 
							AND t.table_name like 'sfa_%'
							AND t.table_schema = public.fn_schema_name(`)
	query.WriteString(fmt.Sprintf("%d", tid))
	query.WriteString(`) AND t.table_name NOT IN (`)
	for index, table := range excludedTables {
		query.WriteString(fmt.Sprintf("'%s'", table))
		if index < len(excludedTables) - 1 {
			query.WriteString(`,`)
		}
	}
	query.WriteString(`) ORDER BY 1`)

	if err := conn.Get(&tableName, query.String()); err != nil {
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
