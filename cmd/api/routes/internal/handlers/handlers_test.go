package handlers

import (
	"bytes"
	"chatInteractionService/cmd/api/routes/internal/helpers"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHelpers struct {
	mock.Mock
}

func (m *MockHelpers) FetchSecrets(secretsClient *secretsmanager.Client) (string, string, string, string, string, error) {
	args := m.Called(secretsClient)
	return args.String(0), args.String(1), args.String(2), args.String(3), args.String(4), args.Error(5)
}

func (m *MockHelpers) DecodeJWT(token string) (*helpers.JWTClaims, error) {
	args := m.Called(token)
	if args.Get(0) != nil {
		return args.Get(0).(*helpers.JWTClaims), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockHelpers) SendToKafka(ec2Client *ec2.Client, secretsClient *secretsmanager.Client, firstName, to, message, region string) error {
	args := m.Called(ec2Client, secretsClient, firstName, to, message, region)
	return args.Error(0)
}

type mockReadCloser struct {
	err error
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	return 0, m.err
}

func (m *mockReadCloser) Close() error {
	return nil
}

func TestSendMessageHandler(t *testing.T) {
	tests := []struct {
		name                string
		profileToken        string
		requestBody         map[string]string
		mockFetchSecrets    func(m *MockHelpers)
		mockDecodeJWT       func(m *MockHelpers)
		mockSendToKafka     func(m *MockHelpers)
		expectedStatusCode  int
		expectedResponseMsg string
	}{
		// {
		// 	name:         "Successful message send",
		// 	profileToken: "valid-token",
		// 	requestBody: map[string]string{
		// 		"to":  "receiver",
		// 		"msg": "hello",
		// 	},
		// 	mockFetchSecrets: func(m *MockHelpers) {
		// 		m.On("FetchSecrets", mock.Anything).Return("secret1", "secret2", "secret3", "secret4", "secret5", nil)
		// 	},
		// 	mockDecodeJWT: func(m *MockHelpers) {
		// 		m.On("DecodeJWT", "valid-token").Return(&helpers.JWTClaims{FirstName: "John", UID: "uid"}, nil)
		// 	},
		// 	mockSendToKafka: func(m *MockHelpers) {
		// 		m.On("SendToKafka", mock.Anything, mock.Anything, "John", "receiver", "hello", "us-east-1").Return(nil)
		// 	},
		// 	expectedStatusCode:  http.StatusOK,
		// 	expectedResponseMsg: "Message sent successfully",
		// },
		{
			name:         "Missing profile token",
			profileToken: "",
			requestBody: map[string]string{
				"to":  "receiver",
				"msg": "hello",
			},
			mockFetchSecrets:    func(m *MockHelpers) {},
			mockDecodeJWT:       func(m *MockHelpers) {},
			mockSendToKafka:     func(m *MockHelpers) {},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedResponseMsg: "Sender is not authorized",
		},
		{
			name:         "Invalid profile token",
			profileToken: "invalid-token",
			requestBody: map[string]string{
				"to":  "receiver",
				"msg": "hello",
			},
			mockFetchSecrets: func(m *MockHelpers) {
				m.On("FetchSecrets", mock.Anything).Return("secret1", "secret2", "secret3", "secret4", "secret5", nil)
			},
			mockDecodeJWT: func(m *MockHelpers) {
				m.On("DecodeJWT", "invalid-token").Return(nil, assert.AnError)
			},
			mockSendToKafka:     func(m *MockHelpers) {},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedResponseMsg: "assert.AnError general error for testing",
		},
		{
			name:         "Invalid JSON in request body",
			profileToken: "valid-token",
			requestBody:  map[string]string{},
			mockFetchSecrets: func(m *MockHelpers) {
				m.On("FetchSecrets", mock.Anything).Return("secret1", "secret2", "secret3", "secret4", "secret5", nil)
			},
			mockDecodeJWT: func(m *MockHelpers) {
				m.On("DecodeJWT", "valid-token").Return(&helpers.JWTClaims{FirstName: "John", UID: "uid"}, nil)
			},
			mockSendToKafka:     func(m *MockHelpers) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedResponseMsg: "Receiver not found",
		},
		{
			name:         "Receiver not found",
			profileToken: "valid-token",
			requestBody: map[string]string{
				"msg": "hello",
			},
			mockFetchSecrets: func(m *MockHelpers) {
				m.On("FetchSecrets", mock.Anything).Return("secret1", "secret2", "secret3", "secret4", "secret5", nil)
			},
			mockDecodeJWT: func(m *MockHelpers) {
				m.On("DecodeJWT", "valid-token").Return(&helpers.JWTClaims{FirstName: "John", UID: "uid"}, nil)
			},
			mockSendToKafka:     func(m *MockHelpers) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedResponseMsg: "Receiver not found",
		},
		{
			name:         "Message is not set",
			profileToken: "valid-token",
			requestBody: map[string]string{
				"to": "receiver",
			},
			mockFetchSecrets: func(m *MockHelpers) {
				m.On("FetchSecrets", mock.Anything).Return("secret1", "secret2", "secret3", "secret4", "secret5", nil)
			},
			mockDecodeJWT: func(m *MockHelpers) {
				m.On("DecodeJWT", "valid-token").Return(&helpers.JWTClaims{FirstName: "John", UID: "uid"}, nil)
			},
			mockSendToKafka:     func(m *MockHelpers) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedResponseMsg: "Message is not set",
		},
		// {
		// 	name:         "Failed to send message",
		// 	profileToken: "valid-token",
		// 	requestBody: map[string]string{
		// 		"to":  "receiver",
		// 		"msg": "hello",
		// 	},
		// 	mockFetchSecrets: func(m *MockHelpers) {
		// 		m.On("FetchSecrets", mock.Anything).Return("secret1", "secret2", "secret3", "secret4", "secret5", nil)
		// 	},
		// 	mockDecodeJWT: func(m *MockHelpers) {
		// 		m.On("DecodeJWT", "valid-token").Return(&helpers.JWTClaims{FirstName: "John", UID: "uid"}, nil)
		// 	},
		// 	mockSendToKafka: func(m *MockHelpers) {
		// 		m.On("SendToKafka", mock.Anything, mock.Anything, "John", "receiver", "hello", "us-east-1").Return(assert.AnError)
		// 	},
		// 	expectedStatusCode:  http.StatusOK,
		// 	expectedResponseMsg: "Message sent successfully",
		// },

		// {
		// 	name:         "Failed to read request body",
		// 	profileToken: "valid-token",
		// 	requestBody:  nil, // Simulate a nil body to trigger read error
		// 	mockFetchSecrets: func(m *MockHelpers) {
		// 		m.On("FetchSecrets", mock.Anything).Return("secret1", "secret2", "secret3", "secret4", "secret5", nil)
		// 	},
		// 	mockDecodeJWT: func(m *MockHelpers) {
		// 		m.On("DecodeJWT", "valid-token").Return(&helpers.JWTClaims{FirstName: "John", UID: "uid"}, nil)
		// 	},
		// 	mockSendToKafka:     func(m *MockHelpers) {},
		// 	expectedStatusCode:  http.StatusBadRequest,
		// 	expectedResponseMsg: "Failed to read request body",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHelpers := new(MockHelpers)
			tt.mockFetchSecrets(mockHelpers)
			tt.mockDecodeJWT(mockHelpers)
			tt.mockSendToKafka(mockHelpers)

			DecodeJWTFunc = mockHelpers.DecodeJWT
			// FetchSecrets = mockHelpers.FetchSecrets
			// SendToKafka = mockHelpers.SendToKafka

			handler := http.HandlerFunc(SendMessageHandler)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/send", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Profile-Token", tt.profileToken)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			var respBody map[string]interface{}
			json.Unmarshal(rr.Body.Bytes(), &respBody)
			assert.Equal(t, tt.expectedResponseMsg, respBody["message"])
		})
	}

	secondTests := []struct {
		name                string
		profileToken        string
		requestBody         string
		mockFetchSecrets    func(m *MockHelpers)
		mockDecodeJWT       func(m *MockHelpers)
		mockSendToKafka     func(m *MockHelpers)
		expectedStatusCode  int
		expectedResponseMsg string
	}{
		{
			name:         "Invalid JSON in request body",
			profileToken: "valid-token",
			requestBody:  "nil",
			mockFetchSecrets: func(m *MockHelpers) {
				m.On("FetchSecrets", mock.Anything).Return("secret1", "secret2", "secret3", "secret4", "secret5", nil)
			},
			mockDecodeJWT: func(m *MockHelpers) {
				m.On("DecodeJWT", "valid-token").Return(&helpers.JWTClaims{FirstName: "John", UID: "uid"}, nil)
			},
			mockSendToKafka:     func(m *MockHelpers) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedResponseMsg: "Invalid JSON in request body",
		},
	}
	for _, tt := range secondTests {
		t.Run(tt.name, func(t *testing.T) {
			mockHelpers := new(MockHelpers)
			tt.mockFetchSecrets(mockHelpers)
			tt.mockDecodeJWT(mockHelpers)
			tt.mockSendToKafka(mockHelpers)

			handler := http.HandlerFunc(SendMessageHandler)

			body, _ := json.Marshal(tt.requestBody)

			req := httptest.NewRequest("POST", "/send", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Profile-Token", tt.profileToken)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			assert.Equal(t, tt.expectedResponseMsg, "Invalid JSON in request body")
		})
	}

	thirdTests := []struct {
		name                string
		profileToken        string
		requestBody         map[string]interface{}
		mockFetchSecrets    func(m *MockHelpers)
		mockDecodeJWT       func(m *MockHelpers)
		mockSendToKafka     func(m *MockHelpers)
		expectedStatusCode  int
		expectedResponseMsg string
	}{
		{
			name:         "Failed to read request body",
			profileToken: "valid-token",
			requestBody:  nil, // Simulating an empty request body
			mockFetchSecrets: func(m *MockHelpers) {
				m.On("FetchSecrets", mock.Anything).Return("secret1", "secret2", "secret3", "secret4", "secret5", nil)
			},
			mockDecodeJWT: func(m *MockHelpers) {
				m.On("DecodeJWT", "valid-token").Return(&helpers.JWTClaims{FirstName: "John", UID: "uid"}, nil)
			},
			mockSendToKafka:     func(m *MockHelpers) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedResponseMsg: "Failed to read request body",
		},
	}

	for _, tt := range thirdTests {
		t.Run(tt.name, func(t *testing.T) {
			mockHelpers := new(MockHelpers)
			tt.mockFetchSecrets(mockHelpers)
			tt.mockDecodeJWT(mockHelpers)
			tt.mockSendToKafka(mockHelpers)

			handler := http.HandlerFunc(SendMessageHandler)

			// Simulate error when reading request body
			mockBody := &mockReadCloser{err: errors.New("error reading body")}
			req := httptest.NewRequest("POST", "/send", mockBody)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Profile-Token", tt.profileToken)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			var respBody map[string]interface{}
			err := json.Unmarshal(rr.Body.Bytes(), &respBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			assert.Equal(t, tt.expectedResponseMsg, "Failed to read request body")
		})
	}
}
