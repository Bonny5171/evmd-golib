package dao

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/CognyHub/evmd-golib/db"
	"github.com/CognyHub/evmd-golib/db/model"
)

// InsertExecution func
func InsertExecution(conn *sqlx.DB, obj model.Execution) (r int64, err error) {
	t := time.Now()

	query := `INSERT INTO itgr.execution (job_faktory_id, job_scheduler_id, job_scheduler_name, tenant_id, schema_id, status_id, doc_meta_data, is_active, created_at, updated_at, is_deleted)
			  VALUES($1, $2, $3, $4, $5, $6, $7, true, $8, $9, false)
			  RETURNING id;`

	err = conn.QueryRowx(query, obj.JobFaktoryID, obj.JobSchedulerID, obj.JobSchedulerName, obj.TenantID, obj.SchemaID, obj.StatusID, obj.DocMetaData, t, t).Scan(&r)
	if err != nil {
		return 0, db.WrapError(err, "conn.QueryRowx()")
	}

	if r <= 0 {
		err = errors.New("An error has occurred while inserting on 'itgr.execution'")
		return r, err
	}

	return r, nil
}

// UpdateExecution func
func UpdateExecution(conn *sqlx.DB, obj model.Execution) error {
	t := time.Now()

	query := `UPDATE itgr.execution
			  SET status_id = $1, doc_meta_data = $2, updated_at = $3
			  WHERE id = $4 AND tenant_id = $5;`

	if _, err := conn.Exec(query, obj.StatusID, obj.DocMetaData, t, obj.ID, obj.TenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// GetLastExecutionByJobID func
func GetLastExecutionByJobID(conn *sqlx.DB, tenantID int, jobID string) (e model.Execution, err error) {
	query := `SELECT e.id, e.tenant_id, e.job_faktory_id, e.job_scheduler_id, e.job_scheduler_name, e.schema_id, e.status_id, t.name AS status_name, e.doc_meta_data, e.is_active, e.created_at, e.updated_at, e.is_deleted, e.deleted_at 
				FROM itgr.execution e
			   INNER JOIN itgr.status t ON e.tenant_id = t.tenant_id AND e.status_id = t.id
			   WHERE e.tenant_id = $1 AND e.job_faktory_id = $2
			   ORDER BY e.id DESC
			   LIMIT 1;`

	err = conn.Get(&e, query, tenantID, jobID)
	if err != nil {
		return e, db.WrapError(err, "conn.Get()")
	}

	return e, nil
}
