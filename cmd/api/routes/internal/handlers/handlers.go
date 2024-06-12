package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"chatInteractionService/cmd/api/routes/internal/helpers"
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
	_, _, _, _, _, err := helpers.FetchSecrets()
	if err != nil {
		log.Println("Couldn't retrieve the secrets")
		writeJSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	log.Println("Secret retrieved from AWS Secrets Manager")
	profileToken := r.Header.Get("Profile-Token")
	if profileToken != "" {
		claims, err := helpers.DecodeJWT(profileToken)
		if err == nil {
			log.Println("Profile token is valid")

			body, err := io.ReadAll(r.Body)
			if err == nil {
				err = json.Unmarshal(body, &requestData)
				if err == nil {
					to, toExists := requestData["to"]
					message, msgExists := requestData["msg"]

					if toExists && to != "" {
						if msgExists && message != "" {
							err = helpers.SendToKafka(claims.FirstName, to, message)
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
}
