package Menu

import (
	"Forum/Fonctions/Connexion"
	"database/sql"
	"encoding/base64"
	"html/template"
	"math/rand"
	"net/http"
)

func generateToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}

// ...

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	session, err := Connexion.Store.Get(r, Connexion.SessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if username, ok := session.Values["username"].(string); ok {
		tmpl, err := template.ParseFiles("templates/html/Menu/profile.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var storedToken string
		err = db.QueryRow("SELECT reset_token FROM utilisateurs WHERE username = ?", username).Scan(&storedToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Stocker le nouveau token dans la base de données
		nextResetToken, err := generateToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Mettre à jour le token dans la base de données
		_, err = db.Exec("UPDATE utilisateurs SET reset_token = ? WHERE username = ?", nextResetToken, username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var dateInscription interface{}
		err = db.QueryRow("SELECT date_inscription FROM utilisateurs WHERE username = ?", username).Scan(&dateInscription)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			DateInscription interface{}
			Token           string
			Username        string
			IsLoggedIn      bool
		}{
			DateInscription: dateInscription,
			Username:        username,
			IsLoggedIn:      true,
			Token:           nextResetToken,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	tmpl, err := template.ParseFiles("templates/html/Menu/profile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
