module Forum

replace Fonction/Fonction => ./Fonction

replace Fonction/Connexion => ./Fonctions/Connexion

replace Fonction/Level => ./Level

require (
	github.com/gorilla/sessions v1.2.2
	github.com/gorilla/websocket v1.5.1
	github.com/mattn/go-sqlite3 v1.14.18
	golang.org/x/crypto v0.15.0
)

require (
	github.com/gorilla/securecookie v1.1.2 // indirect
	golang.org/x/net v0.17.0 // indirect
)

go 1.21
