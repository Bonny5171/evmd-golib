package dao

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
	m "bitbucket.org/everymind/evmd-golib/modelbase"
)

//GetDevices func
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

//GetDevicesByGroup func
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

//GetDeviceByIdGroupedByGroup func
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

//GetDeviceDataTables func
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
				LEFT JOIN itgr.vw_schemas_objects so ON o.tenant_id = so.tenant_id AND o.id = so.sf_object_id AND so.is_active = TRUE
				WHERE EXISTS (SELECT DISTINCT table_name FROM public.device_data WHERE o.sfa_name = table_name AND so.tenant_id = $1 AND is_active = TRUE AND execution_id = $2);`

	err = conn.Select(&t, query, tid, execID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return t, nil
}

//GetDeviceDataIDs func
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

//GetDeviceDataIDsByGroupID func
func GetDeviceDataIDsByGroupID(conn *sqlx.DB, tid int, deviceID string, groupID m.NullString, execID int64) (d []string, err error) {
	var (
		query  = strings.Builder{}
		params = []interface{}{tid, deviceID, execID}
	)

	query.WriteString("SELECT d.id FROM public.device_data d WHERE d.tenant_id = $1 AND d.is_active = TRUE AND d.is_deleted = FALSE AND d.device_id = $2 AND d.execution_id = $3 ")
	if groupID.Valid {
		query.WriteString("AND d.group_id = $4 ")
		params = append(params, groupID.String)
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

//GetDeviceData func
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

//GetDeviceDataUsersToProcess func
func GetDeviceDataUsersToProcess(conn *sqlx.DB, tid int, execID int64) (d []string, err error) {
	query := `SELECT DISTINCT user_id FROM public.device_data
			  WHERE tenant_id = $1 AND execution_id = $2 AND is_active = TRUE AND is_deleted = FALSE;`

	err = conn.Select(&d, query, tid, execID)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return d, nil
}

//SetDeviceDataToExecution func
func SetDeviceDataToExecution(conn *sqlx.DB, tid int, execID int64, retry int, maxWorkers int64) error {
	query := `
		WITH a AS (
			SELECT * FROM (
			SELECT DISTINCT ON (group_id) group_id, created_at
			FROM public.device_data
			WHERE tenant_id = $3 AND is_active = TRUE AND is_deleted = FALSE
			ORDER BY group_id, created_at ASC NULLS LAST
			) sub
		ORDER BY created_at ASC NULLS LAST LIMIT $4
		)
		UPDATE
			public.device_data d
		SET execution_id = CASE public.fn_check_retry (try,$1) WHEN TRUE THEN $2 ELSE execution_id END, updated_at = now()
		FROM a
		WHERE
		a.group_id = d.group_id;
		`

	// query := `UPDATE public.device_data
	//           SET execution_id = CASE public.fn_check_retry(try,$1) WHEN TRUE THEN $2 ELSE execution_id END, updated_at = now()
	//           WHERE tenant_id = $3 AND is_active = TRUE AND is_deleted = FALSE;`

	if _, err := conn.Exec(query, retry, execID, tid, maxWorkers); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

//DeactivateDeviceDataRows func
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

//SetTryDeviceDataRows func
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

//SetDeviceDataToDelete func
func SetDeviceDataToDelete(conn *sqlx.DB, id string) error {
	query := `UPDATE public.device_data
			  SET is_deleted = TRUE, deleted_at = NOW()
			  WHERE id = $1;`

	if _, err := conn.Exec(query, id); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

//PurgeAllDeviceDataToDelete func
func PurgeAllDeviceDataToDelete(conn *sqlx.DB, tid int) (err error) {
	query := `DELETE FROM public.device_data
			  WHERE tenant_id = $1 AND is_deleted = TRUE;`

	_, err = conn.Exec(query, tid)
	if err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

//InsertDeviceDataLog func
func InsertDeviceDataLog(conn *sqlx.DB, obj model.DeviceData, execID int64, statusID int16, statusName string) (id int64, err error) {
	query := `INSERT INTO itgr.device_data_log (
				original_id,tenant_id,device_created_at,schema_name,table_name,pk,device_id,user_id,action_type,sf_id,original_json_data,
				app_id,execution_id,status_id,status_name,external_id,group_id,sequential,try,created_at,updated_at) 
			  VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, NOW(), NOW())
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
	params = append(params, statusName)          // 15
	params = append(params, obj.ExternalID)      // 16
	params = append(params, obj.GroupID)         // 17
	params = append(params, obj.Sequential)      // 18
	params = append(params, obj.Try)             // 19

	if e := conn.QueryRowx(query, params...).Scan(&id); e != nil {
		err = db.WrapError(e, "conn.QueryRowx(query, params...).Scan(&id)")
		return
	}

	return id, nil
}

//UpdateDeviceDataLog func
func UpdateDeviceDataLog(conn *sqlx.DB, brewedJSON m.JSONB, logID int64, statusID int16, statusName string, try int, err error) error {
	params := make([]interface{}, 0)
	params = append(params, logID)
	params = append(params, statusID)
	params = append(params, statusName)
	params = append(params, brewedJSON)
	params = append(params, try)

	query := strings.Builder{}
	query.WriteString("UPDATE itgr.device_data_log SET status_id = $2, status_name = $3, brewed_json_data = $4, try = $5, ")
	if err != nil {
		query.WriteString("error = $6, ")
		params = append(params, err.Error())
	}
	query.WriteString("updated_at = NOW() WHERE id = $1;")

	if _, err := conn.Exec(query.String(), params...); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

//GetDeviceDataLogReport func
func GetDeviceDataLogReport(conn *sqlx.DB, tid int) (logs []model.DeviceDataLogReport, err error) {
	query := `
		SELECT id, tenant_id, created_at, updated_at, is_active, is_deleted, log_status_name, log_error, device_data_id, device_created_at, table_name, pk, sf_id, action_type, external_id, device_id, user_id, group_id, original_json_data, brewed_json_data, execution_id, execution_status_name, execution_job_faktory_id, reported
		FROM itgr.vw_report_device_data_log_errors WHERE tenant_id = $1`

	err = conn.Select(&logs, query, tid)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return logs, nil
}
