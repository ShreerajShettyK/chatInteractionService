package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// JWTClaims represents the structure of the profile's JWT payload
type JWTClaims struct {
	FirstName string `json:"first_name"`
	UID       string `json:"uid"`
}

// decodeJWT decodes a JWT token and returns the payload as a JWTClaims struct
func decodeJWT(token string) (*JWTClaims, error) {
	var claims JWTClaims
	parts := strings.Split(token, ".") //header,payload,signature
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("unable to decode token payload: %v", err)
	}

	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal token payload: %v", err)
	}

	return &claims, nil
}
