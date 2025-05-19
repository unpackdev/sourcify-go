package sourcify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

// Metadata represents the top-level structure for compiler metadata
// for Ethereum smart contracts.
type Metadata struct {
	Compiler Compiler `json:"compiler"` // Compiler contains information about the compiler used.
	Language string   `json:"language"` // Language of the source code
	Output   Output   `json:"output"`   // Output represents details of the compiled code.
	Settings Settings `json:"settings"` // Settings represent the compiler settings used.
	Sources  Sources  `json:"sources"`  // Sources represents the details of the source code.
	Version  int      `json:"version"`  // Version of the metadata.
}

// Compiler provides information about the compiler used for the smart contract.
type Compiler struct {
	Version string `json:"version"` // Compiler version
}

// Output contains details about the output of the compiled code.
type Output struct {
	Abi     []Abi   `json:"abi"`     // Abi represents the Application Binary Interface (ABI) of the compiled code.
	Devdoc  Devdoc  `json:"devdoc"`  // Devdoc represents the developer documentation.
	Userdoc Userdoc `json:"userdoc"` // Userdoc represents the user documentation.
}

// Abi holds the Application Binary Interface (ABI) of the compiled code.
type Abi struct {
	Inputs          []any          `json:"inputs"`          // Input parameters of the functions
	StateMutability string         `json:"stateMutability"` // State of mutability of the functions
	Type            string         `json:"type"`            // Type of the ABI entry
	Anonymous       bool           `json:"anonymous"`       // Whether the function is anonymous
	Name            string         `json:"name"`            // Name of the function
	Outputs         []OutputDetail `json:"outputs"`         // Output parameters of the functions
}

// OutputDetail holds information about the output parameters of the functions.
type OutputDetail struct {
	InternalType string `json:"internalType"` // Internal type of the parameter
	Name         string `json:"name"`         // Name of the parameter
	Type         string `json:"type"`         // Type of the parameter
}

// Devdoc provides details about the developer documentation.
type Devdoc struct {
	DevMethods map[string]DevMethod `json:"methods"` // Mapping of function signatures to their documentation
}

// DevMethod contains information about a method in the developer documentation.
type DevMethod struct {
	Details string `json:"details"` // Details about the method
}

// Userdoc provides information about the user documentation.
type Userdoc struct {
	Methods map[string]DevMethod `json:"methods"` // Mapping of function signatures to their documentation
}

// Libraries represent the libraries used in the source code.
type Libraries struct {
}

// MetadataDetail provides additional metadata.
type MetadataDetail struct {
	BytecodeHash string `json:"bytecodeHash"` // Hash of the bytecode
}

// Sources provides details about the source code.
type Sources map[string]SourceDetails

// SourceDetails holds the details of the main contract source code.
type SourceDetails struct {
	Keccak256 string   `json:"keccak256"` // Hash of the source code
	License   string   `json:"license"`   // License of the source code
	Urls      []string `json:"urls"`      // URLs of the source code
}

// GetContractMetadata fetches the metadata of a contract from a given client,
// chain ID, contract address, and match type. It returns a Metadata object and
// an error, if any. This function is primarily used to fetch and parse metadata
// from smart contracts.
func GetContractMetadata(client *Client, chainId int, contract common.Address, matchType MethodMatchType) (*Metadata, error) {
	var method Method

	switch matchType {
	case MethodMatchTypeFull:
		method = MethodGetFileFromRepositoryFullMatch
	case MethodMatchTypePartial:
		method = MethodGetFileFromRepositoryPartialMatch
	case MethodMatchTypeAny:
		return nil, fmt.Errorf("type: %s is not implemented", matchType)
	default:
		return nil, fmt.Errorf("invalid match type: %s", matchType)
	}

	method.SetParams(
		MethodParam{Key: ":chain", Value: chainId},
		MethodParam{Key: ":address", Value: contract.Hex()},
		MethodParam{Key: ":filePath", Value: "metadata.json"},
	)

	if err := method.Verify(); err != nil {
		return nil, err
	}

	response, statusCode, err := client.CallMethod(method)
	if err != nil {
		return nil, err
	}

	// Close the io.ReadCloser interface.
	// This is important as CallMethod is NOT closing the response body!
	// You'll have memory leaks if you don't do this!
	defer response.Close()

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	var toReturn Metadata
	if err := json.NewDecoder(response).Decode(&toReturn); err != nil {
		return nil, err
	}

	return &toReturn, nil
}

// GetContractMetadataAsBytes retrieves the metadata of a smart contract as a byte slice.
// The client parameter is a pointer to a Client instance used to make the API call.
// The chainId parameter is the ID of the blockchain network where the smart contract resides.
// The contract parameter is the address of the smart contract.
// The matchType parameter determines the type of method matching to use when retrieving the contract metadata.
// It returns a byte slice containing the contract metadata, or an error if there was an issue sending the request, decoding the response, or if an invalid matchType was provided.
//
// The MethodMatchType enum is used to determine the type of method matching to use:
// - MethodMatchTypeFull: Use full method matching. This will only return a match if the entire method signature matches.
// - MethodMatchTypePartial: Use partial method matching. This will return a match if any part of the method signature matches.
// - MethodMatchTypeAny: This match type is not implemented and will return an error.
//
// This function will send a request to the API endpoint specified in the client parameter, using the method determined by the matchType parameter.
// The method will be set with the chainId, contract address, and file path parameters.
// If the method fails verification, an error will be returned.
// If the API call is successful, the response body will be read and returned as a byte slice.
// If the status code of the response is not 200 OK, an error will be returned.
func GetContractMetadataAsBytes(client *Client, chainId int, contract common.Address, matchType MethodMatchType) ([]byte, error) {
	var method Method

	switch matchType {
	case MethodMatchTypeFull:
		method = MethodGetFileFromRepositoryFullMatch
	case MethodMatchTypePartial:
		method = MethodGetFileFromRepositoryPartialMatch
	case MethodMatchTypeAny:
		return nil, fmt.Errorf("type: %s is not implemented", matchType)
	default:
		return nil, fmt.Errorf("invalid match type: %s", matchType)
	}

	method.SetParams(
		MethodParam{Key: ":chain", Value: chainId},
		MethodParam{Key: ":address", Value: contract.Hex()},
		MethodParam{Key: ":filePath", Value: "metadata.json"},
	)

	if err := method.Verify(); err != nil {
		return nil, err
	}

	response, statusCode, err := client.CallMethod(method)
	if err != nil {
		return nil, err
	}

	// Close the io.ReadCloser interface.
	// This is important as CallMethod is NOT closing the response body!
	// You'll have memory leaks if you don't do this!
	defer response.Close()

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	body, err := io.ReadAll(response)
	if err != nil {
		return nil, err
	}

	return body, nil
}
