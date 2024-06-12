package service

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
)

// writeJSONResponse writes a JSON response with a given status code and message
func writeJSONResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	responseData := map[string]interface{}{
		"status":  statusCode,
		"message": message,
	}
	responseJSON, _ := json.Marshal(responseData)
	w.Write(responseJSON)
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]string
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		authTokenString := authHeader[len("Bearer "):]
		region := os.Getenv("REGION")
		userPoolId := os.Getenv("COGNITO_USER_POOL_ID")
		ctx := context.Background()

		_, err := cognitoJwtAuthenticator.ValidateToken(ctx, region, userPoolId, authTokenString)
		if err == nil {
			log.Println("Authorization token is valid")

			profileToken := r.Header.Get("Profile-Token")
			if profileToken != "" {
				claims, err := decodeJWT(profileToken)
				if err == nil {
					log.Println("Profile token is valid")

					body, err := io.ReadAll(r.Body)
					if err == nil {
						err = json.Unmarshal(body, &requestData)
						if err == nil {
							// toExists and msgExists are boolean
							to, toExists := requestData["to"]
							message, msgExists := requestData["msg"]

							if toExists && to != "" {
								if msgExists && message != "" {
									err = sendToKafka(claims.FirstName, to, message)
									if err == nil {
										writeJSONResponse(w, http.StatusOK, "Message sent successfully")
										log.Println("--------------------------------------------------")
										return
									} else {
										writeJSONResponse(w, http.StatusInternalServerError, err.Error())
										log.Println("Failed to send message")
										return
									}
								} else {
									writeJSONResponse(w, http.StatusBadRequest, "Message is not set")
									log.Println("Missing 'msg' parameter")
									return
								}
							} else {
								writeJSONResponse(w, http.StatusBadRequest, "Receiver not found")
								log.Println("Missing 'to' parameter")
								return
							}
						} else {
							writeJSONResponse(w, http.StatusBadRequest, err.Error())
							log.Println("Invalid JSON in request body")
							return
						}
					} else {
						writeJSONResponse(w, http.StatusBadRequest, err.Error())
						log.Println("Failed to read request body")
						return
					}
				} else {
					writeJSONResponse(w, http.StatusUnauthorized, err.Error())
					log.Println("Invalid Profile token")
					return
				}
			} else {
				writeJSONResponse(w, http.StatusUnauthorized, "Sender is not authorized")
				log.Println("Missing Profile token")
				return
			}
		} else {
			writeJSONResponse(w, http.StatusUnauthorized, err.Error())
			log.Println("Token validation error")
			return
		}
	} else {
		writeJSONResponse(w, http.StatusUnauthorized, "Missing Authorization token")
		return
	}
}
