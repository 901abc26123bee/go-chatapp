package realtime

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
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
	Chat     string `json:"message"`
}

// Chat Define a chat object
type Chat struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Msg       string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for simplicity, adjust as needed
		},
	}
)

// ServeWS upgrade http connection to a WebSocket and register client
func ServeWS(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial http request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error upgrading connection: %v\n", err)
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
			log.Printf("Error reading message: %v\n", err)
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

func Broadcaster() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it to every connected client
		mutex.Lock()
		for client := range clients {
			err := client.Conn.WriteJSON(msg)
			if err != nil {
				log.Errorf("Error writing message: %v\n", err)
				client.Conn.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
