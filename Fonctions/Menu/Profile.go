package Menu

import (
	"Forum/Fonctions/Connexion"
	"database/sql"
	"html/template"
	"net/http"
)

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

		data := struct {
			Username   string
			IsLoggedIn bool
		}{
			Username:   username,
			IsLoggedIn: true,
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
