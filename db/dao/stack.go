package dao

import (
	"strings"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
	"github.com/jmoiron/sqlx"
)

type TenantType int

const (
	EnumTenentJob TenantType = iota
	EnumTenentAPI
	EnumTenentDebug
)

func GetStack(conn *sqlx.DB, stack string, tenantType TenantType) (mid model.Stack, err error) {
	query := strings.Builder{}
	query.WriteString(`
		SELECT s.id, s."name", d.string AS dsn 
		  FROM public.stack   s
		 INNER JOIN public.dsn d ON s.id = d.stack_id
		 WHERE s.is_active = TRUE 
		   AND s.is_deleted = FALSE`)

	switch tenantType {
	case EnumTenentJob:
		query.WriteString(` AND upper(d."type") = 'JOB'`)
	case EnumTenentAPI:
		query.WriteString(` AND upper(d."type") = 'API'`)
	case EnumTenentDebug:
		query.WriteString(` AND upper(d."type") = 'DEBUG'`)
	}

	query.WriteString(` AND upper(s."name") = $1 LIMIT 1;`)

	err = conn.Get(&mid, query.String(), strings.ToUpper(stack))
	if err != nil {
		return mid, db.WrapError(err, "conn.Get()")
	}

	return mid, nil
}

func GetAllStacks(conn *sqlx.DB, tenantType TenantType) (mid []model.Stack, err error) {
	query := strings.Builder{}
	query.WriteString(`
		SELECT s.id, s."name", d.string AS dsn 
		  FROM public.stack   s
		 INNER JOIN public.dsn d ON s.id = d.stack_id
		 WHERE s.is_active = TRUE 
		   AND s.is_deleted = FALSE`)

	switch tenantType {
	case EnumTenentJob:
		query.WriteString(` AND upper(d."type") = 'JOB'`)
	case EnumTenentAPI:
		query.WriteString(` AND upper(d."type") = 'API'`)
	case EnumTenentDebug:
		query.WriteString(` AND upper(d."type") = 'DEBUG'`)
	}

	if err = conn.Select(&mid, query.String()); err != nil {
		return mid, db.WrapError(err, "conn.Select()")
	}

	return mid, nil
}
