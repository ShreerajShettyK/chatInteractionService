package main

import (
	"fmt"
	"log"
	"net/http"
	"practicechat/internal/service"
)

func main() {
	http.HandleFunc("/send-message", service.SendMessageHandler)

	fmt.Println("Server starting at :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
