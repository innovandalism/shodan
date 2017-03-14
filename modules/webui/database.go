package webui

import (
	"database/sql"
	"github.com/innovandalism/shodan"
)

// A ShodanRole represents a managed discord role
type ShodanRole struct {
	ID    int
	Guild string
	Role  string
	Name  string
	Type  string
}

// DBGetAvailableRoles returns the assigned roles for a guild
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
