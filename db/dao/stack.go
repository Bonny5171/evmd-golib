package dao

import (
	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
	"github.com/jmoiron/sqlx"
)

func GetStacks(conn *sqlx.DB, stack, key string) (mid model.Stack, err error) {
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