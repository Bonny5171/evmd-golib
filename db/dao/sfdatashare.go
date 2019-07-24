package dao

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/gopkgs/db/model"
)

func SaveSFDataShare(conn *sqlx.DB, data model.SFDataShare) (id int, err error) {
	t := time.Now()

	query := `INSERT INTO itgr.sf_data_share (tenant_id, execution_id, status_id, sf_object_id, sf_object_name, user_id, doc_original_id, is_active, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, true, $8, $9)
			  RETURNING id;`

	err = conn.QueryRowx(query, data.TenantID, data.ExecutionID, data.StatusID, data.ObjectID, data.ObjectName, data.UserID, data.DocOriginalID, t, t).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "conn.QueryRowx()")
	}

	if id <= 0 {
		err = errors.New("An error has occurred while inserting on 'itgr.sf_data_share'")
		return id, err
	}

	return id, nil
}

func PurgeAllDataShareETLSuccess(conn *sqlx.DB, tid int) (err error) {
	statuses, err := GetStatuses(conn, tid, EnumTypeStatusETL)
	if err != nil {
		return errors.Wrap(err, "dao.GetStatuses()")
	}

	statusEtlSuccess := statuses.GetId(EnumStatusEtlSuccess.String())

	query := `DELETE FROM itgr.sf_data_share
			   WHERE tenant_id = $1
			     AND status_id = $2;`

	_, err = conn.Exec(query, tid, statusEtlSuccess)
	if err != nil {
		return errors.Wrap(err, "conn.Exec()")
	}

	return nil
}
