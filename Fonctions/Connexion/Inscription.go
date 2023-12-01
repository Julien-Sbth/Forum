package Connexion

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	_ "fmt"
	_ "github.com/mattn/go-sqlite3"
	_ "golang.org/x/crypto/bcrypt"
	"html/template"
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

type RegisterData struct {
	ErrorMessage string
}

func HandleInscription(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		if !ValidatePassword(password) {
			http.Error(w, "Le mot de passe ne respecte pas les conditions requises", http.StatusBadRequest)
			return
		}

		exists, err := checkUsernameExists(db, username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if exists {
			data := struct {
				ExistingEmail    bool
				ExistingUsername bool
				ErrorMessage     string
				ErrorEmail       string
			}{
				ExistingEmail:    true,
				ExistingUsername: true,
				ErrorMessage:     "Le nom d'utilisateur est déjà pris",
				ErrorEmail:       "L'email est déjà pris",
			}

			tmpl, err := template.ParseFiles("templates/html/Connexion/inscription.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return // Arrête le traitement car le nom d'utilisateur existe déjà
		}

		emailExists, err := checkEmailExists(db, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if emailExists {
			data := RegisterData{
				ErrorMessage: "L'email est déjà utilisé",
			}
			tmpl, err := template.ParseFiles("templates/html/Erreur/error.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		hashedPassword, err := hashPassword(password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := generateToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Insertion des données dans la base de données avec le token
		_, err = db.Exec("INSERT INTO utilisateurs (username, password, email, reset_token) VALUES (?, ?, ?, ?)", username, hashedPassword, email, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Si l'insertion s'est faite avec succès, redirige vers la page de connexion
		http.Redirect(w, r, "/connexion", http.StatusFound)
		return
	}

	// Rendu de la page d'inscription si la méthode HTTP n'est pas POST
	tmpl, err := template.ParseFiles("templates/html/Connexion/inscription.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
