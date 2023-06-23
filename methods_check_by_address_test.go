package sourcify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestCheckContractByAddresses(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == MethodCheckByAddresses.URI {
			// Simulate a successful response with sample contract addresses
			contractAddresses := []*CheckContractAddress{
				{
					Address:  common.HexToAddress("0x1234567890abcdef"),
					Status:   "verified",
					ChainIDs: []string{"1", "2"},
				},
			}

			err := json.NewEncoder(w).Encode(contractAddresses)
			if err != nil {
				t.Errorf("failed to encode mock contract addresses: %v", err)
			}
		} else {
			http.NotFound(w, r)
		}
	}))
	defer mockServer.Close()

	// Create a client for the mock server
	client := NewClient(WithBaseURL(mockServer.URL))

	// Define test data
	addresses := []string{"0x054B2223509D430269a31De4AE2f335890be5C8F"}
	chainIds := []int{56}

	// Call the function to check contract addresses
	contractAddresses, err := CheckContractByAddresses(client, addresses, chainIds, MethodMatchTypeFull)

	// Verify the results
	assert.NoError(t, err, "CheckContractByAddresses returned an error")

	expectedContractAddresses := []*CheckContractAddress{
		{
			Address:  common.HexToAddress("0x1234567890abcdef"),
			Status:   "verified",
			ChainIDs: []string{"1", "2"},
		},
	}

	assert.Equal(t, expectedContractAddresses, contractAddresses, "CheckContractByAddresses returned unexpected contract addresses")
}

func TestCheckContractByAddresses_Error(t *testing.T) {
	// Create a mock HTTP server that returns an error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// Create a client for the mock server
	client := NewClient(WithBaseURL(mockServer.URL))

	// Define test data
	addresses := []string{"0x054B2223509D430269a31De4AE2f335890be5C8F"}
	chainIds := []int{56}

	// Call the function to check contract addresses
	contractAddresses, err := CheckContractByAddresses(client, addresses, chainIds, MethodMatchTypeFull)

	// Verify the results
	assert.Error(t, err, "CheckContractByAddresses should return an error")
	assert.Nil(t, contractAddresses, "CheckContractByAddresses should return nil contract addresses")
}

func TestCheckContractByAddresses_InvalidMatchType(t *testing.T) {
	// Create a client for testing
	client := NewClient(WithBaseURL("https://example.com"))

	// Define test data
	addresses := []string{"0x054B2223509D430269a31De4AE2f335890be5C8F"}
	chainIds := []int{56}

	// Call the function with an invalid match type
	contractAddresses, err := CheckContractByAddresses(client, addresses, chainIds, "invalid")

	// Verify the results
	assert.Error(t, err, "CheckContractByAddresses should return an error for invalid match type")
	assert.Nil(t, contractAddresses, "CheckContractByAddresses should return nil contract addresses for invalid match type")
}
