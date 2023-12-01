package Connexion

import "database/sql"

func checkEmailExists(db *sql.DB, email string) (bool, error) {
	row := db.QueryRow("SELECT COUNT(*) FROM utilisateurs WHERE Email = ?", email)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
