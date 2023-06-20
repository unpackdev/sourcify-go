package sourcify

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// MethodGetContractAddressesFullOrPartialMatch represents the API endpoint for getting the contract addresses with full or partial match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches only for full matches.
	// More information: https://docs.sourcify.dev/docs/api/server/get-contract-addresses-all/
	MethodGetContractAddressesFullOrPartialMatch = Method{
		Name:           "Get verified contract addresses for the chain full or partial match",
		URI:            "/files/contracts/:chain",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/get-contract-addresses-all/",
		Method:         "GET",
		ParamType:      MethodParamTypeUri,
		RequiredParams: []string{":chain"},
		Params:         []MethodParam{},
	}
)

// VerifiedContractAddresses represents the structure for the verified contract addresses response.
type VerifiedContractAddresses struct {
	Full    []common.Address `json:"full"`
	Partial []common.Address `json:"partial"`
}

// GetAvailableContractAddresses retrieves the available verified contract addresses for the given chain ID.
func GetAvailableContractAddresses(client *Client, chainId int) (*VerifiedContractAddresses, error) {
	method := MethodGetContractAddressesFullOrPartialMatch
	method.SetParams(
		MethodParam{Key: ":chain", Value: chainId},
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

	var toReturn VerifiedContractAddresses
	if err := json.NewDecoder(response).Decode(&toReturn); err != nil {
		return nil, err
	}

	return &toReturn, nil
}
