package sourcify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestGetContractSourceCode(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/files/any/1/0x0000000000000000000000001234567890aBcdEF" {
			// Simulate a successful response with sample source code
			sourceCodes := &SourceCodes{
				Status: "success",
				Code: []SourceCode{
					{
						Name:    "Contract.sol",
						Path:    "/path/to/contract.sol",
						Content: "contract MyContract { }",
					},
				},
			}

			err := json.NewEncoder(w).Encode(sourceCodes)
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

	// Define test data
	chainID := 1
	contractAddress := common.HexToAddress("0x1234567890abcdef")
	matchType := MethodMatchTypeAny

	// Call the function to get contract source code
	sourceCodes, err := GetContractSourceCode(client, chainID, contractAddress, matchType)

	// Verify the results
	assert.NoError(t, err, "GetContractSourceCode returned an error")

	expectedSourceCodes := &SourceCodes{
		Status: "success",
		Code: []SourceCode{
			{
				Name:    "Contract.sol",
				Path:    "/path/to/contract.sol",
				Content: "contract MyContract { }",
			},
		},
	}

	assert.Equal(t, expectedSourceCodes, sourceCodes, "GetContractSourceCode returned unexpected source codes")
}
