package sourcify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestGetContractMetadata(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/server/repository/contracts/full_match/1/0x0000000000000000000000001234567890aBcdEF/metadata.json" {
			// Simulate a successful response with sample metadata
			metadata := &Metadata{
				Compiler: Compiler{Version: "0.1.0"},
				Language: "Solidity",
				Output: Output{
					Abi: []Abi{
						{
							Inputs: []any{},
							Type:   "function",
							Name:   "myFunction",
							Outputs: []OutputDetail{
								{
									Type: "uint256",
									Name: "result",
								},
							},
						},
					},
				},
				Version: 1,
			}

			err := json.NewEncoder(w).Encode(metadata)
			if err != nil {
				t.Errorf("failed to encode mock metadata: %v", err)
			}
		} else {
			http.NotFound(w, r)
		}
	}))
	defer mockServer.Close()

	// Create a client for the mock server
	client := NewClient(WithBaseURL(mockServer.URL + "/server"))

	// Define test data
	chainID := 1
	contractAddress := common.HexToAddress("0x1234567890abcdef")
	matchType := MethodMatchTypeFull

	// Call the function to get contract metadata
	metadata, err := GetContractMetadata(client, chainID, contractAddress, matchType)

	// Verify the results
	assert.NoError(t, err, "GetContractMetadata returned an error")

	expectedMetadata := &Metadata{
		Compiler: Compiler{Version: "0.1.0"},
		Language: "Solidity",
		Output: Output{
			Abi: []Abi{
				{
					Inputs: []any{},
					Type:   "function",
					Name:   "myFunction",
					Outputs: []OutputDetail{
						{
							Type: "uint256",
							Name: "result",
						},
					},
				},
			},
		},
		Version: 1,
	}

	assert.Equal(t, expectedMetadata, metadata, "GetContractMetadata returned unexpected metadata")
}
