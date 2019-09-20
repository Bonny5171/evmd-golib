package dao

import (
	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

// GetSchedules retorna todos os 'jobs' agendados que deverão ser executadas
func GetSchedules(conn *sqlx.DB, tenantID, middlewareID int) (s []model.JobScheduler, err error) {
	const query = `
	  SELECT j.id, t.org_id, j.tenant_id, t."name" AS tenant_name, j.middleware_id, j.job_name, j.function_name, j.queue, j.cron, j.parameters, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.description, j.is_active, j.is_deleted 
	    FROM public.job_scheduler j
	   INNER JOIN public.tenant   t ON j.tenant_id = t.id
	   WHERE j.tenant_id = $1
	     AND middleware_id = $2
	   ORDER BY j.id;`

	err = conn.Select(&s, query, tenantID, middlewareID)
	if err != nil {
		return nil, db.WrapError(err, "db.Conn.Select()")
	}

	return s, nil
}

// GetSchedulesByOrg retorna todos os 'jobs' agendados que deverão ser executadas
func GetSchedulesByOrg(conn *sqlx.DB, orgID string, middlewareID int) (s []model.JobScheduler, err error) {
	const query = `
	  SELECT j.id, t.org_id, j.tenant_id, t."name" AS tenant_name, j.middleware_id, j.job_name, j.function_name, j.queue, j.cron, j.parameters, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.description, j.is_active, j.is_deleted 
	    FROM public.job_scheduler j
	   INNER JOIN public.tenant   t ON j.tenant_id = t.id
	   WHERE t.org_id = $1
	     AND middleware_id = $2
	   ORDER BY j.id;`

	err = conn.Select(&s, query, orgID, middlewareID)
	if err != nil {
		return nil, db.WrapError(err, "db.Conn.Select()")
	}

	return s, nil
}

// GetJob retorna os dados de um 'job'
func GetJob(conn *sqlx.DB, tenantID, middlewareID int, name string) (s model.JobScheduler, err error) {
	const query = `
	  SELECT j.id, t.org_id, j.tenant_id, t."name" AS tenant_name, j.middleware_id, j.job_name, j.function_name, j.queue, j.cron, j.parameters, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.description, j.is_active, j.is_deleted 
	    FROM public.job_scheduler j
	   INNER JOIN public.tenant   t ON j.tenant_id = t.id
	   WHERE j.tenant_id = $1
	     AND middleware_id = $2
	     AND j.job_name = $3
	   LIMIT 1;`

	err = conn.Get(&s, query, tenantID, middlewareID, name)
	if err != nil {
		return s, db.WrapError(err, "conn.Get()")
	}

	return s, nil
}

// GetJobByID retorna os dados de um 'job'
func GetJobByID(conn *sqlx.DB, jobID int64) (s model.JobScheduler, err error) {
	const query = `
	  SELECT j.id, t.org_id, j.tenant_id, t."name" AS tenant_name, j.middleware_id, j.job_name, j.function_name, j.queue, j.cron, j.parameters, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.description, j.is_active, j.is_deleted 
	    FROM public.job_scheduler j
	   INNER JOIN public.tenant   t ON j.tenant_id = t.id
	   WHERE j.id = $1
	   LIMIT 1;`

	err = conn.Get(&s, query, jobID)
	if err != nil {
		return s, db.WrapError(err, "conn.Get()")
	}

	return s, nil
}

func SetCronJobSchedule(conn *sqlx.DB, jobID int64, cronexpr string) error {
	query := `UPDATE public.job_scheduler SET cron = $1 WHERE id = $2;`

	if _, err := conn.Exec(query, cronexpr, jobID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

func ActiveJobSchedule(conn *sqlx.DB, jobID int64, active bool) error {
	query := `UPDATE public.job_scheduler SET is_active = $1 WHERE id = $2;`

	if _, err := conn.Exec(query, active, jobID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
