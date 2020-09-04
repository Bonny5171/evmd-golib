package dao

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

//GetCommunity func
func GetCommunity(conn *sqlx.DB, tid int, cid string) (c model.Community, err error) {
	const query = `
		SELECT id, tenant_id, "name", description, login_url, site_url, path_prefix
		  FROM public.community
		 WHERE tenant_id = $1
		   AND id = $2
		   AND is_active = TRUE
		   AND is_deleted = FALSE
		 LIMIT 1;`

	if e := conn.QueryRowx(query, tid, cid).StructScan(&c); e != nil {
		err = db.WrapError(e, "conn.QueryRowx()")
		return
	}

	return c, nil
}

//GetCommunities func
func GetCommunities(conn *sqlx.DB, tid int) (c model.Communities, err error) {
	const query = `
		SELECT id, tenant_id, "name", description, login_url, site_url, path_prefix 
		  FROM public.community
		 WHERE tenant_id = $1
		   AND is_active = TRUE
		   AND is_deleted = FALSE;`

	err = conn.Select(&c, query, tid)
	if err != nil {
		err = db.WrapError(err, "conn.Select()")
		return
	}

	return c, nil
}

//SaveCommunity func
func SaveCommunity(conn *sqlx.DB, community model.Community) (err error) {
	const query = `
		INSERT INTO public.community (id, tenant_id, "name", description, login_url, site_url, path_prefix) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id, tenant_id) DO UPDATE 
		SET "name"      = EXCLUDED."name", 
		    description = EXCLUDED.description, 
		    login_url   = EXCLUDED.login_url, 
			site_url    = EXCLUDED.siteurl,
			path_prefix = EXCLUDED.path_prefix,
		    updated_at  = now();`

	if _, err = conn.Exec(query, community.ID, community.TenantID, community.Name, community.Description, community.LoginURL, community.SiteURL, community.PathPrefix); err != nil {
		err = db.WrapError(err, "conn.Exec()")
		return
	}

	return nil
}

//SaveCommunities func
func SaveCommunities(conn *sqlx.DB, communities model.Communities) (err error) {
	query := `INSERT INTO public.community (id, tenant_id, "name", description, login_url, site_url, path_prefix) VALUES`

	i := 1
	vals := []interface{}{}
	for _, row := range communities {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d),", i, i+1, i+2, i+3, i+4, i+5, i+6)
		vals = append(vals, row.ID, row.TenantID, row.Name, row.Description, row.LoginURL, row.SiteURL, row.PathPrefix)
		i += 7
	}
	query = query[0 : len(query)-1]

	query += `
	     ON CONFLICT (id, tenant_id) DO UPDATE 
		SET "name"      = EXCLUDED."name", 
		    description = EXCLUDED.description, 
		    login_url   = EXCLUDED.login_url, 
			site_url    = EXCLUDED.site_url,
			path_prefix = EXCLUDED.path_prefix,
			updated_at  = now();`

	stmt, e := conn.Prepare(query)
	if e != nil {
		err = db.WrapError(e, "conn.Prepare()")
		return
	}

	if _, e := stmt.Exec(vals...); e != nil {
		err = db.WrapError(e, "stmt.Exec()")
		return
	}

	return nil
}
