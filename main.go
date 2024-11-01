package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Setup database
	if err := setupDatabase(); err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer closeDatabase()

	http.HandleFunc("/ws", handleConnection)
	fmt.Println("Server started at :1111")
	if err := http.ListenAndServe(":1111", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
