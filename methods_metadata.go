package sourcify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

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
		if rErr := ToErrorResponse(response); rErr != nil {
			return nil, rErr
		}

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
		if rErr := ToErrorResponse(response); rErr != nil {
			return nil, rErr
		}

		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	body, err := io.ReadAll(response)
	if err != nil {
		return nil, err
	}

	return body, nil
}
