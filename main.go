// package main

// import (
// 	"encoding/base64"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/IBM/sarama"
// )

// // JWTClaims represents the structure of the JWT payload
// type JWTClaims struct {
// 	FirstName string `json:"first_name"`
// 	UID       string `json:"uid"`
// }

// // decodeJWT decodes a JWT token and returns the payload as a JWTClaims struct
// func decodeJWT(token string) (*JWTClaims, error) {
// 	parts := strings.Split(token, ".")
// 	if len(parts) < 2 {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to decode token payload: %v", err)
// 	}

// 	var claims JWTClaims
// 	err = json.Unmarshal(payload, &claims)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to unmarshal token payload: %v", err)
// 	}

// 	return &claims, nil
// }

// // sendToKafka publishes the message to the Kafka topic
// func sendToKafka(from string, to string, message string) error {
// 	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to create Kafka producer: %v", err)
// 	}
// 	defer producer.Close()

// 	msg := &sarama.ProducerMessage{
// 		Topic: "chat-topic-1",
// 		Value: sarama.StringEncoder(fmt.Sprintf("From:%s, To:%s, Message:%s", from, to, message)),
// 	}

// 	_, _, err = producer.SendMessage(msg)
// 	if err != nil {
// 		return fmt.Errorf("failed to send message to Kafka: %v", err)
// 	}

// 	return nil
// }

// // SendMessageHandler is the HTTP handler for sending messages
// func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
// 	// authToken := r.Header.Get("Authorization")
// 	profileToken := r.Header.Get("Profile-Token")
// 	// profileToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoic2hyZWVyYWoiLCJ1aWQiOiI2NjYxYTdjNWZlYmE2YzRkOTA1NDhlZjQifQ.wziWXJ0soKg23jS1hwFHNgc6qTwg1RKPlcwMjiZLgPU"

// 	// if authToken == "" || profileToken == "" {
// 	if profileToken == "" {
// 		http.Error(w, "Missing tokens", http.StatusUnauthorized)
// 		return
// 	}

// 	// For simplicity, assume authToken is valid. Implement your auth logic here.

// 	claims, err := decodeJWT(profileToken)
// 	if err != nil {
// 		http.Error(w, "Invalid profile token", http.StatusUnauthorized)
// 		return
// 	}

// 	to := r.URL.Query().Get("to")
// 	if to == "" {
// 		http.Error(w, "Missing 'to' parameter", http.StatusBadRequest)
// 		return
// 	}

// 	message := "Hi, how are you?"
// 	err = sendToKafka(claims.FirstName, to, message)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Message sent successfully"))
// }

// func main() {
// 	http.HandleFunc("/send-message", SendMessageHandler)

// 	fmt.Println("Server starting at :8000")
// 	if err := http.ListenAndServe(":8000", nil); err != nil {
// 		log.Fatalf("Server error: %v", err)
// 	}
// }

//////shopify sarama

// package main

// import (
// 	"encoding/base64"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/shopify/sarama"
// )

// // JWTClaims represents the structure of the JWT payload
// type JWTClaims struct {
// 	FirstName string `json:"first_name"`
// 	UID       string `json:"uid"`
// }

// // decodeJWT decodes a JWT token and returns the payload as a JWTClaims struct
// func decodeJWT(token string) (*JWTClaims, error) {
// 	parts := strings.Split(token, ".")
// 	if len(parts) < 2 {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to decode token payload: %v", err)
// 	}

// 	var claims JWTClaims
// 	err = json.Unmarshal(payload, &claims)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to unmarshal token payload: %v", err)
// 	}

// 	return &claims, nil
// }

// // sendToKafka publishes the message to the Kafka topic
// func sendToKafka(from string, to string, message string) error {
// 	config := sarama.NewConfig()
// 	config.Producer.Return.Successes = true
// 	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
// 	if err != nil {
// 		return fmt.Errorf("failed to create Kafka producer: %v", err)
// 	}
// 	defer producer.Close()

// 	msg := &sarama.ProducerMessage{
// 		Topic: "chat-topic-2",
// 		Value: sarama.StringEncoder(fmt.Sprintf("From:%s, To:%s, Message:%s", from, to, message)),
// 	}

// 	_, _, err = producer.SendMessage(msg)
// 	if err != nil {
// 		return fmt.Errorf("failed to send message to Kafka: %v", err)
// 	}

// 	return nil
// }

// // writeJSONResponse writes a JSON response with a given status code and message
// func writeJSONResponse(w http.ResponseWriter, statusCode int, message string) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	responseData := map[string]interface{}{
// 		"status":  statusCode,
// 		"message": message,
// 	}
// 	responseJSON, _ := json.Marshal(responseData)
// 	w.Write(responseJSON)
// }

// // SendMessageHandler is the HTTP handler for sending messages
// func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
// 	profileToken := r.Header.Get("Profile-Token")
// 	if profileToken == "" {
// 		writeJSONResponse(w, http.StatusUnauthorized, "Missing tokens")
// 		return
// 	}

// 	claims, err := decodeJWT(profileToken)
// 	if err != nil {
// 		writeJSONResponse(w, http.StatusUnauthorized, "Invalid profile token")
// 		return
// 	}

// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		writeJSONResponse(w, http.StatusBadRequest, "Failed to read request body")
// 		return
// 	}

// 	var requestData map[string]string
// 	err = json.Unmarshal(body, &requestData)
// 	if err != nil {
// 		writeJSONResponse(w, http.StatusBadRequest, "Invalid JSON in request body")
// 		return
// 	}

// 	to, toExists := requestData["to"]
// 	message, msgExists := requestData["msg"]

// 	if !toExists || to == "" {
// 		writeJSONResponse(w, http.StatusBadRequest, "Missing 'to' parameter")
// 		return
// 	}
// 	if !msgExists || message == "" {
// 		writeJSONResponse(w, http.StatusBadRequest, "Missing 'msg' parameter")
// 		return
// 	}

// 	err = sendToKafka(claims.FirstName, to, message)
// 	if err != nil {
// 		writeJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to send message: %v", err))
// 		return
// 	}

// 	w.Header().Set("Status", "200")
// 	responseData := map[string]interface{}{
// 		"status":  200,
// 		"message": "Message sent successfully",
// 	}
// 	responseJSON, _ := json.Marshal(responseData)
// 	w.Write(responseJSON)
// }

// func main() {
// 	http.HandleFunc("/send-message", SendMessageHandler)

// 	fmt.Println("Server starting at :8000")
// 	if err := http.ListenAndServe(":8000", nil); err != nil {
// 		log.Fatalf("Server error: %v", err)
// 	}
// }

// ////ibm sarama
// package main

// import (
// 	"encoding/base64"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/IBM/sarama"
// )

// // JWTClaims represents the structure of the JWT payload
// type JWTClaims struct {
// 	FirstName string `json:"first_name"`
// 	UID       string `json:"uid"`
// }

// // decodeJWT decodes a JWT token and returns the payload as a JWTClaims struct
// func decodeJWT(token string) (*JWTClaims, error) {
// 	parts := strings.Split(token, ".")
// 	if len(parts) < 2 {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to decode token payload: %v", err)
// 	}

// 	var claims JWTClaims
// 	err = json.Unmarshal(payload, &claims)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to unmarshal token payload: %v", err)
// 	}

// 	return &claims, nil
// }

// // sendToKafka publishes the message to the Kafka topic
// func sendToKafka(from string, to string, message string) error {
// 	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
// 	if err != nil {
// 		return fmt.Errorf("failed to create Kafka producer: %v", err)
// 	}
// 	defer producer.Close()

// 	msg := &sarama.ProducerMessage{
// 		Topic: "chat-topic-2",
// 		Value: sarama.StringEncoder(fmt.Sprintf("From:%s, To:%s, Message:%s", from, to, message)),
// 	}

// 	_, _, err = producer.SendMessage(msg)
// 	if err != nil {
// 		return fmt.Errorf("failed to send message to Kafka: %v", err)
// 	}

// 	return nil
// }

// // writeJSONResponse writes a JSON response with a given status code and message
// func writeJSONResponse(w http.ResponseWriter, statusCode int, message string) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
// 	responseData := map[string]interface{}{
// 		"status":  statusCode,
// 		"message": message,
// 	}
// 	responseJSON, _ := json.Marshal(responseData)
// 	w.Write(responseJSON)
// }

// // SendMessageHandler is the HTTP handler for sending messages
// func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
// 	profileToken := r.Header.Get("Profile-Token")
// 	if profileToken == "" {
// 		writeJSONResponse(w, http.StatusUnauthorized, "Missing tokens")
// 		return
// 	}

// 	claims, err := decodeJWT(profileToken)
// 	if err != nil {
// 		writeJSONResponse(w, http.StatusUnauthorized, "Invalid profile token")
// 		return
// 	}

// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		writeJSONResponse(w, http.StatusBadRequest, "Failed to read request body")
// 		return
// 	}

// 	var requestData map[string]string
// 	err = json.Unmarshal(body, &requestData)
// 	if err != nil {
// 		writeJSONResponse(w, http.StatusBadRequest, "Invalid JSON in request body")
// 		return
// 	}

// 	to, toExists := requestData["to"]
// 	message, msgExists := requestData["msg"]

// 	if !toExists || to == "" {
// 		writeJSONResponse(w, http.StatusBadRequest, "Missing 'to' parameter")
// 		return
// 	}
// 	if !msgExists || message == "" {
// 		writeJSONResponse(w, http.StatusBadRequest, "Missing 'msg' parameter")
// 		return
// 	}

// 	err = sendToKafka(claims.FirstName, to, message)
// 	if err != nil {
// 		writeJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to send message: %v", err))
// 		return
// 	}

// 	w.Header().Set("Status", "200")
// 	responseData := map[string]interface{}{
// 		"status":  200,
// 		"message": "Message sent successfully",
// 	}
// 	responseJSON, _ := json.Marshal(responseData)
// 	w.Write(responseJSON)
// }

// func main() {
// 	http.HandleFunc("/send-message", SendMessageHandler)

// 	fmt.Println("Server starting at :8000")
// 	if err := http.ListenAndServe(":8000", nil); err != nil {
// 		log.Fatalf("Server error: %v", err)
// 	}
// }

///////seperation of functions

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
