package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"chatInteractionService/cmd/api/routes/internal/helpers"

	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization token", http.StatusUnauthorized)
			return
		}

		authTokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_, _, _, region, userPoolID, err := helpers.FetchSecrets()
		if err != nil {
			log.Println("Couldn't retrieve the secrets")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		ctx := context.Background()
		// log.Println(userPoolID)

		_, err = cognitoJwtAuthenticator.ValidateToken(ctx, region, userPoolID, authTokenString)
		if err != nil {
			http.Error(w, "Token validation error", http.StatusUnauthorized)
			return
		}

		log.Println("Authorization token is valid")
		next.ServeHTTP(w, r)
	})
}
