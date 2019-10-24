package dao

import (
	"strings"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
	"github.com/jmoiron/sqlx"
)

func GetStack(conn *sqlx.DB, stack, key string) (mid model.Stack, err error) {
	const query = `
		SELECT id, "name", convert_from(decrypt(dsn::bytea,$1,'bf'),'SQL_ASCII') dsn
		  FROM public.stack 
		 WHERE lower("name") = $2
		   AND is_active = TRUE
		   AND is_deleted = FALSE
		 LIMIT 1;`

	err = conn.Get(&mid, query, key, stack)
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
