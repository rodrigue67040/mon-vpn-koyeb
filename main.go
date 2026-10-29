package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Accepter la connexion WebSocket entrante (venant de DarkTunnel)
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erreur upgrade:", err)
		return
	}
	defer wsConn.Close()

	// 2. Se connecter au VPS OVH en TCP
	targetAddr := os.Getenv("TARGET_ADDR")
	if targetAddr == "" {
		targetAddr = "15.235.60.159:443" // Votre IP VPS
	}
	tcpConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Println("Erreur connexion TCP vers le VPS:", err)
		return
	}
	defer tcpConn.Close()

	// 3. Faire le pont entre les deux connexions
	errc := make(chan error, 2)
	go func() {
		_, err := io.Copy(tcpConn, wsConn.UnderlyingConn())
		errc <- err
	}()
	go func() {
		_, err := io.Copy(wsConn.UnderlyingConn(), tcpConn)
		errc <- err
	}()
	<-errc
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.HandleFunc("/", proxyHandler)
	log.Printf("Proxy en écoute sur le port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}