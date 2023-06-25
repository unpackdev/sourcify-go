package sourcify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// MethodCheckByAddresses represents the API endpoint for checking by addresses in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Checks if contract with the desired chain and address is verified and in the repository.
	// More information: https://docs.sourcify.dev/docs/api/server/check-by-addresses/
	MethodCheckByAddresses = Method{
		Name:           "Check By Addresses",
		URI:            "/check-by-addresses",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/check-by-addresses/",
		Method:         "GET",
		ParamType:      MethodParamTypeQueryString,
		RequiredParams: []string{"addresses", "chainIds"},
		Params: []MethodParam{
			{
				Key:   "addresses",
				Value: []string{},
			},
			{
				Key:   "chainIds",
				Value: []int{},
			},
		},
	}

	// MethodCheckAllByAddresses represents the API endpoint for checking all addresses in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Checks if contract with the desired chain and address is verified and in the repository.
	// More information: https://docs.sourcify.dev/docs/api/server/check-all-by-addresses/
	MethodCheckAllByAddresses = Method{
		Name:           "Check All By Addresses",
		URI:            "/check-all-by-addresses",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/check-all-by-addresses/",
		Method:         "GET",
		ParamType:      MethodParamTypeQueryString,
		RequiredParams: []string{"addresses", "chainIds"},
		Params: []MethodParam{
			{
				Key:   "addresses",
				Value: []string{},
			},
			{
				Key:   "chainIds",
				Value: []int{},
			},
		},
	}
)

// CheckContractAddress represents the contract address and associated chain IDs and statuses.
type CheckContractAddress struct {
	Address  common.Address `json:"address"`  // The contract address.
	Status   string         `json:"status"`   // The status of the contract.
	ChainIDs []string       `json:"chainIds"` // The chain ID.
}

// CheckContractAddressMore represents the contract address and associated chain IDs and statuses.
type CheckContractAddressMore struct {
	Address common.Address                 `json:"address"`  // The contract address.
	Info    []CheckContractAddressMoreInfo `json:"chainIds"` // The chain ID.
}

// CheckContractAddressMoreInfo represents the contract address and associated chain IDs and statuses.
type CheckContractAddressMoreInfo struct {
	Status  string `json:"status"`  // The status of the contract.
	ChainID string `json:"chainId"` // The chain ID.
}

// CheckContractByAddresses retrieves the available verified contract addresses for the given chain ID.
func CheckContractByAddresses(client *Client, addresses []string, chainIds []int, matchType MethodMatchType) ([]*CheckContractAddress, error) {
	var method Method

	switch matchType {
	case MethodMatchTypeFull:
		method = MethodCheckByAddresses
	case MethodMatchTypePartial:
		method = MethodCheckAllByAddresses
	case MethodMatchTypeAny:
		method = MethodCheckAllByAddresses
	default:
		return nil, fmt.Errorf("invalid match type: %s", matchType)
	}

	method.SetParams(
		MethodParam{Key: "addresses", Value: addresses},
		MethodParam{Key: "chainIds", Value: chainIds},
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

	var toReturn []*CheckContractAddress
	if err := json.Unmarshal(body, &toReturn); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal object into Go struct field CheckContractAddress.chainIds") {
			var toReturnMore []*CheckContractAddressMore
			if err := json.Unmarshal(body, &toReturnMore); err != nil {
				return nil, err
			}

			for _, v := range toReturnMore {
				for _, info := range v.Info {
					toReturn = append(toReturn, &CheckContractAddress{
						Address:  v.Address,
						Status:   info.Status,
						ChainIDs: []string{info.ChainID},
					})
				}
			}

			return toReturn, nil
		}
		return nil, err
	}

	return toReturn, nil
}
