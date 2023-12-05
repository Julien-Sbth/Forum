package Websocket

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func HandleWebsocketLivres(w http.ResponseWriter, r *http.Request) {
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
		oldMessages, err = getOldLivresMessagesFromDB()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imageList, err := getAllImageURLsFromDBLivres()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var oldLikes []LikesDislikes
		oldLikes, err = getOldLikesDislikesFromDBLivres()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var imageDatas []ImageData
		for _, imageURL := range imageList {
			image := ImageData{
				URL: imageURL,
			}
			imageDatas = append(imageDatas, image)
		}

		tmpl, err := template.ParseFiles("templates/html/Discussion/livres.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			URL           []ImageData
			Username      string
			LikesDislikes []LikesDislikes
			Messages      []Message
			Content       string
			IsLoggedIn    bool
		}{
			URL:           imageDatas,
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

	tmpl, err := template.ParseFiles("templates/html/Discussion/livres.html")
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

func getAllImageURLsFromDBLivres() ([]string, error) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT data FROM LivresImages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var imageList []string
	for rows.Next() {
		var imageURL string
		if err := rows.Scan(&imageURL); err != nil {
			return nil, err
		}
		imageList = append(imageList, imageURL)
	}

	return imageList, nil
}

func saveImageToDBLivres(imageData string) error {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO LivresImages (data) VALUES (?)", imageData)
	if err != nil {
		return err
	}

	return nil
}

func WebSocketHandlerLivres(w http.ResponseWriter, r *http.Request) {
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

		if msg.Image != "" {
			err = saveImageToDBLivres(msg.Image)
			if err != nil {
				log.Println("Error saving image to database:", err)
			}
		}

		err = saveMessageToDBLivres(msg)
		if err != nil {
			log.Println("Error saving message to database:", err)
		}

		messages = append(messages, msg)

		sendNewMessageToAllClientsExceptSender(msg, conn)
	}
}

func saveMessageToDBLivres(msg Message) error {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO Livres (username, content, likes, dislikes) VALUES (?, ?, ?, ?)", msg.Username, msg.Content, msg.Likes, msg.Dislikes)
	if err != nil {
		return err
	}

	return nil
}

func getOldLikesDislikesFromDBLivres() ([]LikesDislikes, error) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, likes, dislikes FROM Livres ORDER BY id")
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

func LikesDislikesHandlerLivres(w http.ResponseWriter, r *http.Request) {
	likesDislikes, err := getOldLikesDislikesFromDBLivres()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	likesDislikesJSON, err := json.Marshal(likesDislikes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(likesDislikesJSON)
}

func getOldLivresMessagesFromDB() ([]Message, error) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, username, content FROM Livres ORDER BY id")
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

func incrementLikesLivres(messageID int) error {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE Livres SET likes = likes + 1 WHERE id = ?", messageID)
	if err != nil {
		return err
	}

	return nil
}

func incrementDislikesLivres(messageID int) error {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE Livres SET dislikes = dislikes + 1 WHERE id = ?", messageID)
	if err != nil {
		return err
	}

	return nil
}

func LikeHandlerLivres(w http.ResponseWriter, r *http.Request) {
	messageID := r.FormValue("id")
	messageIDInt, err := strconv.Atoi(messageID)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = incrementLikesLivres(messageIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/Livres", http.StatusSeeOther)
}

func DislikeHandlerLivres(w http.ResponseWriter, r *http.Request) {
	messageID := r.FormValue("id")
	messageIDInt, err := strconv.Atoi(messageID)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = incrementDislikesLivres(messageIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/Livres", http.StatusSeeOther)
}

func UploadLivres(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	imageBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	_, err = db.Exec("INSERT INTO LivresImages (Data) VALUES (?)", imageBase64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imageURL := fmt.Sprintf("/image?id=%d", handler.Filename)

	response := struct {
		Success  bool   `json:"success"`
		ImageURL string `json:"imageURL"`
	}{
		Success:  true,
		ImageURL: imageURL,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ImageHandlerLivres(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT data FROM LivresImages")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var imageURLs []string

	for rows.Next() {
		var imageURL string
		err := rows.Scan(&imageURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imageURLs = append(imageURLs, imageURL)
	}

	imageJSON, err := json.Marshal(imageURLs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(imageJSON)
}
