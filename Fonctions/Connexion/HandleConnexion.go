package Connexion

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

var (
	Store       = sessions.NewCookieStore([]byte("my-secret-key"))
	SessionName = "my-session"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func compareHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func HandleConnexion(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		// Gérer l'erreur
	}
	defer db.Close()
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		var storedHash string
		err := db.QueryRow("SELECT password FROM utilisateurs WHERE username = ?", username).Scan(&storedHash)
		if errors.Is(err, sql.ErrNoRows) || !compareHash(password, storedHash) {
			fmt.Fprintln(w, "Pseudo ou mot de passe invalide")
			return
		}
		session, err := Store.Get(r, SessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["username"] = username
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	tmpl, err := template.ParseFiles("templates/html/Connexion/Connexion.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
