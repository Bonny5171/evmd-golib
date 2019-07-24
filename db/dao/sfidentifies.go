package dao

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/gopkgs/db/model"
)

func SaveSFIdentifyDocOrg(conn *sqlx.DB, iden model.SFIdentity) (id int, err error) {
	t := time.Now()

	row := conn.QueryRow("SELECT id FROM itgr.sf_identity WHERE tenant_id = $1", iden.TenantID)
	err = row.Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, errors.Wrap(err, "row.Scan()")
		}
	}

	if id == 0 {
		query := `INSERT INTO itgr.sf_identity (tenant_id, execution_id, "name", doc_org, doc_objects, doc_meta_data, is_active, created_at, updated_at, is_deleted, deleted_at)
			      VALUES($1, $2, $3, $4, $5, $6, true, $7, $8, false, null) 
			      RETURNING id;`

		err = conn.QueryRowx(query, iden.TenantID, iden.ExecutionID, iden.Name, iden.DocOrg, iden.DocObjects, iden.DocMetaData, t, t).Scan(&id)
		if err != nil {
			return 0, errors.Wrap(err, "conn.QueryRowx()")
		}

		if id <= 0 {
			err = errors.New("An error has occurred while inserting on 'itgr.sf_identity'")
			return id, err
		}
	} else {
		query := `UPDATE itgr.sf_identity 
		            SET tenant_id = $1, execution_id = $2, "name" = $3, doc_org = $4, doc_objects = $5, doc_meta_data = $6, is_active = true, updated_at = $7, is_deleted = false, deleted_at = null
			      WHERE id = $8;`

		if _, err := conn.Exec(query, iden.TenantID, iden.ExecutionID, iden.Name, iden.DocOrg, iden.DocObjects, iden.DocMetaData, t, id); err != nil {
			return 0, errors.Wrap(err, "conn.Exec()")
		}
	}

	return id, nil
}

func UpdateSFIdentifyDocObjects(conn *sqlx.DB, iden model.SFIdentity) error {
	t := time.Now()

	query := `UPDATE itgr.sf_identity
                 SET doc_objects = $1, updated_at = $2
			   WHERE id = $3;`

	if _, err := conn.Exec(query, iden.DocObjects, t, iden.ID); err != nil {
		return errors.Wrap(err, "conn.Exec()")
	}

	return nil
}
