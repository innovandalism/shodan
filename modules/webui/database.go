package webui

import (
	"database/sql"
	"github.com/innovandalism/shodan"
)

type ShodanRole struct{
	ID int
	Guild string
	Role string
	Name string
	Type string
}

func DBGetAvailableRoles(db *sql.DB, guild string) ([]*ShodanRole, error) {
	row, err := db.Query("SELECT id, guild, role, name, type FROM shodan_roles WHERE guild = $1", guild)
	if err != nil {
		return nil, shodan.WrapError(err)
	}
	roles := []*ShodanRole{}
	for row.Next() {
		role := ShodanRole{}
		err = row.Scan(role.ID, role.Guild, role.Role, role.Name, role.Type)
		if err != nil {
			return nil, shodan.WrapError(err)
		}
		roles = append(roles, &role)
	}
	return roles, nil
}