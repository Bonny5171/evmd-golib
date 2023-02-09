package dao

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/CognyHub/evmd-golib/db"
	"github.com/CognyHub/evmd-golib/db/model"
)

// GetSchedules retorna todos os 'jobs' agendados que deverão ser executadas
func GetSchedules(conn *sqlx.DB, tenantID, stackID int) (s []model.JobScheduler, err error) {
	query := `
	  SELECT j.id, t.org_id, j.tenant_id, t."name" AS tenant_name, j.stack_id, j.job_name, j.function_name, j.queue, 
			 j.cron, j.parameters, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.description, 
			 j.is_active, j.is_deleted, j.appengine_name
	    FROM public.job_scheduler j
	   INNER JOIN public.tenant   t ON j.tenant_id = t.id
	   WHERE j.tenant_id = $1
		 AND j.stack_id = $2
		 AND t.is_deleted = FALSE
	   ORDER BY j.id;`

	err = conn.Select(&s, query, tenantID, stackID)
	if err != nil {
		return nil, db.WrapError(err, "db.Conn.Select()")
	}

	return s, nil
}

// GetSchedulesByOrg retorna todos os 'jobs' agendados que deverão ser executadas
func GetSchedulesByOrg(conn *sqlx.DB, orgID string, stackID, tenantID int) (s []model.JobScheduler, err error) {
	query := `
	  SELECT j.id, t.org_id, j.tenant_id, t."name" AS tenant_name, j.stack_id, j.job_name, j.function_name, j.queue, 
			 j.cron, j.parameters, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.description, 
			 j.is_active, j.is_deleted, j.appengine_name
	    FROM public.job_scheduler j
	   INNER JOIN public.tenant   t ON j.tenant_id = t.id
	   WHERE t.org_id = $1
		 AND j.stack_id = $2
		 AND t.is_deleted = FALSE
		 AND j.tenant_id = $3		 
	   ORDER BY j.id;`

	err = conn.Select(&s, query, orgID, stackID, tenantID)
	if err != nil {
		return nil, db.WrapError(err, "db.Conn.Select()")
	}

	return s, nil
}

// GetJob retorna os dados de um 'job'
func GetJob(conn *sqlx.DB, tenantID, stackID int, name string) (s model.JobScheduler, err error) {
	query := `
	  SELECT j.id, t.org_id, j.tenant_id, t."name" AS tenant_name, j.stack_id, j.job_name, j.function_name, j.queue, 
			 j.cron, j.parameters, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.description, 
			 j.is_active, j.is_deleted, j.appengine_name
	    FROM public.job_scheduler j
	   INNER JOIN public.tenant   t ON j.tenant_id = t.id
	   WHERE j.tenant_id = $1
	     AND j.stack_id = $2
		 AND j.job_name = $3
		 AND t.is_active = TRUE
	   LIMIT 1;`

	err = conn.Get(&s, query, tenantID, stackID, name)
	if err != nil {
		return s, db.WrapError(err, "conn.Get()")
	}

	return s, nil
}

// GetJobByFuncQueue retorna os dados de um 'job' por function name e queue
func GetJobByFuncQueue(conn *sqlx.DB, tenantID int, stackName, funcName, queue, jobName string) (s model.JobScheduler, err error) {
	query := `
	  SELECT j.id, t.org_id, j.tenant_id, t."name" AS tenant_name, j.stack_id, j.job_name, j.function_name, j.queue, 
			 j.cron, j.parameters, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.description, 
			 j.is_active, j.is_deleted, j.appengine_name
	    FROM public.job_scheduler j
	   INNER JOIN public.tenant   t ON j.tenant_id = t.id
	   INNER JOIN public.stack    s ON j.stack_id = s.id
	   WHERE j.tenant_id = $1
	     AND lower(s."name") = $2
		 AND lower(j.function_name) = $3
		 AND lower(j.queue) = $4
		 AND lower(j.job_name) = $5
		 AND t.is_active = TRUE
	   LIMIT 1;`

	err = conn.Get(&s, query, tenantID, strings.ToLower(stackName), strings.ToLower(funcName), strings.ToLower(queue), strings.ToLower(jobName))
	if err != nil {
		return s, db.WrapError(err, "conn.Get()")
	}

	return s, nil
}

// GetJobByID retorna os dados de um 'job'
func GetJobByID(conn *sqlx.DB, jobID int64) (s model.JobScheduler, err error) {
	query := `
	  SELECT j.id, t.org_id, j.tenant_id, t."name" AS tenant_name, j.stack_id, j.job_name, j.function_name, j.queue, 
	         j.cron, j.parameters, j.retry, j.allows_concurrency, j.allows_schedule, j.schedule_time, j.description, 
	         j.is_active, j.is_deleted, j.appengine_name
	    FROM public.job_scheduler j
	   INNER JOIN public.tenant   t ON j.tenant_id = t.id
	   WHERE j.id = $1
	     AND t.is_active = TRUE
	   LIMIT 1;`

	err = conn.Get(&s, query, jobID)
	if err != nil {
		return s, db.WrapError(err, "conn.Get()")
	}

	return s, nil
}

// SetCronJobSchedule func
func SetCronJobSchedule(conn *sqlx.DB, jobID int64, cronexpr string, tenantID int) error {
	query := `UPDATE public.job_scheduler SET cron = $1 WHERE id = $2 AND tenant_id = $3;`

	if _, err := conn.Exec(query, cronexpr, jobID, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// ActiveJobSchedule func
func ActiveJobSchedule(conn *sqlx.DB, jobID int64, active bool, tenantID int) error {
	query := `UPDATE public.job_scheduler SET is_active = $1 WHERE id = $2 AND tenant_id = $3;`

	if _, err := conn.Exec(query, active, jobID, tenantID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
