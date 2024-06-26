package main

import (
	"log"
	"net/http"

	"chatInteractionService/cmd/api/routes/internal/handlers"
	"chatInteractionService/cmd/api/routes/internal/helpers"
	"chatInteractionService/cmd/api/routes/internal/middleware"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func main() {

	fetchSecrets := func(client *secretsmanager.Client) (string, string, string, string, string, error) {
		return helpers.FetchSecrets(client)
	}

	http.Handle("/send-message", middleware.AuthMiddleware(http.HandlerFunc(handlers.SendMessageHandler), fetchSecrets))

	log.Println("Server starting at :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
