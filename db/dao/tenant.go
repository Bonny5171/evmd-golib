package dao

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

//GetTenantID func
func GetTenantID(conn *sqlx.DB, orgID, alias string) (tid int, err error) {
	const query = `
		SELECT id
		FROM public.tenant
		WHERE LEFT(org_id, 15) = LEFT($1, 15)
		AND alias = $2
		AND is_active = TRUE
		AND is_deleted = FALSE
		LIMIT 1;`

	row := conn.QueryRow(query, orgID, alias)

	if e := row.Scan(&tid); e != nil {
		err = db.WrapError(e, "row.Scan()")
		return
	}

	return
}

//GetTenantByID func
func GetTenantByID(conn *sqlx.DB, tid int) (tenant model.Tenant, err error) {
	const query = `
		SELECT id, company_id, name, org_id, organization_type, custom_domain, is_sandbox, is_active, created_at, updated_at, is_deleted, deleted_at, last_modified_by_id, is_cloned, sf_client_id, sf_client_secret, sf_callback_token_url
		FROM public.tenant
		WHERE id = $1
		AND is_active = TRUE
		AND is_deleted = FALSE
		LIMIT 1;
	`

	row := conn.QueryRowx(query, tid)
	if e := row.StructScan(&tenant); e != nil {
		err = db.WrapError(e, "row.StructScan()")
		return
	}

	return
}

//GetTenantByCustomDomain func
func GetTenantByCustomDomain(conn *sqlx.DB, domain string) (tenant model.Tenant, err error) {
	const query = `
		SELECT id, company_id, name, org_id, organization_type, custom_domain, is_sandbox, is_active, created_at, updated_at, is_deleted, deleted_at, last_modified_by_id, is_cloned, sf_client_id, sf_client_secret, sf_callback_token_url
		FROM public.tenant
		WHERE custom_domain = $1
		AND is_active = TRUE
		AND is_deleted = FALSE
		LIMIT 1;
	`

	row := conn.QueryRowx(query, domain)
	if e := row.StructScan(&tenant); e != nil {
		err = db.WrapError(e, "row.StructScan()")
		return
	}
	return
}

//GetTenantByCommunityUrl func
func GetTenantByCommunityUrl(conn *sqlx.DB, communityURL string) (tenant model.Tenant, err error) {
	const query = `
		SELECT id, company_id, name, org_id, organization_type, custom_domain, is_sandbox, is_active, created_at, updated_at, is_deleted, deleted_at, last_modified_by_id, is_cloned, sf_client_id, sf_client_secret, sf_callback_token_url
		FROM public.tenant
		WHERE community_url = $1
		AND is_active = TRUE
		AND is_deleted = FALSE
		LIMIT 1;
	`

	row := conn.QueryRowx(query, communityURL)
	if e := row.StructScan(&tenant); e != nil {
		err = db.WrapError(e, "row.StructScan()")
		return
	}
	return
}

//GetTenant func
func GetTenant(conn *sqlx.DB, orgID string) (tenant model.Tenant, err error) {
	const query = `
		SELECT id, company_id, name, org_id, organization_type, custom_domain, is_sandbox, is_active, is_deleted
		FROM public.tenant
		WHERE LOWER(LEFT(org_id, 15)) = LOWER(LEFT($1, 15)) AND is_deleted = FALSE
		ORDER BY id DESC
		LIMIT 1;`

	if e := conn.QueryRowx(query, orgID).StructScan(&tenant); e != nil {
		err = db.WrapError(e, "conn.QueryRowx().StructScan()")
		return
	}

	return
}

