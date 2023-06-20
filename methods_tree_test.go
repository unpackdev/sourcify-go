package sourcify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestGetContractFiles(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/files/tree/any/1/0x0000000000000000000000001234567890aBcdEF" {
			// Simulate a successful response with sample file tree
			fileTree := &FileTree{
				Status: "success",
				Files:  []string{"/path/to/file1.sol", "/path/to/file2.sol"},
			}

			err := json.NewEncoder(w).Encode(fileTree)
			if err != nil {
				t.Errorf("failed to encode mock file tree: %v", err)
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

	// Call the function to get contract files
	fileTree, err := GetContractFiles(client, chainID, contractAddress, matchType)

	// Verify the results
	assert.NoError(t, err, "GetContractFiles returned an error")

	expectedFileTree := &FileTree{
		Status: "success",
		Files:  []string{"/path/to/file1.sol", "/path/to/file2.sol"},
	}

	assert.Equal(t, expectedFileTree, fileTree, "GetContractFiles returned unexpected file tree")
}

func TestGetContractFiles_Error(t *testing.T) {
	// Create a mock HTTP server that always returns 404 Not Found
	mockServer := httptest.NewServer(http.NotFoundHandler())
	defer mockServer.Close()

	// Create a client for the mock server
	client := NewClient(WithBaseURL(mockServer.URL))

	// Define test data
	chainID := 1
	contractAddress := common.HexToAddress("0x1234567890abcdef")
	matchType := MethodMatchTypeAny

	// Call the function to get contract files
	fileTree, err := GetContractFiles(client, chainID, contractAddress, matchType)

	// Verify the results
	assert.Error(t, err, "GetContractFiles should return an error")
	assert.Nil(t, fileTree, "File tree should be nil")
}
