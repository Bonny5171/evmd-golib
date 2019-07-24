package dao

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/gopkgs/db/model"
)

// InsertExecution
func InsertExecution(conn *sqlx.DB, obj model.Execution) (r int64, err error) {
	t := time.Now()

	query := `INSERT INTO itgr.execution (job_faktory_id, job_scheduler_id, tenant_id, schema_id, status_id, doc_meta_data, is_active, created_at, updated_at, is_deleted)
			  VALUES($1, $2, $3, $4, $5, $6, true, $7, $8, false)
			  RETURNING id;`

	err = conn.QueryRowx(query, obj.JobFaktoryID, obj.JobSchedulerID, obj.TenantID, obj.SchemaID, obj.StatusID, obj.DocMetaData, t, t).Scan(&r)
	if err != nil {
		return 0, errors.Wrap(err, "conn.QueryRowx()")
	}

	if r <= 0 {
		err = errors.New("An error has occurred while inserting on 'itgr.execution'")
		return r, err
	}

	return r, nil
}

func UpdateExecution(conn *sqlx.DB, obj model.Execution) error {
	t := time.Now()

	query := `UPDATE itgr.execution
			  SET status_id = $1, doc_meta_data = $2, updated_at = $3
			  WHERE id = $4;`

	if _, err := conn.Exec(query, obj.StatusID, obj.DocMetaData, t, obj.ID); err != nil {
		return errors.Wrap(err, "conn.Exec()")
	}

	return nil
}
