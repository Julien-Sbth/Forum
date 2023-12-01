package Connexion

import "database/sql"

func UserExists(db *sql.DB, username string) (bool, error) {
	row := db.QueryRow("SELECT COUNT(*) FROM utilisateurs WHERE Nom = ?", username)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
