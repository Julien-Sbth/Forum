package Connexion

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

func HandleForgetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/html/Connexion/ForgetPassword.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		newPassword := r.FormValue("new_password")
		token := r.FormValue("token")

		db, err := sql.Open("sqlite3", "database.sqlite")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var storedToken string
		err = db.QueryRow("SELECT reset_token FROM utilisateurs WHERE username = ?", username).Scan(&storedToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Vérification du token
		if token != storedToken {
			http.Error(w, "Token invalide", http.StatusUnauthorized)
			return
		}

		// Générer un nouveau token pour la prochaine réinitialisation
		nextResetToken, err := generateToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Hasher le nouveau mot de passe avec bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Mise à jour du mot de passe dans la base de données
		_, err = db.Exec("UPDATE utilisateurs SET password = ?, reset_token = ? WHERE username = ?", hashedPassword, nextResetToken, username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Passer le token généré au modèle
		data := struct {
			Token string
		}{
			Token: nextResetToken,
		}

		tmpl, err := template.ParseFiles("templates/html/Connexion/ForgetPassword.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Afficher la page avec le token généré
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
