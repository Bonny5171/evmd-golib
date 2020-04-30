package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

type (
	SchemaType int
)

const (
	EnumSchemaTypeInboud SchemaType = iota
	EnumSchemaTypeOutbound
)

func (t SchemaType) String() string {
	n := [...]string{"inbound", "outbound"}
	if t < EnumSchemaTypeInboud || t > EnumSchemaTypeOutbound {
		return ""
	}
	return n[t]
}

// GetSchemaObjects
func GetSchemaObjects(conn *sqlx.DB, tenantID, schemaID int) (s model.SchemaObjects, err error) {
	const query = `SELECT id, schema_id, sf_object_id, sf_object_name, sequence, raw_command 
				     FROM itgr.schema_object 
				    WHERE tenant_id = $1 
					  AND schema_id = $2
					  AND is_active = TRUE
					  AND is_deleted = FALSE
				    ORDER BY "sequence", id;`

	err = conn.Select(&s, query, tenantID, schemaID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return s, nil
}

// GetAllSchemaObjectsToProcess
func GetAllSchemaObjectsToProcess(conn *sqlx.DB, tenantID int, schemaObjectName string, schemaType SchemaType) (s model.SchemaObjectToProcesses, err error) {
	const query = `
		SELECT v.id, v.schema_id, v.schema_name, v.tenant_id, t."name" AS tenant_name, v."type", v.sf_object_id, v.sf_object_name, v.doc_fields, 
		       v."sequence", v.filter, v.raw_command, v.sf_last_modified_date, v.layoutable, v.compactlayoutable, v.listviewable, v.sfa_pks
		  FROM itgr.vw_schemas_objects v
		 INNER JOIN public.tenant t ON v.tenant_id = t.id
		 WHERE v.tenant_id = $1 
		   AND v.schema_name = $2 
		   AND v."type" = $3 
		   AND v.is_active = TRUE
		   AND v.is_deleted = FALSE
		   AND v.doc_fields IS NOT NULL
		   AND (v.sf_object_id IS NOT NULL AND v.raw_command IS NULL) OR (v.sf_object_id IS NULL AND v.raw_command IS NOT NULL)
		 ORDER BY v."sequence", v.sf_object_id;`

	err = conn.Select(&s, query, tenantID, schemaObjectName, schemaType.String())
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return s, nil
}

// GetSchemaObjectsToProcess
func GetSchemaObjectsToProcess(conn *sqlx.DB, tenantID int, objectName []string) (s model.SchemaObjectToProcesses, err error) {
	query, args, err := sqlx.In(`
		SELECT v.id, v.schema_id, v.schema_name, v.tenant_id, t."name" AS tenant_name, v."type", v.sf_object_id, v.sf_object_name, v.doc_fields, 
		       v."sequence", v.filter, v.raw_command, v.sf_last_modified_date, v.layoutable, v.compactlayoutable, v.listviewable, v.sfa_pks
		  FROM itgr.vw_schemas_objects v
		 INNER JOIN public.tenant t ON v.tenant_id = t.id
		 WHERE v.tenant_id = ? 
		   AND v.sf_object_name IN (?) 
		   AND v.is_active = TRUE
		   AND v.is_deleted = FALSE
		   AND v.doc_fields IS NOT NULL
		   AND (v.sf_object_id IS NOT NULL AND v.raw_command IS NULL) OR (v.sf_object_id IS NULL AND v.raw_command IS NOT NULL)
		 ORDER BY v."sequence", v.sf_object_id;`, tenantID, objectName)

	if err != nil {
		return nil, db.WrapError(err, "sqlx.In()")
	}

	query = conn.Rebind(query)

	err = conn.Select(&s, query, args...)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return s, nil
}

// GetSchemaShareObjectsToProcess
func GetSchemaShareObjectsToProcess(conn *sqlx.DB, tenantID int) (o model.SFObjectToProcesses, err error) {
	const query = `
			SELECT DISTINCT o.id, o.sf_object_name, o.tenant_id, t."name" AS tenant_name, s."filter" 
			  FROM itgr.sf_object o
			 INNER JOIN public.tenant  t ON o.tenant_id = t.id
			  LEFT JOIN itgr.schema_object s ON o.tenant_id = s.tenant_id AND o.id = s.sf_object_id
			 WHERE o.tenant_id = $1
			   AND o.is_active = TRUE
			   AND o.is_deleted = FALSE
			   AND o.get_share_data = TRUE
			 ORDER BY o.id;`

	err = conn.Select(&o, query, tenantID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return o, nil
}

// UpdateSfObjectIDs
func UpdateSfObjectIDs(conn *sqlx.DB) error {
	const query = `UPDATE itgr.schema_object AS so 
	                  SET sf_object_id = o.id 
					 FROM itgr.sf_object AS o 
					WHERE so.sf_object_name = o.sf_object_name 
					  AND so.tenant_id = o.tenant_id
					  AND so.sf_object_id IS NULL;`

	if _, err := conn.Exec(query); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// UpdateLastModifiedDate
func UpdateLastModifiedDate(conn *sqlx.DB, schemaObjectID int, lastModifiedDate pq.NullTime) error {
	const query = `UPDATE itgr.schema_object
	                  SET sf_last_modified_date = $1 
					WHERE id = $2;`

	if _, err := conn.Exec(query, lastModifiedDate, schemaObjectID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
