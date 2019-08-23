package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

func GetFieldsBase64(conn *sqlx.DB, tenantId int, objID int) (f []model.SFObjectField, err error) {
	query := `
		SELECT id, tenant_id, sf_object_id, sf_field_name
		  FROM itgr.sf_object_field
		 WHERE tenant_id = $1
		   AND sf_object_id = $2
		   AND sf_type = 'base64';`

	err = conn.Select(&f, query, tenantId, objID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return
}
