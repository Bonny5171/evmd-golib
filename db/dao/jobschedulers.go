package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/evmd-golib/db/model"
)

// GetSchedules retorna todos os 'jobs' agendados que deverão ser executadas
func GetSchedules(conn *sqlx.DB, tenantID int) (s []model.JobScheduler, err error) {
	const query = `
	  SELECT j.id, j."name", j.tenant_id, t."name" AS tenant_name, t.org_id, j.queue, j.job_name, j.description, j.parameters, j.cron, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.is_active, j.is_deleted 
	    FROM itgr.job_scheduler j
	   INNER JOIN public.tenant t ON j.tenant_id = t.id
	   WHERE j.tenant_id = $1
	   ORDER BY j.id;`

	err = conn.Select(&s, query, tenantID)
	if err != nil {
		return nil, errors.Wrap(err, "db.Conn.Select()")
	}

	return s, nil
}

// GetSchedulesByOrg retorna todos os 'jobs' agendados que deverão ser executadas
func GetSchedulesByOrg(conn *sqlx.DB, orgID string) (s []model.JobScheduler, err error) {
	const query = `
	  SELECT j.id, j."name", j.tenant_id, t."name" AS tenant_name, t.org_id, j.queue, j.job_name, j.description, j.parameters, j.cron, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.is_active, j.is_deleted 
	    FROM itgr.job_scheduler j
	   INNER JOIN public.tenant t ON j.tenant_id = t.id
	   WHERE t.org_id = $1
	   ORDER BY j.id;`

	err = conn.Select(&s, query, orgID)
	if err != nil {
		return nil, errors.Wrap(err, "db.Conn.Select()")
	}

	return s, nil
}

// GetJob retorna os dados de um 'job'
func GetJob(conn *sqlx.DB, tenantID int, name string) (s model.JobScheduler, err error) {
	const query = `
	  SELECT j.id, j."name", j.tenant_id, t."name" AS tenant_name, t.org_id, j.queue, j.job_name, j.description, j.parameters, j.cron, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.doc_meta_data, j.is_active, j.is_deleted 
	    FROM itgr.job_scheduler j
	   INNER JOIN public.tenant t ON j.tenant_id = t.id
	   WHERE j.tenant_id = $1
	     AND j.name = $2
	   LIMIT 1;`

	err = conn.Get(&s, query, tenantID, name)
	if err != nil {
		return s, errors.Wrap(err, "conn.Get()")
	}

	return s, nil
}

// GetJobByID retorna os dados de um 'job'
func GetJobByID(conn *sqlx.DB, jobID int64) (s model.JobScheduler, err error) {
	const query = `
	  SELECT j.id, j."name", j.tenant_id, t."name" AS tenant_name, t.org_id, j.queue, j.job_name, j.description, j.parameters, j.cron, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.doc_meta_data, j.is_active, j.is_deleted 
	    FROM itgr.job_scheduler j
	   INNER JOIN public.tenant t ON j.tenant_id = t.id
	   WHERE j.id = $1
	   LIMIT 1;`

	err = conn.Get(&s, query, jobID)
	if err != nil {
		return s, errors.Wrap(err, "conn.Get()")
	}

	return s, nil
}

func SetCronJobSchedule(conn *sqlx.DB, jobID int64, cronexpr string) error {
	query := `UPDATE itgr.job_scheduler SET cron = $1 WHERE id = $2;`

	if _, err := conn.Exec(query, cronexpr, jobID); err != nil {
		return errors.Wrap(err, "conn.Exec()")
	}

	return nil
}

func ActiveJobSchedule(conn *sqlx.DB, jobID int64, active bool) error {
	query := `UPDATE itgr.job_scheduler SET is_active = $1 WHERE id = $2;`

	if _, err := conn.Exec(query, active, jobID); err != nil {
		return errors.Wrap(err, "conn.Exec()")
	}

	return nil
}
