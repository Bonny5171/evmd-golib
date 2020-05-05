package dao

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
	m "bitbucket.org/everymind/evmd-golib/modelbase"
)

func GetDevices(conn *sqlx.DB, tid int, execID int64) (d []model.Device, err error) {
	query := `SELECT d.device_id, count(*) AS qty
			    FROM public.device_data d
			   WHERE d.tenant_id = $1
			     AND d.is_active = TRUE
			     AND d.is_deleted = FALSE
			     AND d.execution_id = $2
			   GROUP BY d.device_id;`

	err = conn.Select(&d, query, tid, execID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return d, nil
}

func GetDevicesByGroup(conn *sqlx.DB, tid int, execID int64) (d []model.Device, err error) {
	query := `SELECT d.device_id, d.group_id, count(*) AS qty
			    FROM public.device_data d
			   WHERE d.tenant_id = $1
			     AND d.is_active = TRUE
			     AND d.is_deleted = FALSE
			     AND d.execution_id = $2
			   GROUP BY d.device_id, d.group_id;`

	err = conn.Select(&d, query, tid, execID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return d, nil
}

func GetDeviceByIdGroupedByGroup(conn *sqlx.DB, tid int, execID int64, deviceID string) (d []model.Device, err error) {
	query := `SELECT d.device_id, d.group_id, count(*) AS qty
			    FROM public.device_data d
			   WHERE d.tenant_id = $1
			     AND d.is_active = TRUE
			     AND d.is_deleted = FALSE
			     AND d.execution_id = $2
			     AND d.device_id = $3
			   GROUP BY d.device_id, d.group_id;`

	err = conn.Select(&d, query, tid, execID, deviceID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return d, nil
}

func GetDeviceDataTables(conn *sqlx.DB, tid int, execID int64) (t []*model.DeviceTableField, err error) {
	query := `SELECT o.id AS sf_object_id, 
			         o.sf_object_name, 
			         o.sfa_name AS sfa_table_name, 
			         f.from_to_fields, 
			         pk.sf_field_name AS primary_key, 
					 e.sf_field_name AS external_id,
					 so.sfa_pks
			    FROM itgr.sf_object o
			   INNER JOIN itgr.vw_sf_object_fields_from_to f ON o.tenant_id = f.tenant_id AND o.id = f.sf_object_id
			   INNER JOIN itgr.sf_object_field pk ON o.tenant_id = pk.tenant_id AND o.id = pk.sf_object_id AND  sf_type = 'id'
				LEFT JOIN itgr.sf_object_field e ON o.tenant_id = e.tenant_id AND o.id = e.sf_object_id AND e.sf_external_id = TRUE AND e.sfa_external_id = TRUE
<<<<<<< HEAD
				LEFT JOIN itgr.vw_schemas_object so ON o.tenant_id = so.tenant_id AND o.sf_object_id = so.object_id
=======
				LEFT JOIN itgr.vw_schemas_objects so ON o.tenant_id = so.tenant_id AND o.id = so.sf_object_id
>>>>>>> develop
			   WHERE o.sfa_name IN (SELECT DISTINCT table_name FROM public.device_data WHERE tenant_id = $1 AND is_active = TRUE AND execution_id = $2);`

	err = conn.Select(&t, query, tid, execID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return t, nil
}

func GetDeviceDataIDs(conn *sqlx.DB, tid int, device string, execID int64) (d []string, err error) {
	query := `SELECT d.id
			    FROM public.device_data d
			   WHERE d.tenant_id = $1
				 AND d.device_id = $2
				 AND d.is_active = TRUE
			     AND d.is_deleted = FALSE
			     AND d.execution_id = $3
			   ORDER BY d.sequential ASC;`

	err = conn.Select(&d, query, tid, device, execID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return d, nil
}

func GetDeviceDataIDsByGroupID(conn *sqlx.DB, tid int, device_id string, group_id m.NullString, execID int64) (d []string, err error) {
	var (
		query  = strings.Builder{}
		params = []interface{}{tid, device_id, execID}
	)

	query.WriteString("SELECT d.id FROM public.device_data d WHERE d.tenant_id = $1 AND d.is_active = TRUE AND d.is_deleted = FALSE AND d.device_id = $2 AND d.execution_id = $3 ")
	if group_id.Valid {
		query.WriteString("AND d.group_id = $4 ")
		params = append(params, group_id.String)
	} else {
		query.WriteString("AND d.group_id = NULL ")
	}
	query.WriteString("ORDER BY d.sequential ASC;")

	err = conn.Select(&d, query.String(), params...)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return d, nil
}

func GetDeviceData(conn *sqlx.DB, id string) (d model.DeviceData, err error) {
	query := `SELECT d.id, d.tenant_id, d.schema_name, d.table_name, o.id AS sf_object_id, o.sf_object_name, d.user_id, d.pk, d.external_id, d.sf_id, d.action_type,
					 to_jsonb(regexp_replace(d.json_data, E'[\\n\\r\\f\\u000B\\u0085\\u2028\\u2029]+', ' ', 'g')::jsonb) AS json_data, 
					 d.app_id, d.device_id, d.device_created_at, d.group_id, d.sequential, d.try, d.is_active, d.is_deleted
			  FROM public.device_data d
			  INNER JOIN itgr.sf_object o ON d.tenant_id = o.tenant_id AND d.table_name = o.sfa_name
			  WHERE d.id = $1
			  LIMIT 1;`

	if err = conn.Get(&d, query, id); err != nil {
		err = db.WrapError(err, "conn.Get()")
	}

	return
}

func GetDeviceDataUsersToProcess(conn *sqlx.DB, tid int, execID int64) (d []string, err error) {
	query := `SELECT DISTINCT user_id FROM public.device_data
			  WHERE tenant_id = $1 AND execution_id = $2 AND is_active = TRUE AND is_deleted = FALSE;`

	err = conn.Select(&d, query, tid, execID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return d, nil
}

func SetDeviceDataToExecution(conn *sqlx.DB, tid int, execID int64, retry int) error {
	query := `UPDATE public.device_data
	          SET execution_id = CASE public.fn_check_retry(try,$1) WHEN TRUE THEN $2 ELSE execution_id END, updated_at = now()
  	          WHERE tenant_id = $3 AND is_active = TRUE AND is_deleted = FALSE;`

	if _, err := conn.Exec(query, retry, execID, tid); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func DeactivateDeviceDataRows(conn *sqlx.DB, tid int, retry int) error {
	query := `
		WITH b AS (
			WITH a AS (
				SELECT id,group_id,public.fn_check_retry(try,$1) AS retry
				FROM public.device_data
				WHERE tenant_id = $2
				AND is_active = TRUE
				AND is_deleted = FALSE
			)
			SELECT DISTINCT a.group_id
			FROM a
			WHERE a.retry = FALSE
		)
		UPDATE public.device_data d
		SET is_active = FALSE 
		FROM b
		WHERE b.group_id = d.group_id;`

	if _, err := conn.Exec(query, retry, tid); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func SetTryDeviceDataRows(conn *sqlx.DB, id string, retry int) (try int, err error) {
	query := `UPDATE public.device_data 
			  SET try = CASE public.fn_check_retry(try,$1) WHEN TRUE THEN try + 1 ELSE try END, updated_at = now()
			  WHERE id = $2 
			  RETURNING try;`

	if e := conn.QueryRowx(query, retry, id).Scan(&try); e != nil {
		err = db.WrapError(e, "conn.QueryRowx(query, retry, id).Scan(&try)")
		return
	}

	return try, nil
}

func SetDeviceDataToDelete(conn *sqlx.DB, id string) error {
	query := `UPDATE public.device_data
			  SET is_deleted = TRUE, deleted_at = NOW()
			  WHERE id = $1;`

	if _, err := conn.Exec(query, id); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func PurgeAllDeviceDataToDelete(conn *sqlx.DB, tid int) (err error) {
	query := `DELETE FROM public.device_data
			  WHERE tenant_id = $1 AND is_deleted = TRUE;`

	_, err = conn.Exec(query, tid)
	if err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func InsertDeviceDataLog(conn *sqlx.DB, obj model.DeviceData, execID int64, statusID int16) (id int64, err error) {
	query := `INSERT INTO itgr.device_data_log (
				original_id,tenant_id,device_created_at,schema_name,table_name,pk,device_id,user_id,action_type,sf_id,original_json_data,
				app_id,execution_id,status_id,external_id,group_id,sequential,try,created_at,updated_at) 
			  VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, NOW(), NOW())
			  RETURNING id;`

	params := make([]interface{}, 0)
	params = append(params, obj.ID)              // 1
	params = append(params, obj.TenantID)        // 2
	params = append(params, obj.DeviceCreatedAt) // 3
	params = append(params, obj.SchemaName)      // 4
	params = append(params, obj.TableName)       // 5
	params = append(params, obj.PK)              // 6
	params = append(params, obj.DeviceID)        // 7
	params = append(params, obj.UserID)          // 8
	params = append(params, obj.ActionType)      // 9
	params = append(params, obj.SfID)            // 10
	params = append(params, obj.JSONData)        // 11
	params = append(params, obj.AppID)           // 12
	params = append(params, execID)              // 13
	params = append(params, statusID)            // 14
	params = append(params, obj.ExternalID)      // 15
	params = append(params, obj.GroupID)         // 16
	params = append(params, obj.Sequential)      // 17
	params = append(params, obj.Try)             // 18

	if e := conn.QueryRowx(query, params...).Scan(&id); e != nil {
		err = db.WrapError(e, "conn.QueryRowx(query, params...).Scan(&id)")
		return
	}

	return id, nil
}

func UpdateDeviceDataLog(conn *sqlx.DB, brewedJSON m.JSONB, logID int64, statusID int16, try int, err error) error {
	params := make([]interface{}, 0)
	params = append(params, logID)
	params = append(params, statusID)
	params = append(params, brewedJSON)
	params = append(params, try)

	query := strings.Builder{}
	query.WriteString("UPDATE itgr.device_data_log SET status_id = $2, brewed_json_data = $3, try = $4, ")
	if err != nil {
		query.WriteString("error = $5, ")
		params = append(params, err.Error())
	}
	query.WriteString("updated_at = NOW() WHERE id = $1;")

	if _, err := conn.Exec(query.String(), params...); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
