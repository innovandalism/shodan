package dal

func (sp *Database) AddToken(token string, discord_uid string) (int, error) {
	row := sp.DB.QueryRow("INSERT INTO shodan_jwt (access_token, userid) VALUES ($1, $2) ON CONFLICT (userid) DO UPDATE SET access_token = $1 RETURNING id", token, discord_uid)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (sp *Database) GetToken(id int) (string, error) {
	row := sp.DB.QueryRow("SELECT access_token FROM shodan_jwt WHERE id=$1", id)
	var token string
	err := row.Scan(&token)
	if err != nil {
		return "",err
	}
	return token,nil
}