//SaveConfigTenant func
func SaveConfigTenant(conn *sqlx.DB, name, companyID, orgID, instanceUrl, organizationType, userID string, isSandbox bool) (tid int, err error) {
	const query = `
		INSERT INTO public.tenant (id, "name", company_id, org_id, custom_domain, organization_type, is_sandbox, last_modified_by_id, is_active) 
		VALUES(fn_next_tenant_id(), $1, $2, $3, $4, $5, $6, $7, true) 
		ON CONFLICT (org_id, is_active) DO 
		UPDATE SET "name" = EXCLUDED."name", custom_domain = EXCLUDED.custom_domain, organization_type = EXCLUDED.organization_type, is_sandbox = EXCLUDED.is_sandbox, updated_at = NOW()
		RETURNING id;`

	var customDomain string
	if len(instanceUrl) > 0 {
		u, err := url.Parse(instanceUrl)
		if err != nil {
			return 0, fmt.Errorf("url.Parse(): %w", err)
		}
		h := strings.Split(u.Hostname(), ".")
		customDomain = h[0]
	}

	err = conn.QueryRowx(query, name, companyID, orgID, customDomain, organizationType, isSandbox, userID).Scan(&tid)
	if err != nil {
		return 0, db.WrapError(err, "conn.QueryRowx()")
	}

	if tid <= 0 {
		err = errors.New("An error has occurred while inserting on 'itgr.execution'")
		return 0, err
	}

	return
}

//SaveConfigTenantTx func
func SaveConfigTenantTx(conn *sqlx.Tx, name, companyID, orgID, instanceURL, organizationType, userID string, isSandbox bool, clientID, clientSecret, callbackURL, alias string, isCloned bool, clonedFrom int) (tid int, err error) {
	const query = `
		INSERT INTO public.tenant (id, "name", company_id, org_id, custom_domain, organization_type, is_sandbox, last_modified_by_id, is_active, sf_client_id, sf_client_secret, sf_callback_token_url, alias, is_cloned, cloned_from) 
		VALUES(fn_next_tenant_id(), $1, $2, $3, $4, $5, $6, $7, true, $8, $9, $10, $11, $12, $13) 
		RETURNING id;`

	var customDomain string
	if len(instanceURL) > 0 {
		u, err := url.Parse(instanceURL)
		if err != nil {
			return 0, fmt.Errorf("url.Parse(): %w", err)
		}
		h := strings.Split(u.Hostname(), ".")
		customDomain = h[0]
	}

	err = conn.QueryRowx(query, name, companyID, orgID, customDomain, organizationType, isSandbox, userID, clientID, clientSecret, callbackURL, alias, isCloned, clonedFrom).Scan(&tid)
	if err != nil {
		return 0, db.WrapError(err, "conn.QueryRowx()")
	}

	if tid <= 0 {
		err = errors.New("An error has occurred while inserting on 'itgr.execution'")
		return 0, err
	}

	return
}

//SaveConfigTenantWithIDTx func
func SaveConfigTenantWithIDTx(conn *sqlx.Tx, name, companyID, orgID, instanceURL, organizationType, userID string, isSandbox bool, clientID, clientSecret, callbackURL, alias string, isCloned bool, clonedFrom, tenantID int) (tid int, err error) {
	const query = `
		INSERT INTO public.tenant (id, "name", company_id, org_id, custom_domain, organization_type, is_sandbox, last_modified_by_id, is_active, sf_client_id, sf_client_secret, sf_callback_token_url, alias, is_cloned, cloned_from) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, true, $9, $10, $11, $12, $13, $14) 
		RETURNING id;`

	var customDomain string
	if len(instanceURL) > 0 {
		u, err := url.Parse(instanceURL)
		if err != nil {
			return 0, fmt.Errorf("url.Parse(): %w", err)
		}
		h := strings.Split(u.Hostname(), ".")
		customDomain = h[0]
	}

	err = conn.QueryRowx(query, tenantID, name, companyID, orgID, customDomain, organizationType, isSandbox, userID, clientID, clientSecret, callbackURL, alias, isCloned, clonedFrom).Scan(&tid)
	if err != nil {
		return 0, db.WrapError(err, "conn.QueryRowx()")
	}

	if tid <= 0 {
		err = errors.New("An error has occurred while inserting on 'itgr.execution'")
		return 0, err
	}

	return
}

