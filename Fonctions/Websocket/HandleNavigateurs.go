package Websocket

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func HandleWebsocketNavigateurs(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var oldMessages []Message
	if username, ok := session.Values["username"].(string); ok {
		// Récupérer les anciens messages depuis la base de données
		oldMessages, err = getOldNavigateursMessagesFromDB()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var oldLikes []LikesDislikes
		oldLikes, err = getOldLikesDislikesFromDBNavigateurs()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/html/Discussion/navigateurs.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Likes         int
			Dislikes      int
			Username      string
			Messages      []Message
			Content       string
			IsLoggedIn    bool
			LikesDislikes []LikesDislikes
		}{
			Username:      username,
			LikesDislikes: oldLikes,
			Messages:      oldMessages,
			Content:       "",
			IsLoggedIn:    true,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	tmpl, err := template.ParseFiles("templates/html/Menu/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func WebSocketHandlerNavigateurs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	clients[conn] = true
	defer delete(clients, conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("WebSocket JSON unmarshal error:", err)
			continue
		}

		msg.SocketID = conn.LocalAddr().String()

		err = saveMessageToDBNavigateurs(msg)
		if err != nil {
			log.Println("Error saving message to database:", err)
		}

		messages = append(messages, msg)

		sendNewMessageToAllClientsExceptSender(msg, conn) // Send message to all clients except the sender
	}
}

func saveMessageToDBNavigateurs(msg Message) error {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	// Exécute la requête SQL pour insérer le message dans la table des messages
	_, err = db.Exec("INSERT INTO navigateurs (username, content, likes, dislikes) VALUES (?, ?, ?, ?)", msg.Username, msg.Content, msg.Likes, msg.Dislikes)
	if err != nil {
		return err
	}

	return nil
}

func getOldLikesDislikesFromDBNavigateurs() ([]LikesDislikes, error) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, likes, dislikes FROM navigateurs ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var oldLikesDislikes []LikesDislikes
	for rows.Next() {
		var id, likes, dislikes int
		err := rows.Scan(&id, &likes, &dislikes)
		if err != nil {
			return nil, err
		}

		ld := LikesDislikes{
			ID:       id,
			Likes:    likes,
			Dislikes: dislikes,
		}
		oldLikesDislikes = append(oldLikesDislikes, ld)
	}

	return oldLikesDislikes, nil
}

func LikesDislikesHandlerNavigateurs(w http.ResponseWriter, r *http.Request) {
	// Appel de la deuxième fonction pour récupérer les likes et dislikes à partir de la DB
	likesDislikes, err := getOldLikesDislikesFromDBNavigateurs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convertir les likes/dislikes en JSON et les renvoyer comme réponse
	likesDislikesJSON, err := json.Marshal(likesDislikes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(likesDislikesJSON)
}

// Fonction pour récupérer les anciens messages depuis la base de données
func getOldNavigateursMessagesFromDB() ([]Message, error) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, username, content FROM navigateurs ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var oldMessages []Message
	for rows.Next() {
		var id int
		var username, content string
		err := rows.Scan(&id, &username, &content)
		if err != nil {
			return nil, err
		}

		msg := Message{
			ID:       id,
			Username: username,
			Content:  content,
		}
		oldMessages = append(oldMessages, msg)
	}

	return oldMessages, nil
}

func incrementLikesNavigateurs(messageID int) error {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE navigateurs SET likes = likes + 1 WHERE id = ?", messageID)
	if err != nil {
		return err
	}

	return nil
}

func incrementDislikesNavigateurs(messageID int) error {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE navigateurs SET dislikes = dislikes + 1 WHERE id = ?", messageID)
	if err != nil {
		return err
	}

	return nil
}

func LikeHandlerNavigateurs(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID du message depuis les données du formulaire
	messageID := r.FormValue("id")
	messageIDInt, err := strconv.Atoi(messageID)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = incrementLikesNavigateurs(messageIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/bacterie", http.StatusSeeOther)
}

func DislikeHandlerNavigateurs(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID du message depuis les données du formulaire
	messageID := r.FormValue("id")
	messageIDInt, err := strconv.Atoi(messageID)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = incrementDislikesNavigateurs(messageIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/navigateurs", http.StatusSeeOther)
}