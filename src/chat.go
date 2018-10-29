package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
var serverBroadcast = make(chan Message)
var connect = make(chan Server)
var upgrader = websocket.Upgrader{}

// holds websocket of newly joined server(only used by central server)
var newlyJoinedServer *websocket.Conn

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
		serverBroadcast <- msg
	}
}

func emitToClients(msg Message) {
	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func emitToServers(msg Message) {
	for server := range servers {
		err := server.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			server.Close()
			delete(servers, server)
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		emitToClients(msg)
		emitToServers(msg)
	}
}

func handleServerMessages() {
	for {
		msg := <-serverBroadcast
		emitToClients(msg)
	}
}

func handleCentralConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("error: %v", err)
	}
	servers[ws] = true
	newlyJoinedServer = ws
	remoteAddr := strings.Split(r.RemoteAddr, ":")[0] + ":8080"
	log.Println(remoteAddr + " joined as server...")
	server := Server{Address: remoteAddr}
	connect <- server
}

func emitServers() {
	for {
		newServer := <-connect
		log.Println("Emitting " + newServer.Address + " to connected servers...")

		for server := range servers {
			if server == newlyJoinedServer {
				continue
			}
			err := server.WriteJSON(newServer)

			if err != nil {
				log.Printf("error: %v", err)
			}
		}
		log.Printf("Done: Emitted")
	}
}

// Serve clients: central server
func Serve() {
	http.HandleFunc("/ws", handleConnections)
	go emitServers()
}

func main() {
	flag.Parse()

	if *mode == "central" {
		log.Println("Starting central server...")
		http.HandleFunc("/ws", handleCentralConnections)
		go emitServers()
		startServer(8090)
		return
	}
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/ws/servers", handleServerConnections)
	go listenForServers()
	go handleServers()
	go handleMessages()
	go handleServerMessages()
	startServer(8080)
}

func startServer(port int) {
	log.Printf("Listening at port %d", port)
	portString := ":" + strconv.Itoa(port)
	err := http.ListenAndServe(portString, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
