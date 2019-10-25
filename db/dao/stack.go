package dao

import (
	"fmt"
	"strings"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
	"github.com/jmoiron/sqlx"
)

func GetStack(conn *sqlx.DB, stack, key string, debug bool) (mid model.Stack, err error) {
	query := strings.Builder{}
	query.WriteString("SELECT id, \"name\", convert_from(decrypt(dsn::bytea,$1,'bf'),'SQL_ASCII') dsn FROM public.stack WHERE is_deleted = FALSE")

	if debug {
		query.WriteString(fmt.Sprintf(" AND is_active = FALSE AND lower(\"name\") LIKE ('debug:%s%%')", stack))
	} else {
		query.WriteString(fmt.Sprintf(" AND is_active = TRUE AND lower(\"name\") = '%s')", stack))
	}

	query.WriteString(" LIMIT 1;")

	q := query.String()
	err = conn.Get(&mid, q, key)
	if err != nil {
		return mid, db.WrapError(err, "conn.Get()")
	}

	return mid, nil
}

func GetAllStacks(conn *sqlx.DB, key string, debug bool) (mid []model.Stack, err error) {
	query := strings.Builder{}
	query.WriteString(`SELECT id, "name", convert_from(decrypt(dsn::bytea,$1,'bf'),'SQL_ASCII') dsn FROM public.stack WHERE is_deleted = FALSE`)

	if debug {
		query.WriteString(` AND is_active = FALSE AND lower("name") LIKE ('debug:%');`)
	} else {
		query.WriteString(` AND is_active = TRUE;`)
	}

	if err = conn.Select(&mid, query.String(), key); err != nil {
		return mid, db.WrapError(err, "conn.Select()")
	}

	return mid, nil
}
