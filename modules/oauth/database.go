package oauth

import "database/sql"

// DBAddToken adds a token to the database
func DBAddToken(db *sql.DB, token string, discord_uid string) (int, error) {
	row := db.QueryRow("INSERT INTO shodan_jwt (access_token, userid) VALUES ($1, $2) ON CONFLICT (userid) DO UPDATE SET access_token = $1 RETURNING id", token, discord_uid)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DBGetToken returns the token based on the user ID
func DBGetToken(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT access_token FROM shodan_jwt WHERE id=$1", id)
	var token string
	err := row.Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
