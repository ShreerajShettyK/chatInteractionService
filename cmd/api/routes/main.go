// package main

// import (
// 	"log"
// 	"net/http"

// 	"chatInteractionService/cmd/api/routes/internal/handlers"
// 	"chatInteractionService/cmd/api/routes/internal/helpers"
// 	"chatInteractionService/cmd/api/routes/internal/middleware"

// 	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
// )

// func main() {

// 	fetchSecrets := func(client *secretsmanager.Client) (string, string, string, string, string, error) {
// 		return helpers.FetchSecrets(client)
// 	}

// 	http.Handle("/send-message", middleware.AuthMiddleware(http.HandlerFunc(handlers.SendMessageHandler), fetchSecrets))

// 	log.Println("Server starting at :8000")
// 	if err := http.ListenAndServe(":8000", nil); err != nil {
// 		log.Fatalf("Server error: %v", err)
// 	}
// }

package main

import (
	"context"
	"log"
	"net/http"

	"chatInteractionService/cmd/api/routes/internal/handlers"
	"chatInteractionService/cmd/api/routes/internal/helpers"
	"chatInteractionService/cmd/api/routes/internal/middleware"

	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func main() {
	// cfg, err := config.LoadDefaultConfig(context.Background())
	// if err != nil {
	// 	log.Fatalf("Error loading AWS SDK config: %v", err)
	// }

	// secretsClient := secretsmanager.NewFromConfig(cfg)

	fetchSecrets := func(client *secretsmanager.Client) (string, string, string, string, string, error) {
		return helpers.FetchSecrets(client)
	}

	validateToken := func(ctx context.Context, region, userPoolID, token string) (*cognitoJwtAuthenticator.AWSCognitoClaims, error) {
		return cognitoJwtAuthenticator.ValidateToken(ctx, region, userPoolID, token)
	}

	http.Handle("/send-message", middleware.AuthMiddleware(http.HandlerFunc(handlers.SendMessageHandler), fetchSecrets, validateToken))

	log.Println("Server starting at :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
