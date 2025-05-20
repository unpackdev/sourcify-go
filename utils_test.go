package sourcify

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"path/filepath"
	"strconv"
)

// LoadContract loads a contract response from a JSON file in the testdata directory.
// The filename format is chainid:address.json
func LoadContract(chainID int, address common.Address) (*ContractResponse, error) {
	// Create the filename pattern
	filename := fmt.Sprintf("%d:%s.json", chainID, address.Hex())
	
	// Find the testdata directory
	testdataPath := filepath.Join("testdata", filename)
	
	// Read the file
	fileData, err := os.ReadFile(testdataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test data file: %w", err)
	}
	
	// Unmarshal the JSON data
	var contractResponse ContractResponse
	if err := json.Unmarshal(fileData, &contractResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contract data: %w", err)
	}
	
	return &contractResponse, nil
}

// SaveContract saves a contract response to a JSON file in the testdata directory.
// This is useful for creating test data from actual API responses.
func SaveContract(contractResponse *ContractResponse) error {
	// Make sure the testdata directory exists
	testdataPath := "testdata"
	if err := os.MkdirAll(testdataPath, 0755); err != nil {
		return fmt.Errorf("failed to create testdata directory: %w", err)
	}
	
	// Convert chainID to int
	chainID, err := strconv.Atoi(contractResponse.ChainID)
	if err != nil {
		return fmt.Errorf("failed to convert chainID to int: %w", err)
	}
	
	// Create the filename pattern
	filename := fmt.Sprintf("%d:%s.json", chainID, contractResponse.Address)
	filePath := filepath.Join(testdataPath, filename)
	
	// Marshal the contract response to JSON
	data, err := json.MarshalIndent(contractResponse, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal contract data: %w", err)
	}
	
	// Write the file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write test data file: %w", err)
	}
	
	return nil
}
