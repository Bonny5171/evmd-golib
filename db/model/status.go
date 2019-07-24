package model

type Status struct {
	ID       int16  `db:"id"`
	TenantID int    `db:"tenant_id"`
	Name     string `db:"name"`
	Type     string `db:"type"`
}

type Statuses []Status

func (ss *Statuses) GetId(n string) (id int16) {
	for _, item := range *(ss) {
		if item.Name == n {
			return item.ID
		}
	}

	return
}
