package dao

import (
	"net/url"
	"strings"

	"bitbucket.org/everymind/evmd-golib/db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func GetTenant(conn *sqlx.DB, orgID string) (tid int, err error) {
	const query = `
		SELECT id
		FROM public.tenant
		WHERE org_id = $1
		LIMIT 1;`

	row := conn.QueryRow(query, orgID)

	if e := row.Scan(&tid); e != nil {
		err = db.WrapError(e, "row.Scan()")
		return
	}

	return
}

func SaveConfigTenant(conn *sqlx.DB, name, companyID, orgID, instanceUrl, organizationType string, isSandbox bool) (tid int, err error) {
	const query = `
		INSERT INTO public.tenant (id, "name", company_id, org_id, custom_domain, organization_type, is_sandbox) 
		VALUES(fn_next_tenant_id(), $1, $2, $3, $4, $5, $6) 
		ON CONFLICT (org_id) DO 
		UPDATE SET "name" = EXCLUDED."name", custom_domain = EXCLUDED.custom_domain, organization_type = EXCLUDED.organization_type, is_sandbox = EXCLUDED.is_sandbox, updated_at = NOW()
		RETURNING id;`

	var customDomain string
	if len(instanceUrl) > 0 {
		u, err := url.Parse(instanceUrl)
		if err != nil {
			return 0, errors.Wrap(err, "url.Parse()")
		}
		h := strings.Split(u.Hostname(), ".")
		customDomain = h[0]
	}

	err = conn.QueryRowx(query, name, companyID, orgID, customDomain, organizationType, isSandbox).Scan(&tid)
	if err != nil {
		return 0, db.WrapError(err, "conn.QueryRowx()")
	}

	if tid <= 0 {
		err = errors.New("An error has occurred while inserting on 'itgr.execution'")
		return 0, err
	}

	return
}

func SaveBusinessTenant(conn *sqlx.DB, tenantID int, name, orgID string) error {
	const query = `
		INSERT INTO public.tenant (id, "name", org_id) 
		VALUES($1, $2, $3, $4)
		ON CONFLICT (org_id) DO 
		UPDATE SET "name" = EXCLUDED."name", updated_at = NOW()
		RETURNING id;`

	if _, err := conn.Exec(query, tenantID, name, orgID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
