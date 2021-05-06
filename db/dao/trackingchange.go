package dao

func SelectRebuildTables(conn *sqlx.DB, tid int) ([]string, error) {
	var tableName []string
	query := fmt.Sprintf("SELECT table_name FROM itgr.vw_tenant_clone WHERE table_name LIKE 'sfa_%' AND table_schema = tn_%03d", tid)

	if err := conn.Get(&tableName, query); err != nil {
		return "", db.WrapError(err, "conn.Get()")
	}

	return tableName, err
}

func CountTableRows(conn *sqlx.DB, tid int, tableName string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT count(*) FROM tn_%03d.%s", tid, tableName)

	if err := conn.Get(&count, query); err != nil {
		return "", db.WrapError(err, "conn.Get()")
	}

	return count, err
}

//RebuildTrackingChange func
func RebuildTrackingChange(conn *sqlx.DB, tid int, targetTable string) error {
	query := `SELECT sync.fn_rebuild_tracking_change($1, $2);`

	if _, err := conn.Exec(query, tid, targetTable); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
