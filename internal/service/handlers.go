package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	profileToken := r.Header.Get("Profile-Token")
	if profileToken == "" {
		writeJSONResponse(w, http.StatusUnauthorized, "Missing tokens")
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

	w.Header().Set("Status", "200")
	responseJSON := map[string]interface{}{
		"status":  200,
		"message": "Message sent successfully",
	}
	json.NewEncoder(w).Encode(responseJSON)
}
