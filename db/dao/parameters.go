package dao

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/gopkgs/db/model"
)

type ParameterType int

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
func GetParameters(conn *sqlx.DB, tenantId int, pType ParameterType) (p model.Parameters, err error) {
	sb := strings.Builder{}
	var a []interface{}

	sb.WriteString(`
		SELECT p.id, p.tenant_id, t.org_id, p."name", p."type", p.value, p.description 
		  FROM public."parameter" p
		  JOIN public.tenant      t ON p.tenant_id = t.id
		 WHERE p.tenant_id = $1 
		   AND p.is_active = true 
		   AND p.is_deleted = false`)

	a = append(a, tenantId)

	if pType != EnumParamNil {
		sb.WriteString(" AND type = $2")
		a = append(a, pType.String())
	}

	err = conn.Select(&p, sb.String(), a...)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Select()")
	}

	return p, nil
}

// GetParameter retorna o parametro informado (paramName) de uma determinada org (tenant_id) do Salesforce
func GetParameter(conn *sqlx.DB, tenantId int, paramName string) (p model.Parameter, err error) {
	query := `
		SELECT p.id, p.tenant_id, t.org_id, p."name", p."type", p.value, p.description 
		  FROM public."parameter" p
		  JOIN public.tenant      t ON p.tenant_id = t.id
		 WHERE p.tenant_id = $1 
		   AND p."name" = $2
		   AND p.is_active = true 
		   AND p.is_deleted = false
		 LIMIT 1;`

	err = conn.Get(&p, query, tenantId, paramName)
	if err != nil {
		return p, errors.Wrap(err, "conn.Get()")
	}

	return p, nil
}

// GetParametersByOrgID retorna os parametros de uma determinada org (orgID) do Salesforce
func GetParametersByOrgID(conn *sqlx.DB, orgID string) (p model.Parameters, err error) {
	query := `
		SELECT p.id, p.tenant_id, t.org_id, p."name", p."type", p.value, p.description 
		  FROM public."parameter" p 
		  JOIN public.tenant      t ON p.tenant_id = t.id
		 WHERE t.org_id = $1
		   AND p.is_active = true 
		   AND p.is_deleted = false;`

	err = conn.Select(&p, query, orgID)
	if err != nil {
		return nil, errors.Wrap(err, "conn.Select()")
	}

	return p, nil
}

// GetParameterByOrgID retorna o parametro informado (paramName) de uma determinada org (orgID) do Salesforce
func GetParameterByOrgID(conn *sqlx.DB, orgID, paramName string) (p model.Parameter, err error) {
	query := `
		SELECT p.id, p.tenant_id, t.org_id, p."name", p."type", p.value, p.description 
		  FROM public."parameter" p 
		  JOIN public.tenant      t ON p.tenant_id = t.id
		 WHERE t.org_id = $1
		   AND p."name" = $2
		   AND p.is_active = true 
		   AND p.is_deleted = false
		 LIMIT 1;`

	err = conn.Get(&p, query, orgID, paramName)
	if err != nil {
		return p, errors.Wrap(err, "conn.Get()")
	}

	return p, nil
}

// UpdateParameter atualiza o parametro de uma determinada org (tenant_id)
func UpdateParameter(conn *sqlx.DB, param model.Parameter) error {
	query := `
		INSERT INTO public."parameter" (tenant_id, "name", "type", value, description) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (tenant_id, "name")
		DO UPDATE SET 
		  value = EXCLUDED.value,
		  updated_at = now();`

	if _, err := conn.Exec(query, param.TenantID, param.Name, param.Type, param.Value, param.Description); err != nil {
		return errors.Wrap(err, "conn.Exec()")
	}

	return nil
}

// UpdateParameters atualiza o parametro de uma determinada org (tenant_id)
func UpdateParameters(conn *sqlx.DB, params []model.Parameter) error {
	if len(params) == 0 {
		return errors.New("no parameters to save")
	}

	query := `
		INSERT INTO public."parameter" (tenant_id, "name", "type", value, description) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (tenant_id, "name")
		DO UPDATE SET 
		  value = EXCLUDED.value,
		  updated_at = now();`

	stmt, err := conn.Preparex(query)
	if err != nil {
		return errors.Wrap(err, "conn.Preparex()")
	}

	for _, p := range params {
		if _, err := stmt.Exec(p.TenantID, p.Name, p.Type, p.Value, p.Description); err != nil {
			return errors.Wrap(err, "stmt.Exec()")
		}
	}

	return nil
}
