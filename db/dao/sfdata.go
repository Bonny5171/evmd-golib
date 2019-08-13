package dao

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/evmd-golib/db/model"
)

func SaveSFData(conn *sqlx.DB, data model.SFData) (id int, err error) {
	t := time.Now()

	query := `INSERT INTO itgr.sf_data (tenant_id, execution_id, record_type_id, status_id, sf_object_name, sf_object_id, doc_id, doc_record, is_active, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true, $9, $10)
			  RETURNING id;`

	err = conn.QueryRowx(query, data.TenantID, data.ExecutionID, data.RecordTypeID, data.StatusID, data.ObjectName, data.ObjectID, data.DocID, data.DocRecord, t, t).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "conn.QueryRowx()")
	}

	if id <= 0 {
		err = errors.New("An error has occurred while inserting on 'itgr.sf_data'")
		return id, err
	}

	return id, nil
}

func PurgeAllDataETLSuccess(conn *sqlx.DB, tid int) (err error) {
	statuses, err := GetStatuses(conn, tid, EnumTypeStatusETL)
	if err != nil {
		return errors.Wrap(err, "dao.GetStatuses()")
	}

	statusEtlSuccess := statuses.GetId(EnumStatusEtlSuccess.String())

	query := `DELETE FROM itgr.sf_data
			   WHERE tenant_id = $1
			     AND status_id = $2;`

	_, err = conn.Exec(query, tid, statusEtlSuccess)
	if err != nil {
		return errors.Wrap(err, "conn.Exec()")
	}

	return nil
}

func GetSfData(conn *sqlx.DB, tenantID int, execID int64, objID int) (d []model.SFData, err error) {
	query := `
		SELECT DISTINCT d.id, d.doc_id, o.sf_object_name
		  FROM itgr.sf_data d
		 INNER JOIN itgr.sf_object o ON d.sf_object_id = o.id
		 WHERE d.tenant_id = $1
		   AND d.execution_id = $2
		   AND d.sf_object_id = $3
	     ORDER BY d.id LIMIT 1;`

	err = conn.Select(&d, query, tenantID, execID, objID)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Select()")
	}

	return
}

func UpdateStatusSfData(conn *sqlx.DB, tenantID int, execID, objectID int64, statusIDFrom, statusIDTo Status) (err error) {
	statuses, err := GetStatuses(conn, tenantID, EnumTypeStatusETL)
	if err != nil {
		return errors.Wrap(err, "dao.GetStatuses()")
	}

	statusFrom := statuses.GetId(statusIDFrom.String())
	statusTo := statuses.GetId(statusIDTo.String())

	query := `UPDATE itgr.sf_data SET status_id = $5, updated_at = NOW() WHERE tenant_id = $1 AND execution_id = $2 AND sf_object_id = $3 AND status_id = $4;`

	if _, err := conn.Exec(query, tenantID, execID, objectID, statusFrom, statusTo); err != nil {
		return errors.Wrap(err, "conn.Exec()")
	}

	return
}
