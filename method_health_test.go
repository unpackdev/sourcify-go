package sourcify

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetHealth checks the server health by mocking the HTTP request and response.
func TestGetHealth(t *testing.T) {
	// Create a test server with a handler that always returns HTTP status OK (200)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create a new Sourcify client with the test server's URL
	client := NewClient(WithBaseURL(server.URL))

	// Call GetHealth
	isHealthy, err := GetHealth(client)
	assert.True(t, isHealthy, "Expected server to be healthy, got unhealthy")
	assert.Nil(t, err, "Unexpected error: %v", err)
}
