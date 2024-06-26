package helpers

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDecoder mocks the Decoder interface for testing
type MockDecoder struct {
	mock.Mock
}

// DecodeString is a mocked method to decode a string
func (m *MockDecoder) DecodeString(s string) ([]byte, error) {
	args := m.Called(s)
	return args.Get(0).([]byte), args.Error(1)
}

func TestDecodeJWT(t *testing.T) {
	// Positive test case
	t.Run("ValidToken", func(t *testing.T) {
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiSm9obiIsInVpZCI6InVpZCIsImV4cCI6MTUxNjIzOTAyMn0"

		// Mock Decoder
		mockDecoder := new(MockDecoder)
		mockDecoder.On("DecodeString", mock.Anything).Return([]byte(`{"first_name":"John","uid":"uid"}`), nil)

		// Call the function being tested
		claims, err := DecodeJWT(token)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "John", claims.FirstName)
		assert.Equal(t, "uid", claims.UID)
	})

	// Negative test case - invalid token structure
	t.Run("InvalidTokenStructure", func(t *testing.T) {
		token := "invalidToken"

		// Mock Decoder
		mockDecoder := new(MockDecoder)
		mockDecoder.On("DecodeString", mock.Anything).Return(nil, errors.New("unable to decode token payload: illegal base64 data at input byte 7"))

		// Call the function being tested
		claims, err := DecodeJWT(token)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	// Negative test case - unable to decode token payload
	t.Run("InvalidBase64Payload", func(t *testing.T) {
		token := "header.invalidBase64.signature"

		// Call the function being tested
		claims, err := DecodeJWT(token)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "unable to decode token payload")
	})

	// Negative test case - unable to unmarshal token payload
	t.Run("InvalidJSONPayload", func(t *testing.T) {
		token := "header.payload.signature"
		invalidJSON := base64.RawURLEncoding.EncodeToString([]byte("invalid json"))

		// Replace the payload part with the invalid JSON payload
		tokenParts := strings.Split(token, ".")
		tokenParts[1] = invalidJSON
		token = strings.Join(tokenParts, ".")

		// Call the function being tested
		claims, err := DecodeJWT(token)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "unable to unmarshal token payload")
	})
}
