package dao

import (
	"strings"

	"bitbucket.org/everymind/evmd-golib/db"
	"github.com/jmoiron/sqlx"
)

func ExecSFEtlData(conn *sqlx.DB, execID int64, tenantID int, objID int64, reprocessAll bool) error {
	query := "SELECT itgr.sf_etl_data($1, $2, $3, $4);"

	if _, err := conn.Exec(query, execID, tenantID, objID, reprocessAll); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func ExecSfEtlJsonData(conn *sqlx.DB, execID int64, tenantID, recordTypeID int) error {
	query := "SELECT itgr.sf_etl_data_json($1, $2, $3);"

	if _, err := conn.Exec(query, execID, tenantID, recordTypeID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func ExecSFEtlShareData(conn *sqlx.DB, execID int64, tenantID int, userID string) error {
	query := "SELECT itgr.sf_etl_data_share($1, $2, $3);"

	if _, err := conn.Exec(query, execID, tenantID, userID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func ExecSFEtlSyncData(conn *sqlx.DB, execID int64, tenantID int, objID int64) error {
	query := "SELECT itgr.sf_etl_data_sync($1, $2, $3);"

	if _, err := conn.Exec(query, execID, tenantID, objID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func ExecSFCreateAllTables(conn *sqlx.DB) error {
	query := "DO $$ BEGIN PERFORM itgr.sf_create_all_tables(); END $$;"

	if _, err := conn.Exec(query); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func ExecSFPurgePublicSFTables(conn *sqlx.DB, tenantID int) error {
	query := "SELECT itgr.sf_purge_sf_tables($1);"

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func ExecSFPurgePublicSFShare(conn *sqlx.DB, tenantID int) error {
	query := "SELECT itgr.sf_purge_sf_share($1);"

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func ExecSFCheckJobsExection(conn *sqlx.DB, tenantID int, jobName, statusName string) (result bool, err error) {
	query := "SELECT itgr.fn_check_jobs($1, $2, $3);"

	row := conn.QueryRow(query, tenantID, jobName, statusName)

	if err := row.Scan(&result); err != nil {
		return false, db.WrapError(err, "row.Scan()")
	}

	return result, nil
}

func ExecSFAfterEtl(conn *sqlx.DB, tenantID int) error {
	return ExecSFAExecEtls(conn, tenantID, "")
}

func ExecSFAExecEtls(conn *sqlx.DB, tenantID int, tableName string) error {
	params := make([]interface{}, 0)
	params = append(params, tenantID)

	sb := strings.Builder{}
	sb.WriteString("SELECT itgr.fn_exec_etls($1")
	if len(tableName) > 0 {
		sb.WriteString(", $2")
		params = append(params, tableName)
	}
	sb.WriteString(");")

	if _, err := conn.Exec(sb.String(), params...); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
