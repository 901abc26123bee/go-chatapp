package realtime_service

import (
	"fmt"
	"net/http"
)

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
	clients[ws] = true
	mutex.Unlock()

	// Listen for new messages from the client
	for {
			var msg Message
			err := ws.ReadJSON(&msg)
			if err != nil {
					fmt.Printf("Error reading message: %v\n", err)
					mutex.Lock()
					delete(clients, ws)
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
					err := client.WriteJSON(msg)
					if err != nil {
							fmt.Printf("Error writing message: %v\n", err)
							client.Close()
							delete(clients, client)
					}
			}
			mutex.Unlock()
	}
}