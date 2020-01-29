package dao

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

// GetSFObject
func GetSFObject(conn *sqlx.DB, tid int, objectName string) (o model.SFObject, err error) {
	const query = `
		SELECT id, tenant_id, sf_object_name, doc_describe, doc_meta_data
		  FROM itgr.sf_object 
		 WHERE tenant_id = $1 AND is_active = TRUE AND is_deleted = FALSE AND sf_object_name = $2 
		 LIMIT 1;`

	if err = conn.Get(&o, query, tid, objectName); err != nil {
		err = db.WrapError(err, "conn.Get()")
		return
	}

	return
}

func SaveSFObject(conn *sqlx.DB, obj model.SFObject) (id int, err error) {
	t := time.Now()

	row := conn.QueryRow("SELECT id FROM itgr.sf_object WHERE tenant_id = $1 AND sf_object_name = $2", obj.TenantID, obj.Name)
	err = row.Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, db.WrapError(err, "row.Scan()")
		}
	}

	if obj.DocDescribe.IsNull() {
		obj.DocDescribe = []byte("{}")
	}

	if id == 0 {
		query := `
		INSERT INTO itgr.sf_object (tenant_id, execution_id, sf_object_name, sfa_name, doc_describe, doc_meta_data, is_package, is_active, created_at, updated_at, is_deleted, deleted_at)
		VALUES($1, $2, $3, 'sf_'||public.fn_snake_case(public.fn_remove_namespace($1, $3)), $4, $5, public.fn_is_package($1, $3), true, $6, $6, false, null) 
		RETURNING id;`

		err = conn.QueryRowx(query, obj.TenantID, obj.ExecutionID, obj.Name, obj.DocDescribe, obj.DocMetaData, t).Scan(&id)
		if err != nil {
			return 0, db.WrapError(err, "conn.QueryRowx()")
		}

		if id <= 0 {
			err = errors.New("An error has occurred while inserting on 'itgr.sf_object'")
			return 0, err
		}
	} else {
		query := `
		UPDATE itgr.sf_object 
		SET tenant_id = $1, 
			execution_id = $2, 
			sf_object_name = $3, 
			sfa_name = 'sf_'||public.fn_snake_case(public.fn_remove_namespace($1, $3)), 
			doc_describe = $4, 
			doc_meta_data = $5, 
			is_package = public.fn_is_package($1, $3), 
			is_active = true, 
			created_at = $6, 
			updated_at = $6, 
			is_deleted = false, 
			deleted_at = null
		WHERE id = $7;`

		if _, err := conn.Exec(query, obj.TenantID, obj.ExecutionID, obj.Name, obj.DocDescribe, obj.DocMetaData, t, id); err != nil {
			return 0, db.WrapError(err, "conn.Exec()")
		}
	}

	return id, nil
}
