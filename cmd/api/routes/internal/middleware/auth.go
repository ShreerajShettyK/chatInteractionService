// package middleware

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
// )

// var cfg aws.Config

// func init() {
// 	var err error
// 	cfg, err = config.LoadDefaultConfig(context.Background())
// 	if err != nil {
// 		log.Fatalf("Error loading AWS SDK config: %v", err)
// 	}
// }

// func AuthMiddleware(next http.Handler, fetchSecrets func(client *secretsmanager.Client) (string, string, string, string, string, error)) http.Handler {
// 	secretsClient := secretsmanager.NewFromConfig(cfg)
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			http.Error(w, "Missing Authorization token", http.StatusUnauthorized)
// 			return
// 		}

// 		authTokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 		_, _, _, region, userPoolID, err := fetchSecrets(secretsClient)
// 		if err != nil {
// 			log.Println("Couldn't retrieve the secrets")
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 			return
// 		}
// 		ctx := context.Background()

// 		_, err = cognitoJwtAuthenticator.ValidateToken(ctx, region, userPoolID, authTokenString)
// 		if err != nil {
// 			http.Error(w, "Token validation error", http.StatusUnauthorized)
// 			return
// 		}

// 		log.Println("Authorization token is valid")
// 		next.ServeHTTP(w, r)
// 	})
// }

package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var cfg aws.Config

func init() {
	var err error
	cfg, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Error loading AWS SDK config: %v", err)
	}
}

func AuthMiddleware(next http.Handler, fetchSecrets func(client *secretsmanager.Client) (string, string, string, string, string, error), validateToken func(ctx context.Context, region, userPoolID, token string) (*cognitoJwtAuthenticator.AWSCognitoClaims, error)) http.Handler {
	secretsClient := secretsmanager.NewFromConfig(cfg)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization token", http.StatusUnauthorized)
			return
		}

		authTokenString := strings.TrimPrefix(authHeader, "Bearer ")

		_, _, _, region, userPoolID, err := fetchSecrets(secretsClient)
		if err != nil {
			log.Println("Couldn't retrieve the secrets")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		ctx := context.Background()

		_, err = validateToken(ctx, region, userPoolID, authTokenString)
		if err != nil {
			http.Error(w, "Token validation error", http.StatusUnauthorized)
			return
		}

		log.Println("Authorization token is valid")
		next.ServeHTTP(w, r)
	})
}
