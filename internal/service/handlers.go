package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

// SendMessageHandler is the HTTP handler for sending messages
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeJSONResponse(w, http.StatusUnauthorized, "Missing Authorization token")
		return
	}

	tokenString := authHeader[len("Bearer "):]
	region := "us-east-1"
	userPoolId := "us-east-1_icXeg2eiv"
	ctx := context.Background()

	_, err := cognitoJwtAuthenticator.ValidateToken(ctx, region, userPoolId, tokenString)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, fmt.Sprintf("Token validation error: %s", err))
		return
	}

	log.Println("Authorization token is valid")
	profileToken := r.Header.Get("Profile-Token")
	if profileToken == "" {
		writeJSONResponse(w, http.StatusUnauthorized, "Missing Profile token")
		return
	}

	claims, err := decodeJWT(profileToken)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, "Invalid profile token")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	log.Println("Profile token is valid")

	var requestData map[string]string
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, "Invalid JSON in request body")
		return
	}

	to, toExists := requestData["to"]
	message, msgExists := requestData["msg"]

	if !toExists || to == "" {
		writeJSONResponse(w, http.StatusBadRequest, "Missing 'to' parameter")
		return
	}
	if !msgExists || message == "" {
		writeJSONResponse(w, http.StatusBadRequest, "Missing 'msg' parameter")
		return
	}

	err = sendToKafka(claims.FirstName, to, message)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to send message: %v", err))
		return
	}
	log.Println("Sent message to kafka")

	w.Header().Set("Status", "200")
	responseJSON := map[string]interface{}{
		"status":  200,
		"message": "Message sent successfully",
	}
	json.NewEncoder(w).Encode(responseJSON)
}
