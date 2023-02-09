package dao

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/CognyHub/evmd-golib/db"
	"github.com/CognyHub/evmd-golib/logger"
)

// ExecSFEtlData func
func ExecSFEtlData(conn *sqlx.DB, execID int64, tenantID int, objID int64, reprocessAll bool) error {
	query := "SELECT itgr.sf_etl_data($1, $2, $3, $4);"

	if _, err := conn.Exec(query, execID, tenantID, objID, reprocessAll); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSfEtlJsonData func
func ExecSfEtlJsonData(conn *sqlx.DB, execID int64, tenantID, recordTypeID int) error {
	query := "SELECT itgr.sf_etl_data_json($1, $2, $3);"

	if _, err := conn.Exec(query, execID, tenantID, recordTypeID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFEtlShareData func
func ExecSFEtlShareData(conn *sqlx.DB, execID int64, tenantID int, userID string) error {
	query := "SELECT itgr.sf_etl_data_share($1, $2, $3);"

	if _, err := conn.Exec(query, execID, tenantID, userID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFEtlSyncData func
func ExecSFEtlSyncData(conn *sqlx.DB, execID int64, tenantID int, objID int64) error {
	query := "SELECT itgr.sf_etl_data_sync($1, $2, $3);"

	if _, err := conn.Exec(query, execID, tenantID, objID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFCreateAllTables func
func ExecSFCreateAllTables(conn *sqlx.DB, tenantID int) error {
	query := "SELECT itgr.sf_create_all_tables($1);"

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFPurgePublicSFTables func
func ExecSFPurgePublicSFTables(conn *sqlx.DB, tenantID int) error {
	query := "SELECT itgr.sf_purge_sf_tables($1);"

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFPurgePublicSFShare func
func ExecSFPurgePublicSFShare(conn *sqlx.DB, tenantID int) error {
	query := "SELECT itgr.sf_purge_sf_share($1);"

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFCheckJobsExecution func
func ExecSFCheckJobsExecution(conn *sqlx.DB, tenantID int, jobID int64, statusName string) (result bool, err error) {
	query := "SELECT itgr.fn_check_jobs($1, $2, $3);"

	row := conn.QueryRow(query, tenantID, jobID, statusName)

	if err := row.Scan(&result); err != nil {
		return false, db.WrapError(err, "row.Scan()")
	}

	return result, nil
}

// ExecSFAfterEtl func
func ExecSFAfterEtl(conn *sqlx.DB, tenantID int) error {
	return ExecSFAExecEtls(conn, tenantID, "")
}

// ExecSFAfterEtlSerial func
func ExecSFAfterEtlSerial(conn *sqlx.DB, tentantID int) error {
	query := fmt.Sprintf("SELECT COALESCE(etl.output_table_name, id) FROM itgr.fn_sf_etl_config_tables(%d) etl ORDER BY etl.order_by;", tentantID)

	var etls []string
	var etl string

	rows, err := conn.Queryx(query)
	if err != nil {
		return db.WrapError(err, "conn.Exec()")
	}
	for rows.Next() {
		err := rows.Scan(&etl)
		if err != nil {
			return db.WrapError(err, "conn.Exec()")
		}
		etls = append(etls, etl)
	}

	for _, etlID := range etls {
		logger.Debugf("Transaction ETL Begin: %s", etlID)
		tx, err := conn.Begin()
		if err != nil {
			return db.WrapError(err, "conn.Exec()")
		}
		logger.Debugf("Set Isolation Mode for ETL: %s", etlID)
		_, err = tx.Exec(`set transaction isolation level read uncommitted`)
		if err != nil {
			tx.Rollback()
			return db.WrapError(err, "conn.Exec()")
		}
		etlQuery := fmt.Sprintf("SELECT itgr.fn_exec_etls(%d, '%s');", tentantID, etlID)
		logger.Debugf("Execute of query: %s", etlQuery)
		_, err = tx.Query(etlQuery)
		if err != nil {
			tx.Rollback()
			return db.WrapError(err, "conn.Exec()")
		}
		logger.Debugf("Finish of Execute query: %s", etlQuery)
		tx.Commit()
		logger.Debugf("Commit of Execute query: %s", etlQuery)
	}

	return nil
}

// ExecSFAExecEtls func
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

// ExecSFAExecEtlsTx func
func ExecSFAExecEtlsTx(conn *sqlx.Tx, tenantID int, tableName string) error {
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

// ExecSFCreateJobScheduler func
func ExecSFCreateJobScheduler(conn *sqlx.DB, tenantID int, templateTenantID int) error {
	if templateTenantID > 0 {
		query := "SELECT public.fn_create_job_scheduler($1, $2);"
		if _, err := conn.Exec(query, tenantID, templateTenantID); err != nil {
			return db.WrapError(err, "conn.Exec()")
		}
	} else {
		query := "SELECT public.fn_create_job_scheduler($1);"
		if _, err := conn.Exec(query, tenantID); err != nil {
			return db.WrapError(err, "conn.Exec()")
		}
	}

	return nil
}

// ExecSFCreateJobSchedulerTx func
func ExecSFCreateJobSchedulerTx(conn *sqlx.Tx, tenantID int, templateTenantID int) error {
	if templateTenantID > 0 {
		query := "SELECT public.fn_create_job_scheduler($1, $2);"
		if _, err := conn.Exec(query, tenantID, templateTenantID); err != nil {
			return db.WrapError(err, "conn.Exec()")
		}
	} else {
		query := "SELECT public.fn_create_job_scheduler($1);"
		if _, err := conn.Exec(query, tenantID); err != nil {
			return db.WrapError(err, "conn.Exec()")
		}
	}

	return nil
}

// ExecSFSchemaBuild func
func ExecSFSchemaBuild(conn *sqlx.DB, tenantID int) error {
	query := "SELECT public.fn_schema_build($1);"

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFSchemaBuildTx func
func ExecSFSchemaBuildTx(conn *sqlx.Tx, tenantID int) error {
	query := "SELECT public.fn_schema_build($1);"

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFSchemaCreate func
func ExecSFSchemaCreate(conn *sqlx.DB, tenantID, templateTenantID int, tenantName, orgID string) error {
	query := "SELECT public.fn_schema_create($1, $2, $3, $4);"

	if _, err := conn.Exec(query, tenantName, orgID, tenantID, templateTenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFSchemaCreateTx func
func ExecSFSchemaCreateTx(conn *sqlx.Tx, tenantID, templateTenantID int, tenantName, orgID string) error {
	query := "SELECT public.fn_schema_create($1, $2, $3, $4);"

	if _, err := conn.Exec(query, tenantName, orgID, tenantID, templateTenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFSchemaRemove func
func ExecSFSchemaRemove(conn *sqlx.DB, tenantID int) error {
	query := "SELECT public.fn_schema_remove($1);"

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFSchemaRemoveTx func
func ExecSFSchemaRemoveTx(conn *sqlx.Tx, tenantID int) error {
	query := "SELECT public.fn_schema_remove($1);"

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFDataCreateFromTemplates func
func ExecSFDataCreateFromTemplates(conn *sqlx.DB, tenantID int) error {
	query := "SELECT public.fn_data_create_from_templates($1);"

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFDataCreateFromTemplatesTx func
func ExecSFDataCreateFromTemplatesTx(conn *sqlx.Tx, tenantID int) error {
	query := "SELECT public.fn_data_create_from_templates($1);"

	if _, err := conn.Exec(query, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFSchemaCreateCloned func
func ExecSFSchemaCreateCloned(conn *sqlx.DB, tenantName, orgID string, tenantID, templateTenantID int) error {
	query := "SELECT public.fn_schema_create_cloned($1, $2, $3, $4)"

	if _, err := conn.Exec(query, tenantName, orgID, tenantID, templateTenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFSchemaCreateClonedTx func
func ExecSFSchemaCreateClonedTx(conn *sqlx.Tx, tenantName, orgID string, tenantID, templateTenantID int) error {
	query := "SELECT public.fn_schema_create_cloned($1, $2, $3, $4)"

	if _, err := conn.Exec(query, tenantName, orgID, tenantID, templateTenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFTablesCloneCreate func
func ExecSFTablesCloneCreate(conn *sqlx.DB, tenantID, templateTenantID int) error {
	query := "SELECT public.fn_sf_tables_clone_create($1, $2)"

	if _, err := conn.Exec(query, tenantID, templateTenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFTablesCloneCreateTx func
func ExecSFTablesCloneCreateTx(conn *sqlx.Tx, tenantID, templateTenantID int) error {
	query := "SELECT public.fn_sf_tables_clone_create($1, $2)"

	if _, err := conn.Exec(query, tenantID, templateTenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFTablesCloneData func
func ExecSFTablesCloneData(conn *sqlx.DB, tenantID, templateTenantID int) error {
	query := "SELECT public.fn_sf_tables_clone_data($1, $2)"

	if _, err := conn.Exec(query, tenantID, templateTenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFTablesCloneDataTx func
func ExecSFTablesCloneDataTx(conn *sqlx.Tx, tenantID, templateTenantID int) error {
	query := "SELECT public.fn_sf_tables_clone_data($1, $2)"

	if _, err := conn.Exec(query, tenantID, templateTenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFForeignSchemaClone func
func ExecSFForeignSchemaClone(conn *sqlx.DB, tenantID, templateTenantID int) error {
	query := "SELECT public.fn_foreign_schema_clone($1, $2)"

	if _, err := conn.Exec(query, tenantID, templateTenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFForeignSchemaCloneTx func
func ExecSFForeignSchemaCloneTx(conn *sqlx.Tx, tenantID, templateTenantID int) error {
	query := "SELECT public.fn_foreign_schema_clone($1, $2)"

	if _, err := conn.Exec(query, tenantID, templateTenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ExecSFCloneUsersTx func
func ExecSFCloneUsersTx(conn *sqlx.Tx, tenantID, templateTenantID int) error {
	query := "SELECT public.fn_clone_users($1, $2)"

	if _, err := conn.Exec(query, tenantID, templateTenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
