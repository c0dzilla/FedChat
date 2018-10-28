package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// Message 'go-enforced' comment -_-
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// Server holds information about connected servers
type Server struct {
	Address string `json:"address"`
}

var address = flag.String("addr", "localhost:8080", "Central Server Address")
var clients = make(map[*websocket.Conn]bool)
var servers = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var connect = make(chan Server)
var upgrader = websocket.Upgrader{}

func listenForServers() {
	centralServerURL := url.URL{Scheme: "ws", Host: *address, Path: "/ws"}
	ws, _, err := websocket.DefaultDialer.Dial(centralServerURL.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	for {
		var server Server
		err := ws.ReadJSON(&server)

		if err != nil {
			log.Printf("error: %v", err)
		}
		connect <- server
	}
}

func handleServers() {
	for {
		server := <-connect
		newServerURL := url.URL{Scheme: "ws", Host: server.Address, Path: "/ws/servers"}
		ws, _, err := websocket.DefaultDialer.Dial(newServerURL.String(), nil)

		if err != nil {
			log.Printf("error: %v", err)
		}
		defer ws.Close()
		servers[ws] = true
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleServerConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	servers[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(servers, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}

		for server := range servers {
			err := server.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				server.Close()
				delete(servers, server)
			}
		}
	}
}

func main() {
	flag.Parse()
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/ws/servers", handleServerConnections)

	go listenForServers()
	go handleServers()
	go handleMessages()

	log.Println("Listening at port 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
