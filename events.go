package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Handle save-promo event
func handleSavePromo(conn *websocket.Conn, promo string) {
	// Insert the new promo code into the database
	if err := InsertPromoCode(promo); err != nil {
		log.Println("Error inserting promo code:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error saving promo code"))
		return
	}

	// Send a success response
	successResponse := map[string]string{"status": "success", "message": "Promo code saved successfully"}
	responseBytes, err := json.Marshal(successResponse)
	if err != nil {
		log.Println("Error marshaling success response:", err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, responseBytes)

	// Create a broadcast message with the promo code
	broadcastMessage := map[string]string{"promo": promo}
	broadcastResponseBytes, err := json.Marshal(broadcastMessage)
	if err != nil {
		log.Println("Error marshaling broadcast message:", err)
		return
	}

	// Broadcast the promo code to all clients
	broadcast(broadcastResponseBytes)
}
