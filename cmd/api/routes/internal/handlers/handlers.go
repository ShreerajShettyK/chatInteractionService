package handlers

import (
	"chatInteractionService/cmd/api/routes/internal/helpers"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	DecodeJWTFunc = helpers.DecodeJWT
	FetchSecrets  = helpers.FetchSecrets
	SendToKafka   = helpers.SendToKafka
)

var cfg aws.Config

func init() {
	var err error
	cfg, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Error loading AWS SDK config: %v", err)
	}
}

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

// SendMessageHandler handles the sending of a message to Kafka
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]string
	secretsClient := secretsmanager.NewFromConfig(cfg)
	region := "us-east-1" // Set your AWS region here

	ec2Client := ec2.NewFromConfig(cfg)

	_, _, _, _, _, err := FetchSecrets(secretsClient)
	if err != nil {
		log.Println("Couldn't retrieve the secrets")
		writeJSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	log.Println("Secret retrieved from AWS Secrets Manager")

	profileToken := r.Header.Get("Profile-Token")
	if profileToken != "" {
		claims, err := DecodeJWTFunc(profileToken)
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
							err = SendToKafka(ec2Client, secretsClient, claims.FirstName, to, message, region)
							if err == nil {
								writeJSONResponse(w, http.StatusOK, "Message sent successfully")
								log.Println("--------------------------------------------------")
								return
							} else {
								writeJSONResponse(w, http.StatusInternalServerError, "Failed to send message")
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
