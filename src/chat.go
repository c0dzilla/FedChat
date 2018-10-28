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

var mode = flag.String("mode", "default", "Server mode: Central/Regular")
var address = flag.String("addr", "default", "Central Server Address")
var clients = make(map[*websocket.Conn]bool)
var servers = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var connect = make(chan Server)
var upgrader = websocket.Upgrader{}

func listenForServers() {
	if *address == "default" {
		log.Printf("No argument provided: Not connecting to central server...\n Enjoy as a standalone chat!")
		return
	}
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

func handleCentralConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("error: %v", err)
	}
	servers[ws] = true
	server := Server{Address: r.RemoteAddr}
	connect <- server
}

// don't send to new server
func emitServers() {
	for {
		newServer := <-connect

		for server := range servers {
			err := server.WriteJSON(newServer)

			if err != nil {
				log.Printf("error: %v", err)
			}
		}
	}
}

func sendServers() {

}

// Serve clients: central server
func Serve() {
	http.HandleFunc("/ws", handleConnections)
	go sendServers()
}

func main() {
	flag.Parse()

	if *mode == "central" {
		http.HandleFunc("/ws/central", handleCentralConnections)
		go emitServers()
	}
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
