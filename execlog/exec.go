package execlog

import (
	"encoding/json"

	"go.uber.org/multierr"

	"bitbucket.org/everymind/gopkgs/db/dao"
	"bitbucket.org/everymind/gopkgs/db/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Exec struct {
	ID             int64
	JFID           string
	JobSchedulerID int64
	TenantID       int
	SchemaID       int
	StatusID       int16
	Connection     *sqlx.DB
	StatusList     model.Statuses
}

func NewExec(conn *sqlx.DB, jfid string, jsid int64, tid, sid int, st dao.StatusType) (exe Exec, err error) {
	exe.JFID = jfid
	exe.JobSchedulerID = jsid
	exe.TenantID = tid
	exe.SchemaID = sid
	exe.Connection = conn
	exe.StatusList, err = dao.GetStatuses(conn, tid, st)
	if err != nil {
		return exe, errors.Wrap(err, "dao.GetStatuses()")
	}
	return
}

func (e *Exec) LogExecution(s dao.Status) {
	e.log(s, nil)
}

func (e *Exec) LogError(err error) error {
	e.log(dao.EnumStatusExecError, err)
	return err
}

func (e *Exec) LogStackErrors(errs []error) error {
	e.logStack(dao.EnumStatusExecError, errs)

	var err error

	for _, e := range errs {
		err = multierr.Append(err, e)
	}

	return err
}

func (e *Exec) log(s dao.Status, r error) error {
	sn := s.String()
	e.StatusID = e.StatusList.GetId(sn)

	obj := model.Execution{
		ID:             e.ID,
		JobSchedulerID: e.JobSchedulerID,
		TenantID:       e.TenantID,
		StatusID:       e.StatusID,
	}

	if len(e.JFID) > 0 {
		obj.JobFaktoryID.Valid = true
		obj.JobFaktoryID.String = e.JFID
	}

	if e.SchemaID > 0 {
		obj.SchemaID.Valid = true
		obj.SchemaID.Int64 = int64(e.SchemaID)
	}

	if r != nil {
		nerr := &struct {
			Type    string
			Details string
		}{
			Type:    "Error",
			Details: r.Error(),
		}
		jerr, err := json.MarshalIndent(nerr, "", "    ")
		if err != nil {
			return err
		}
		obj.DocMetaData = jerr
	}

	if e.ID == 0 {
		id, err := dao.InsertExecution(e.Connection, obj)
		if err != nil {
			return errors.Wrap(err, "dao.InsertExecution()")
		}

		e.ID = id
	} else {
		err := dao.UpdateExecution(e.Connection, obj)
		if err != nil {
			return errors.Wrap(err, "dao.UpdateExecution()")
		}
	}

	return nil
}

func (e *Exec) logStack(s dao.Status, r []error) error {
	sn := s.String()
	e.StatusID = e.StatusList.GetId(sn)

	obj := model.Execution{
		ID:             e.ID,
		JobSchedulerID: e.JobSchedulerID,
		TenantID:       e.TenantID,
		StatusID:       e.StatusID,
	}

	if len(e.JFID) > 0 {
		obj.JobFaktoryID.Valid = true
		obj.JobFaktoryID.String = e.JFID
	}

	if e.SchemaID > 0 {
		obj.SchemaID.Valid = true
		obj.SchemaID.Int64 = int64(e.SchemaID)
	}

	if len(r) > 0 {
		nerr := &struct {
			Type    string
			Details []string
		}{
			Type: "Errors",
		}

		for _, i := range r {
			nerr.Details = append(nerr.Details, i.Error())
		}

		jerr, err := json.MarshalIndent(nerr, "", "    ")
		if err != nil {
			return err
		}
		obj.DocMetaData = jerr
	}

	if e.ID == 0 {
		id, err := dao.InsertExecution(e.Connection, obj)
		if err != nil {
			return errors.Wrap(err, "dao.InsertExecution()")
		}

		e.ID = id
	} else {
		err := dao.UpdateExecution(e.Connection, obj)
		if err != nil {
			return errors.Wrap(err, "dao.UpdateExecution()")
		}
	}

	return nil
}
