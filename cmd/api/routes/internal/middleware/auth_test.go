package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ShreerajShettyK/cognitoJwtAuthenticator"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCognitoJwtAuthenticator struct {
	mock.Mock
}

func (m *MockCognitoJwtAuthenticator) ValidateToken(ctx context.Context, region, userPoolID, token string) (*cognitoJwtAuthenticator.AWSCognitoClaims, error) {
	args := m.Called(ctx, region, userPoolID, token)
	return args.Get(0).(*cognitoJwtAuthenticator.AWSCognitoClaims), args.Error(1)
}

type MockHelpers struct {
	mock.Mock
}

func (m *MockHelpers) FetchSecrets(client *secretsmanager.Client) (string, string, string, string, string, error) {
	args := m.Called(client)
	return args.String(0), args.String(1), args.String(2), args.String(3), args.String(4), args.Error(5)
}

func TestAuthMiddleware(t *testing.T) {
	// mockCognito := new(MockCognitoJwtAuthenticator)
	mockHelpers := new(MockHelpers)

	mockFetchSecrets := func(client *secretsmanager.Client) (string, string, string, string, string, error) {
		return mockHelpers.FetchSecrets(client)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	authMiddleware := AuthMiddleware(handler, mockFetchSecrets)

	t.Run("MissingAuthorizationHeader", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		authMiddleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Missing Authorization token\n", rr.Body.String())
	})

	// t.Run("InvalidToken", func(t *testing.T) {
	// 	req, _ := http.NewRequest("GET", "/", nil)
	// 	req.Header.Set("Authorization", "Bearer invalidtoken")
	// 	rr := httptest.NewRecorder()

	// 	mockHelpers.On("FetchSecrets", mock.Anything).Return("", "", "", "us-west-2", "userpoolID", nil).Once()
	// 	mockCognito.On("ValidateToken", mock.Anything, "us-west-2", "userpoolID", "invalidtoken").Return(nil, assert.AnError).Once()

	// 	authMiddleware.ServeHTTP(rr, req)

	// 	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	// 	assert.Equal(t, "Token validation error\n", rr.Body.String())
	// 	mockHelpers.AssertExpectations(t)
	// 	mockCognito.AssertExpectations(t)
	// })

	// t.Run("ValidToken", func(t *testing.T) {
	// 	req, _ := http.NewRequest("GET", "/", nil)
	// 	req.Header.Set("Authorization", "Bearer validtoken")
	// 	rr := httptest.NewRecorder()

	// 	mockHelpers.On("FetchSecrets", mock.Anything).Return("", "", "", "us-west-2", "userpoolID", nil).Once()
	// 	mockCognito.On("ValidateToken", mock.Anything, "us-west-2", "userpoolID", "validtoken").Return(&cognitoJwtAuthenticator.AWSCognitoClaims{}, nil).Once()

	// 	authMiddleware.ServeHTTP(rr, req)

	// 	assert.Equal(t, http.StatusOK, rr.Code)
	// 	assert.Equal(t, "Token validation error\n", rr.Body.String())
	// 	mockHelpers.AssertExpectations(t)
	// 	mockCognito.AssertExpectations(t)
	// })

	t.Run("FetchSecretsError", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer sometoken")
		rr := httptest.NewRecorder()

		mockHelpers.On("FetchSecrets", mock.Anything).Return("", "", "", "", "", assert.AnError).Once()

		authMiddleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "Internal Server Error\n", rr.Body.String())
		mockHelpers.AssertExpectations(t)
	})
}
