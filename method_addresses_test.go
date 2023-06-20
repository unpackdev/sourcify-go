package sourcify

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestGetAvailableContractAddresses(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a sample JSON response
		response := `{"full": ["0x1234567890123456789012345678901234567890"], "partial": []}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	// Create a client with the mock server URL
	client := NewClient(WithBaseURL(server.URL))

	// Call the function being tested
	addresses, err := GetAvailableContractAddresses(client, 123)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, addresses)
	expectedAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
	assert.Equal(t, []common.Address{expectedAddress}, addresses.Full)
	assert.Empty(t, addresses.Partial)
}
