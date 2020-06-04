package model

//Parameter model struct
type Parameter struct {
	ID       int    `db:"id"`
	TenantID int    `db:"tenant_id"`
	OrgID    string `db:"org_id"`
	Name     string `db:"name"`
	Value    string `db:"value"`
	Type     string `db:"type"`
}

//Parameters Slice
type Parameters []Parameter

//ByName Method of Parameters to find a Parameter by Name
func (p *Parameters) ByName(name string) (result string) {
	for _, item := range *(p) {
		if item.Name == name {
			result = item.Value
			break
		}
	}

	return
}
