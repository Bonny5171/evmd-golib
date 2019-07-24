package dao

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/gopkgs/db/model"
)

func SaveSFObject(conn *sqlx.DB, obj model.SFObject) (id int, err error) {
	t := time.Now()

	row := conn.QueryRow("SELECT id FROM itgr.sf_object WHERE tenant_id = $1 AND sf_object_name = $2", obj.TenantID, obj.Name)
	err = row.Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, errors.Wrap(err, "row.Scan()")
		}
	}

	if obj.DocDescribe.IsNull() {
		obj.DocDescribe = []byte("{}")
	}

	if id == 0 {
		query := `INSERT INTO itgr.sf_object (tenant_id, execution_id, sf_object_name, doc_describe, doc_meta_data, is_active, created_at, updated_at, is_deleted, deleted_at)
			  VALUES($1, $2, $3, $4, $5, true, $6, $7, false, null) 
			  RETURNING id;`

		err = conn.QueryRowx(query, obj.TenantID, obj.ExecutionID, obj.Name, obj.DocDescribe, obj.DocMetaData, t, t).Scan(&id)
		if err != nil {
			return 0, errors.Wrap(err, "conn.QueryRowx()")
		}

		if id <= 0 {
			err = errors.New("An error has occurred while inserting on 'itgr.sf_object'")
			return 0, err
		}
	} else {
		query := `UPDATE itgr.sf_object 
		            SET tenant_id = $1, execution_id = $2, sf_object_name = $3, doc_describe = $4, doc_meta_data = $5, is_active = true, created_at = $6, updated_at = $7, is_deleted = false, deleted_at = null
			      WHERE id = $8;`

		if _, err := conn.Exec(query, obj.TenantID, obj.ExecutionID, obj.Name, obj.DocDescribe, obj.DocMetaData, t, t, id); err != nil {
			return 0, errors.Wrap(err, "conn.Exec()")
		}
	}

	return id, nil
}
