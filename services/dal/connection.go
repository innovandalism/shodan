package dal

import (
	"database/sql"
	_ "github.com/lib/pq"
)


type Database struct{
	DB *sql.DB
}

func (sp *Database) Connect(connStr string) error {
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