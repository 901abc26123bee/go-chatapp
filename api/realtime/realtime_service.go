package realtime_service

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Store connected clients
var clients = make(map[*Client]bool)
var broadcast = make(chan Message)
var mutex = sync.Mutex{}

// Client defines a websocket connection client
type Client struct {
	Conn   *websocket.Conn
	UserID string
}

// Message Define a message object
type Message struct {
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_name"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %v\n", err)
		return
	}
	defer ws.Close()

	// Register new client
	mutex.Lock()
	client := &Client{Conn: ws}
	clients[client] = true
	mutex.Unlock()

	// Listen for new messages from the client
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			mutex.Lock()
			client := &Client{Conn: ws}
			delete(clients, client)
			mutex.Unlock()
			break
		}
		// Send the new message to the broadcast channel
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it to every connected client
		mutex.Lock()
		for client := range clients {
			err := client.Conn.WriteJSON(msg)
			if err != nil {
				fmt.Printf("Error writing message: %v\n", err)
				client.Conn.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
