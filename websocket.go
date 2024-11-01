package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Define the structure for incoming messages
type Message struct {
	Event string `json:"event"`
	Data  struct {
		UserID int    `json:"user_id,omitempty"` // UserID is optional for some events
		Promo  string `json:"promo,omitempty"`   // Promo is optional for check-promo
	} `json:"data,omitempty"` // Data can be empty for some events
}

var upgrader = websocket.Upgrader{}

// Client represents a single WebSocket connection
type Client struct {
	conn *websocket.Conn
}

// Clients manager
var clients = struct {
	sync.Mutex
	connections map[*Client]struct{}
}{
	connections: make(map[*Client]struct{}),
}

// Handle WebSocket connections
func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	client := &Client{conn: conn}

	// Register new client
	clients.Lock()
	clients.connections[client] = struct{}{}
	clients.Unlock()

	log.Println("Client connected")

	defer func() {
		clients.Lock()
		delete(clients.connections, client)
		clients.Unlock()
		log.Println("Client disconnected")
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error unmarshaling message:", err)
			continue // Skip this message if it can't be parsed
		}

		log.Printf("Received event: %s\n", message.Event)

		switch message.Event {
		case "save-promo":
			handleSavePromo(conn, message.Data.Promo) // Call the function for this event
		default:
			log.Println("Unknown event type:", message.Event)
		}
	}
}

// Broadcast sends a message to all connected clients
func broadcast(message []byte) {
	clients.Lock()
	defer clients.Unlock()
	for client := range clients.connections {
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing message to client:", err)
			client.conn.Close()                 // Close the connection on error
			delete(clients.connections, client) // Remove client from list
		}
	}
}
