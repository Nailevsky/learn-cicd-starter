package auth

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

var ErrNoAuthHeaderIncluded = errors.New("no authorization header included")

// GetAPIKey extracts the API key from the Authorization header
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	return "broken-key", nil
}

// Test case for a valid API key extraction
func TestGetAPIKey_ValidKey(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "ApiKey test-key")

	apiKey, err := GetAPIKey(headers)

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if apiKey != "test-key" {
		t.Errorf("Expected 'test-key', but got '%s'", apiKey)
	}
}

// Test case for missing Authorization header
func TestGetAPIKey_MissingHeader(t *testing.T) {
	headers := http.Header{} // No Authorization header

	apiKey, err := GetAPIKey(headers)

	if err == nil {
		t.Fatalf("Expected an error, but got none")
	}
	if !errors.Is(err, ErrNoAuthHeaderIncluded) {
		t.Errorf("Expected error '%v', but got '%v'", ErrNoAuthHeaderIncluded, err)
	}
	if apiKey != "" {
		t.Errorf("Expected empty API key, but got '%s'", apiKey)
	}
}

// Test case for a malformed Authorization header (missing "ApiKey" prefix)
func TestGetAPIKey_MalformedHeader(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-key") // Wrong prefix

	apiKey, err := GetAPIKey(headers)

	if err == nil {
		t.Fatalf("Expected an error, but got none")
	}
	expectedErr := "malformed authorization header"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error containing '%s', but got '%v'", expectedErr, err)
	}
	if apiKey != "" {
		t.Errorf("Expected empty API key, but got '%s'", apiKey)
	}
}
