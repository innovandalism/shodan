package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)


type ShodanPostgres struct{
	DB *sql.DB
}

func (sp *ShodanPostgres) Connect(connStr string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	sp.DB = db
	return nil
}

func (sp *ShodanPostgres) AddToken(token string, discord_uid string) (int, error) {
	row := sp.DB.QueryRow("INSERT INTO shodan_jwt (access_token, userid) VALUES ($1, $2) RETURNING id", token, discord_uid)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}