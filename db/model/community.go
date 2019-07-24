package model

type Community struct {
	ID          string      `db:"id"`
	TenantID    int         `db:"tenant_id"`
	Name        string      `db:"name"`
	Description interface{} `db:"description"`
	LoginURL    string      `db:"login_url"`
	SiteURL     string      `db:"site_url"`
	PathPrefix  string      `db:"path_prefix"`
}

type Communities []Community
