package main

import (
	"Forum/Fonctions/Connexion"
	"Forum/Fonctions/Menu"
	"Forum/Fonctions/Websocket"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	certFile := "KeyHTTPS/certificat.crt"
	keyFile := "KeyHTTPS/privatekey.key"

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		Websocket.UploadImage(w, r, db)
	})

	http.HandleFunc("/getImages", func(w http.ResponseWriter, r *http.Request) {
		Websocket.DisplayImage(w, r, db)
	})

	http.HandleFunc("/", Menu.HandlePlay)
	http.HandleFunc("/profile", Menu.HandleProfile)
	http.HandleFunc("/postes", Menu.HandlePostes)
	http.HandleFunc("/services", Menu.HandleService)
	http.HandleFunc("/about", Menu.HandleAbout)
	http.HandleFunc("/logout", Connexion.HandleLogout)
	http.HandleFunc("/connexion", Connexion.HandleConnexion)
	http.HandleFunc("/inscription", Connexion.HandleInscription)
	http.HandleFunc("/password", Connexion.HandleForgetPassword)
	http.HandleFunc("/websocket", Websocket.HandleWebsocket)
	http.HandleFunc("/Echec", Websocket.HandleWebsocketEchec)

	http.HandleFunc("/Foot", Websocket.HandleWebsocketFoot)
	http.HandleFunc("/LikesDislikesFoot", Websocket.LikesDislikesHandlerFoot)
	http.HandleFunc("/likesFoot", Websocket.LikeHandlerFoot)
	http.HandleFunc("/dislikesFoot", Websocket.DislikeHandlerFoot)
	http.HandleFunc("/wsFoot", Websocket.WebSocketHandlerFoot)

	http.HandleFunc("/Bacterie", Websocket.HandleWebsocketBacterie)
	http.HandleFunc("/LikesDislikesBacterie", Websocket.LikesDislikesHandlerBacterie)
	http.HandleFunc("/likesBacterie", Websocket.LikeHandlerBacterie)
	http.HandleFunc("/dislikesBacterie", Websocket.DislikeHandlerBacterie)
	http.HandleFunc("/wsBacterie", Websocket.WebSocketHandlerBacterie)

	http.HandleFunc("/Cyber", Websocket.HandleWebsocketCyber)
	http.HandleFunc("/LikesDislikesCyber", Websocket.LikesDislikesHandlerCyber)
	http.HandleFunc("/likesCyber", Websocket.LikeHandlerCyber)
	http.HandleFunc("/dislikesCyber", Websocket.DislikeHandlerCyber)
	http.HandleFunc("/wsCyber", Websocket.WebSocketHandlerCyber)

	http.HandleFunc("/Emploie", Websocket.HandleWebsocketEmploie)
	http.HandleFunc("/LikesDislikesEmploie", Websocket.LikesDislikesHandlerEmploie)
	http.HandleFunc("/likesEmploie", Websocket.LikeHandlerEmploie)
	http.HandleFunc("/dislikesEmploie", Websocket.DislikeHandlerEmploie)
	http.HandleFunc("/wsEmploie", Websocket.WebSocketHandlerEmploie)

	http.HandleFunc("/Histoire", Websocket.HandleWebsocketHistoire)
	http.HandleFunc("/LikesDislikesHistoire", Websocket.LikesDislikesHandlerHistoire)
	http.HandleFunc("/likesHistoire", Websocket.LikeHandlerHistoire)
	http.HandleFunc("/dislikesHistoire", Websocket.DislikeHandlerHistoire)
	http.HandleFunc("/wsHistoire", Websocket.WebSocketHandlerHistoire)

	http.HandleFunc("/Lit", Websocket.HandleWebsocketLit)
	http.HandleFunc("/LikesDislikesLit", Websocket.LikesDislikesHandlerLit)
	http.HandleFunc("/likesLit", Websocket.LikeHandlerLit)
	http.HandleFunc("/dislikesLit", Websocket.DislikeHandlerLit)
	http.HandleFunc("/wsLit", Websocket.WebSocketHandlerLit)

	http.HandleFunc("/Livres", Websocket.HandleWebsocketLivres)
	http.HandleFunc("/LikesDislikesLivres", Websocket.LikesDislikesHandlerLivres)
	http.HandleFunc("/likesLivres", Websocket.LikeHandlerLivres)
	http.HandleFunc("/dislikesLivres", Websocket.DislikeHandlerLivres)
	http.HandleFunc("/wsLivres", Websocket.WebSocketHandlerLivres)

	http.HandleFunc("/Meuble", Websocket.HandleWebsocketMeuble)
	http.HandleFunc("/LikesDislikesMeuble", Websocket.LikesDislikesHandlerMeuble)
	http.HandleFunc("/likesMeuble", Websocket.LikeHandlerMeuble)
	http.HandleFunc("/dislikesMeuble", Websocket.DislikeHandlerMeuble)
	http.HandleFunc("/wsMeuble", Websocket.WebSocketHandlerMeuble)

	http.HandleFunc("/Mirroir", Websocket.HandleWebsocketMirroir)
	http.HandleFunc("/LikesDislikesMirroir", Websocket.LikesDislikesHandlerMirroir)
	http.HandleFunc("/likesMirroir", Websocket.LikeHandlerMirroir)
	http.HandleFunc("/dislikesMirroir", Websocket.DislikeHandlerMirroir)
	http.HandleFunc("/wsMirroir", Websocket.WebSocketHandlerMirroir)

	http.HandleFunc("/MMA", Websocket.HandleWebsocketMMA)
	http.HandleFunc("/LikesDislikesMMA", Websocket.LikesDislikesHandlerMMA)
	http.HandleFunc("/likesMMA", Websocket.LikeHandlerMMA)
	http.HandleFunc("/dislikesMMA", Websocket.DislikeHandlerMMA)
	http.HandleFunc("/wsMMA", Websocket.WebSocketHandlerMMA)

	http.HandleFunc("/Musique", Websocket.HandleWebsocketMusique)
	http.HandleFunc("/LikesDislikesMusique", Websocket.LikesDislikesHandlerMusique)
	http.HandleFunc("/likesMusique", Websocket.LikeHandlerMusique)
	http.HandleFunc("/dislikesMusique", Websocket.DislikeHandlerMusique)
	http.HandleFunc("/wsMusique", Websocket.WebSocketHandlerMusique)

	http.HandleFunc("/Navigateurs", Websocket.HandleWebsocketNavigateurs)
	http.HandleFunc("/LikesDislikesNavigateurs", Websocket.LikesDislikesHandlerNavigateurs)
	http.HandleFunc("/likesNavigateurs", Websocket.LikeHandlerNavigateurs)
	http.HandleFunc("/dislikesNavigateurs", Websocket.DislikeHandlerNavigateurs)
	http.HandleFunc("/wsNavigateurs", Websocket.WebSocketHandlerNavigateurs)

	http.HandleFunc("/Nourriture", Websocket.HandleWebsocketNourriture)
	http.HandleFunc("/LikesDislikesNourriture", Websocket.LikesDislikesHandlerNourriture)
	http.HandleFunc("/likesNourriture", Websocket.LikeHandlerNourriture)
	http.HandleFunc("/dislikesNourriture", Websocket.DislikeHandlerNourriture)
	http.HandleFunc("/wsNourriture", Websocket.WebSocketHandlerNourriture)

	http.HandleFunc("/Nucleaire", Websocket.HandleWebsocketNucleaire)
	http.HandleFunc("/LikesDislikesNucleaire", Websocket.LikesDislikesHandlerNucleaire)
	http.HandleFunc("/likesNucleaire", Websocket.LikeHandlerNucleaire)
	http.HandleFunc("/dislikesNucleaire", Websocket.DislikeHandlerNucleaire)
	http.HandleFunc("/wsNucleaire", Websocket.WebSocketHandlerNucleaire)

	http.HandleFunc("/PC", Websocket.HandleWebsocketPC)
	http.HandleFunc("/LikesDislikesPC", Websocket.LikesDislikesHandlerPC)
	http.HandleFunc("/likesPC", Websocket.LikeHandlerPC)
	http.HandleFunc("/dislikesPC", Websocket.DislikeHandlerPC)
	http.HandleFunc("/wsPC", Websocket.WebSocketHandlerPC)

	http.HandleFunc("/Romans", Websocket.HandleWebsocketRomans)
	http.HandleFunc("/LikesDislikesRomans", Websocket.LikesDislikesHandlerRomans)
	http.HandleFunc("/likesRomans", Websocket.LikeHandlerRomans)
	http.HandleFunc("/dislikesRomans", Websocket.DislikeHandlerRomans)
	http.HandleFunc("/wsRomans", Websocket.WebSocketHandlerRomans)

	http.HandleFunc("/Rugby", Websocket.HandleWebsocketRugby)
	http.HandleFunc("/LikesDislikesRugby", Websocket.LikesDislikesHandlerRugby)
	http.HandleFunc("/likesRugby", Websocket.LikeHandlerRugby)
	http.HandleFunc("/dislikesRugby", Websocket.DislikeHandlerRugby)
	http.HandleFunc("/wsRugby", Websocket.WebSocketHandlerRugby)

	http.HandleFunc("/Stage", Websocket.HandleWebsocketStage)
	http.HandleFunc("/LikesDislikesStage", Websocket.LikesDislikesHandlerStage)
	http.HandleFunc("/likesStage", Websocket.LikeHandlerStage)
	http.HandleFunc("/dislikesStage", Websocket.DislikeHandlerStage)
	http.HandleFunc("/wsStage", Websocket.WebSocketHandlerStage)

	http.HandleFunc("/Virus", Websocket.HandleWebsocketVirus)
	http.HandleFunc("/LikesDislikesVirus", Websocket.LikesDislikesHandlerVirus)
	http.HandleFunc("/likesVirus", Websocket.LikeHandlerVirus)
	http.HandleFunc("/dislikesVirus", Websocket.DislikeHandlerVirus)
	http.HandleFunc("/wsVirus", Websocket.WebSocketHandlerVirus)

	http.HandleFunc("/Youtubeurs", Websocket.HandleWebsocketYoutubeurs)
	http.HandleFunc("/LikesDislikesYoutubeurs", Websocket.LikesDislikesHandlerYoutubeurs)
	http.HandleFunc("/likesYoutubeurs", Websocket.LikeHandlerYoutubeurs)
	http.HandleFunc("/dislikesYoutubeurs", Websocket.DislikeHandlerYoutubeurs)
	http.HandleFunc("/wsYoutubeurs", Websocket.WebSocketHandlerYoutubeurs)

	http.HandleFunc("/Programmation", Websocket.HandleWebsocketProgrammation)
	http.HandleFunc("/LikesDislikesProgrammation", Websocket.LikesDislikesHandlerProgrammation)
	http.HandleFunc("/likesProgrammation", Websocket.LikeHandlerProgrammantion)
	http.HandleFunc("/dislikesProgrammation", Websocket.DislikeHandlerProgrammation)
	http.HandleFunc("/wsProgrammation", Websocket.WebSocketHandlerProgrammation)

	http.HandleFunc("/likes", Websocket.LikeHandlerEchec)
	http.HandleFunc("/dislikes", Websocket.DislikeHandlerEchec)
	http.HandleFunc("/LikesDislikesEchec", Websocket.LikesDislikesHandlerEchec)

	fmt.Println("Server started on port :443 https://localhost:443")
	err = http.ListenAndServeTLS(":443", certFile, keyFile, nil)
	if err != nil {
		log.Fatal("Erreur de démarrage du serveur HTTPS : ", err)
	}
}