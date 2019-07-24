package dao

import (
	"bitbucket.org/everymind/evmd-golib/db/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func GetDevices(conn *sqlx.DB, tid int, execID int64) (d []model.Device, err error) {
	query := `SELECT d.device_id, count(*) AS qty
			    FROM public.device_data d
			   WHERE d.tenant_id = $1
			     AND d.is_deleted = false
			     AND d.execution_id = $2
			   GROUP BY d.device_id;`

	err = conn.Select(&d, query, tid, execID)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Select()")
	}

	return d, nil
}

func GetDeviceDataTables(conn *sqlx.DB, tid int, execID int64) (t []*model.DeviceTableField, err error) {
	query := `SELECT o.id AS sf_object_id, 
			         o.sf_object_name, 
			         'sf_'::text || fn_snake_case(o.sf_object_name) AS sfa_table_name, 
			         f.from_to_fields, 
			         pk.sf_field_name AS primary_key, 
			         e.sf_field_name AS external_id
			    FROM itgr.sf_object o
			   INNER JOIN itgr.vw_sf_object_fields_from_to f ON o.tenant_id = f.tenant_id AND o.id = f.sf_object_id
			   INNER JOIN itgr.sf_object_field pk ON o.tenant_id = pk.tenant_id AND o.id = pk.sf_object_id AND  sf_type = 'id'
			    LEFT JOIN itgr.sf_object_field e ON o.tenant_id = e.tenant_id AND o.id = e.sf_object_id AND e.sf_external_id = TRUE AND e.sfa_external_id = TRUE
			   WHERE 'sf_'||fn_snake_case(o.sf_object_name) 
			      IN (SELECT DISTINCT table_name FROM public.device_data WHERE tenant_id = $1 AND execution_id = $2);`

	err = conn.Select(&t, query, tid, execID)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Select()")
	}

	return t, nil
}

func GetDeviceDataIDs(conn *sqlx.DB, tid int, device string, execID int64) (d []string, err error) {
	query := `SELECT d.id
			    FROM public.device_data d
			   WHERE d.tenant_id = $1
			     AND d.device_id = $2
			     AND d.is_deleted = FALSE
			     AND d.execution_id = $3
			   ORDER BY d.device_created_at, d.created_at, d.updated_at;`

	err = conn.Select(&d, query, tid, device, execID)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Select()")
	}

	return d, nil
}

func GetDeviceData(conn *sqlx.DB, id string) (d model.DeviceData, err error) {
	query := `SELECT d.id, d.tenant_id, d.schema_name, d.table_name, o.id AS sf_object_id, o.sf_object_name, d.user_id, d.pk, d.sf_id, to_jsonb(d.json_data::jsonb) AS json_data, d.app_id, d.device_id, d.device_created_at
			    FROM public.device_data d
			   INNER JOIN itgr.sf_object o ON d.tenant_id = o.tenant_id AND d.table_name = 'sf_'::text || fn_snake_case(o.sf_object_name)
			   WHERE d.id = $1
			   LIMIT 1;`

	if err = conn.Get(&d, query, id); err != nil {
		err = errors.Wrap(err, "conn.Get()")
	}

	return
}

func SetDeviceDatasToExecution(conn *sqlx.DB, tid int, execID int64) error {
	query := `UPDATE public.device_data
			     SET execution_id = $1
			   WHERE tenant_id = $2 
			     AND is_deleted = FALSE;`

	if _, err := conn.Exec(query, execID, tid); err != nil {
		return errors.Wrap(err, "conn.Exec()")
	}

	return nil
}

func SetDeviceDataToDelete(conn *sqlx.DB, id string) error {
	query := `UPDATE public.device_data
			  SET is_deleted = TRUE, deleted_at = NOW()
			  WHERE id = $1;`

	if _, err := conn.Exec(query, id); err != nil {
		return errors.Wrap(err, "conn.Exec()")
	}

	return nil
}

func InsertDeviceDataLogError(conn *sqlx.DB, obj model.DeviceData, execID int64, err error) error {
	query := `INSERT INTO itgr.device_data_log_error (original_id, tenant_id, device_created_at, schema_name, table_name, pk, device_id, user_id, sf_id, original_json_data, brewed_json_data, app_id, execution_id, error_description, created_at, updated_at) 
	          VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW());`

	_, e := conn.Exec(query, obj.ID, obj.TenantID, obj.DeviceCreatedAt, obj.SchemaName, obj.TableName, obj.PK, obj.DeviceID, obj.UserID, obj.SfID, obj.JSONData, obj.BrewedJSONData, obj.AppID, execID, err.Error())
	if e != nil {
		return errors.Wrap(e, "conn.Exec()")
	}

	return nil
}

func PurgeAllDeviceDataToDelete(conn *sqlx.DB, tid int) (err error) {
	query := `DELETE FROM public.device_data
			   WHERE tenant_id = $1
			     AND is_deleted = TRUE;`

	_, err = conn.Exec(query, tid)
	if err != nil {
		return errors.Wrap(err, "conn.Exec()")
	}

	return nil
}
