package sourcify

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestGetContractByChainIdAndAddress(t *testing.T) {
	// Test it via the Ethereum USDT contract
	chainID := 1
	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	// Load the expected response from testdata file
	localResponse, err := LoadContract(chainID, contractAddress)
	require.NoError(t, err)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fmt.Sprintf("/v2/contract/%d/%s", chainID, contractAddress.Hex()) {
			// Simulate a successful response with sample source code
			err := json.NewEncoder(w).Encode(localResponse)
			if err != nil {
				t.Errorf("failed to encode mock source codes: %v", err)
			}
		} else {
			http.NotFound(w, r)
		}
	}))
	defer mockServer.Close()

	// Create a client for the mock server
	client := NewClient(WithBaseURL(mockServer.URL))

	// Call the function to get contract source code
	contractResponse, err := GetContractByChainIdAndAddress(client, chainID, contractAddress, []string{}, []string{})

	// Verify the results
	assert.NoError(t, err, "GetContractByChainIdAndAddress returned an error")

	assert.Equal(t, localResponse, contractResponse, "GetContractByChainIdAndAddress returned unexpected source codes")
}

func TestUpstreamGetContractByChainIdAndAddress(t *testing.T) {
	// Skip this test in automated testing environments as it requires internet connection
	if testing.Short() {
		t.Skip("Skipping test in short mode as it requires internet connection")
	}

	// Create a custom HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create a new Sourcify client with custom options
	client := NewClient(
		WithHTTPClient(httpClient),
		WithBaseURL("https://sourcify.dev/server"),
		WithRetryOptions(
			WithMaxRetries(3),
			WithDelay(2*time.Second),
		),
	)

	// Test it via the Ethereum USDT contract
	chainID := 1
	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	// Load the expected response from testdata file
	localResponse, err := LoadContract(chainID, contractAddress)
	if err != nil {
		// If the test data file doesn't exist yet, we can create it
		t.Logf("Test data file not found: %v", err)
		t.Log("To generate test data, run this test with the GENERATE_TEST_DATA=true environment variable")

		if os.Getenv("GENERATE_TEST_DATA") == "true" {
			// Call the function to get contract source code
			contractResponse, err := GetContractByChainIdAndAddress(client, chainID, contractAddress, []string{}, []string{})
			require.NoError(t, err, "GetContractByChainIdAndAddress returned an error")

			// Save the contract response for future tests
			err = SaveContract(contractResponse)
			require.NoError(t, err, "Failed to save contract response")
			t.Logf("Generated test data for %s on chain %d", contractAddress.Hex(), chainID)
			return
		}

		t.SkipNow()
	}

	// Call the function to get contract source code
	contractResponse, err := GetContractByChainIdAndAddress(client, chainID, contractAddress, []string{}, []string{})

	// Verify the results
	require.NoError(t, err, "GetContractByChainIdAndAddress returned an error")

	// Compare the responses
	assert.Equal(t, localResponse, contractResponse, "GetContractByChainIdAndAddress returned unexpected contract response")
}
