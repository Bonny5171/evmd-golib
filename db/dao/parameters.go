package dao

import (
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/db/model"
)

//ParameterType type
type ParameterType int

//EnumParam types
const (
	EnumParamNil ParameterType = iota
	EnumParamA
	EnumParamB
	EnumParamDS
	EnumParamN
	EnumParamO
	EnumParamS
	EnumParamDJ
)

func (t ParameterType) String() string {
	return [...]string{"", "a", "b", "ds", "n", "o", "s", "dj"}[t]
}

// GetParameters retorna os parametros de uma determinada org (tenant_id) do Salesforce
func GetParameters(conn *sqlx.DB, tenantID int, pType ParameterType) (p model.Parameters, err error) {
	sb := strings.Builder{}
	var a []interface{}

	sb.WriteString(`
		SELECT p.id, p.tenant_id, t.org_id, p."name", p."type", p.value 
		  FROM public."parameter" p
		  JOIN public.tenant      t ON p.tenant_id = t.id
		 WHERE p.tenant_id = $1 
		   AND p.is_active = true 
		   AND p.is_deleted = false`)

	a = append(a, tenantID)

	if pType != EnumParamNil {
		sb.WriteString(" AND type = $2")
		a = append(a, pType.String())
	}

	err = conn.Select(&p, sb.String(), a...)
	if err != nil {
		return nil, db.WrapError(err, "conn.Select()")
	}

	return p, nil
}

// GetParameterByOrgID retorna o parametro informado (paramName) de uma determinada org (orgID) do Salesforce
func GetParameterByOrgID(conn *sqlx.DB, orgID, paramName string) (p model.Parameter, err error) {
	query := `
		SELECT p.id, p.tenant_id, t.org_id, p."name", p."type", p.value 
		  FROM public."parameter" p 
		  JOIN public.tenant      t ON p.tenant_id = t.id
		 WHERE t.org_id = $1
		   AND p."name" = $2
		   AND p.is_active = true 
		   AND p.is_deleted = false
		 LIMIT 1;`

	err = conn.Get(&p, query, orgID, paramName)
	if err != nil {
		return p, db.WrapError(err, "conn.Get()")
	}

	return p, nil
}

// UpdateParameter atualiza o parametro de uma determinada org (tenant_id)
func UpdateParameter(conn *sqlx.DB, param model.Parameter) error {
	query := `
		INSERT INTO public."parameter" (tenant_id, "name", value) 
		VALUES ($1, $2, $3)
		ON CONFLICT (tenant_id, "name")
		DO UPDATE SET 
		  value = EXCLUDED.value,
		  updated_at = now();`

	if _, err := conn.Exec(query, param.TenantID, param.Name, param.Value); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// UpdateParameters atualiza o parametro de uma determinada org (tenant_id)
func UpdateParameters(conn *sqlx.DB, params []model.Parameter) error {
	if len(params) == 0 {
		return errors.New("no parameters to save")
	}

	query := `
		INSERT INTO public."parameter" (tenant_id, "name", value) 
		VALUES ($1, $2, $3)
		ON CONFLICT (tenant_id, "name")
		DO UPDATE SET 
		  value = EXCLUDED.value,
		  updated_at = now();`

	stmt, err := conn.Preparex(query)
	if err != nil {
		return db.WrapError(err, "conn.Preparex()")
	}

	for _, p := range params {
		if _, err := stmt.Exec(p.TenantID, p.Name, p.Value); err != nil {
			return db.WrapError(err, "stmt.Exec()")
		}
	}

	return nil
}

// UpdateParametersTx atualiza o parametro de uma determinada org (tenant_id)
func UpdateParametersTx(conn *sqlx.Tx, params []model.Parameter) error {
	if len(params) == 0 {
		return errors.New("no parameters to save")
	}

	query := `
		INSERT INTO public."parameter" (tenant_id, "name", value) 
		VALUES ($1, $2, $3)
		ON CONFLICT (tenant_id, "name")
		DO UPDATE SET 
		  value = EXCLUDED.value,
		  updated_at = now();`

	stmt, err := conn.Preparex(query)
	if err != nil {
		return db.WrapError(err, "conn.Preparex()")
	}

	for _, p := range params {
		if _, err := stmt.Exec(p.TenantID, p.Name, p.Value); err != nil {
			return db.WrapError(err, "stmt.Exec()")
		}
	}

	return nil
}

// UpdateStackParameter atualiza o parametro de uma determinada org (tenant_id)
func UpdateStackParameter(conn *sqlx.DB, param model.Parameter) error {
	query := `
		INSERT INTO public."parameter" (id, tenant_id, value, "type") 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id, tenant_id, app_id, record_type_id)
		DO UPDATE SET 
		  value  = EXCLUDED.value,
		  "type" = EXCLUDED."type",
		  updated_at = now();`

	if _, err := conn.Exec(query, param.Name, param.TenantID, param.Value, param.Type); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// UpdateStackParameterTx atualiza o parametro de uma determinada org (tenant_id)
func UpdateStackParameterTx(conn *sqlx.Tx, param model.Parameter) error {
	query := `
		INSERT INTO public."parameter" (id, tenant_id, value, "type") 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id, tenant_id, app_id, record_type_id)
		DO UPDATE SET 
		  value  = EXCLUDED.value,
		  "type" = EXCLUDED."type",
		  updated_at = now();`

	if _, err := conn.Exec(query, param.Name, param.TenantID, param.Value, param.Type); err != nil {
		return db.WrapError(err, "conn.Exec()")
	}

	return nil
}

// UpdateStackParameters atualiza o parametro de uma determinada org (tenant_id)
func UpdateStackParameters(conn *sqlx.DB, params []model.Parameter) error {
	if len(params) == 0 {
		return errors.New("no parameters to save")
	}

	query := `
		INSERT INTO public."parameter" (id, tenant_id, value, "type") 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id, tenant_id, app_id, record_type_id)
		DO UPDATE SET 
		  value = EXCLUDED.value,
		  "type" = EXCLUDED."type",
		  updated_at = now();`

	stmt, err := conn.Preparex(query)
	if err != nil {
		return db.WrapError(err, "conn.Preparex()")
	}

	for _, p := range params {
		if _, err := stmt.Exec(p.Name, p.TenantID, p.Value, p.Type); err != nil {
			return db.WrapError(err, "stmt.Exec()")
		}
	}

	return nil
}

// UpdateStackParametersTx atualiza o parametro de uma determinada org (tenant_id)
func UpdateStackParametersTx(conn *sqlx.Tx, params []model.Parameter) error {
	if len(params) == 0 {
		return errors.New("no parameters to save")
	}

	query := `
		INSERT INTO public."parameter" (id, tenant_id, value, "type") 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id, tenant_id, app_id, record_type_id)
		DO UPDATE SET 
		  value = EXCLUDED.value,
		  "type" = EXCLUDED."type",
		  updated_at = now();`

	stmt, err := conn.Preparex(query)
	if err != nil {
		return db.WrapError(err, "conn.Preparex()")
	}

	for _, p := range params {
		if _, err := stmt.Exec(p.Name, p.TenantID, p.Value, p.Type); err != nil {
			return db.WrapError(err, "stmt.Exec()")
		}
	}

	return nil
}