//SaveBusinessTenant func
func SaveBusinessTenant(conn *sqlx.DB, tenantID int, name, orgID, userID string) error {
	const query = `
		INSERT INTO public.tenant (id, "name", org_id, last_modified_by_id, is_active) 
		VALUES($1, $2, $3, $4 $5, true)
		ON CONFLICT (org_id, is_active) DO 
		UPDATE SET "name" = EXCLUDED."name", updated_at = NOW()
		RETURNING id;`

	if _, err := conn.Exec(query, tenantID, name, orgID, userID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

//SaveBusinessTenantTx func
func SaveBusinessTenantTx(conn *sqlx.Tx, tenantID int, name, orgID, userID string) error {
	const query = `
		INSERT INTO public.tenant (id, "name", org_id, last_modified_by_id, is_active) 
		VALUES($1, $2, $3, $4 $5, true)
		ON CONFLICT (org_id, is_active) DO 
		UPDATE SET "name" = EXCLUDED."name", updated_at = NOW()
		RETURNING id;`

	if _, err := conn.Exec(query, tenantID, name, orgID, userID); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

//CheckTenantAvailability func
func CheckTenantAvailability(conn *sqlx.DB, orgID string) (inUse bool, err error) {
	query := `SELECT CASE WHEN count(*) > 0 THEN true ELSE false END AS in_use FROM public.tenant WHERE LOWER(LEFT(org_id, 15)) = LOWER(LEFT($1, 15)) AND is_deleted = FALSE;`

	if e := conn.QueryRow(query, orgID).Scan(&inUse); e != nil {
		err = db.WrapError(e, "row.Scan()")
		return
	}

	return
}

//DisableTenant func
func DisableTenant(conn *sqlx.DB, tenantID int, userID string) error {
	query := "UPDATE public.tenant SET is_active = FALSE, updated_at = $3, last_modified_by_id = $2 WHERE id = $1;"

	if _, err := conn.Exec(query, tenantID, userID, time.Now()); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

//DisableTenantTx func
func DisableTenantTx(conn *sqlx.Tx, tenantID int, userID string) error {
	query := "UPDATE public.tenant SET is_active = FALSE, updated_at = $3, last_modified_by_id = $2 WHERE id = $1;"

	if _, err := conn.Exec(query, tenantID, userID, time.Now()); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

//DeleteTenant func
func DeleteTenant(conn *sqlx.DB, tenantID int, userID string) error {
	query := "UPDATE public.tenant SET is_active = FALSE, updated_at = $3, is_deleted = TRUE, deleted_at = $3, last_modified_by_id = $2 WHERE id = $1;"

	if _, err := conn.Exec(query, tenantID, userID, time.Now()); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

//DeleteTenantTx func
func DeleteTenantTx(conn *sqlx.Tx, tenantID int, userID string) error {
	query := "UPDATE public.tenant SET is_active = FALSE, updated_at = $3, is_deleted = TRUE, deleted_at = $3, last_modified_by_id = $2 WHERE id = $1;"

	if _, err := conn.Exec(query, tenantID, userID, time.Now()); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

//EnableTenant func
func EnableTenant(conn *sqlx.DB, tenantID int, userID string) error {
	query := "UPDATE public.tenant SET is_active = TRUE, updated_at = $3, last_modified_by_id = $2 WHERE id = $1;"

	if _, err := conn.Exec(query, tenantID, userID, time.Now()); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

//EnableTenantTx func
func EnableTenantTx(conn *sqlx.Tx, tenantID int, userID string) error {
	query := "UPDATE public.tenant SET is_active = TRUE, updated_at = $3, last_modified_by_id = $2 WHERE id = $1;"

	if _, err := conn.Exec(query, tenantID, userID, time.Now()); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}
