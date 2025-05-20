package sourcify

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"net/http"
	"strings"
	"time"
)

var (
	// MethodGetContractByChainIdAndAddress represents the API endpoint for getting the contract addresses with full or partial match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches only for full matches.
	// More information: https://docs.sourcify.dev/docs/api/#/Contract%20Lookup
	MethodGetContractByChainIdAndAddress = Method{
		Name:           "Get verified contract addresses for the chain full or partial match",
		URI:            "/v2/contract/:chain/:address",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/#/Contract%20Lookup",
		Method:         "GET",
		ParamType:      MethodParamTypeUriAndQueryString,
		RequiredParams: []string{":chain", ":address"},
		Params: []MethodParam{
			{
				Key:   "fields",
				Value: "",
			},
			{
				Key:   "omit",
				Value: ""},
		},
	}

	MethodGetContractByChainId = Method{
		Name:           "Get verified contract addresses for the chain full or partial match",
		URI:            "/v2/contracts/:chain",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/#/Contract%20Lookup",
		Method:         "GET",
		ParamType:      MethodParamTypeUriAndQueryString,
		RequiredParams: []string{":chain"},
		Params: []MethodParam{
			{
				Key:   "fields",
				Value: "",
			},
		},
	}
)

// ContractResponse represents the response from the Sourcify API when retrieving contract information
type ContractResponse struct {
	Abi              []ABIEntry      `json:"abi"`
	Address          string          `json:"address"`
	ChainID          string          `json:"chainId"`
	Compilation      Compilation     `json:"compilation"`
	CreationBytecode Bytecode        `json:"creationBytecode"`
	CreationMatch    string          `json:"creationMatch"`
	Deployment       Deployment      `json:"deployment"`
	DevDoc           DevDoc          `json:"devdoc"`
	Match            string          `json:"match"`
	MatchID          string          `json:"matchId"`
	Metadata         Metadata        `json:"metadata"`
	ProxyResolution  ProxyResolution `json:"proxyResolution"`
	RuntimeBytecode  Bytecode        `json:"runtimeBytecode"`
	RuntimeMatch     string          `json:"runtimeMatch"`
	SourceIds        SourceIds       `json:"sourceIds"`
	Sources          Sources         `json:"sources"`
	StdJSONInput     StdJSONInput    `json:"stdJsonInput"`
	StdJSONOutput    StdJSONOutput   `json:"stdJsonOutput"`
	StorageLayout    *StorageLayout  `json:"storageLayout"`
	UserDoc          UserDoc         `json:"userdoc"`
	VerifiedAt       time.Time       `json:"verifiedAt"`
}

// ContractBaseResponse represents the response from the Sourcify API when retrieving contract information
type ContractBaseResponse struct {
	Address       string    `json:"address"`
	ChainID       string    `json:"chainId"`
	CreationMatch string    `json:"creationMatch"`
	Match         string    `json:"match"`
	MatchID       string    `json:"matchId"`
	RuntimeMatch  string    `json:"runtimeMatch"`
	VerifiedAt    time.Time `json:"verifiedAt"`
}

type ContractsResponse struct {
	Results []ContractBaseResponse `json:"results"`
}

// GetContractsByChainId retrieves the available verified contract addresses for the given chain ID.
func GetContractsByChainId(client *Client, chainId int, sort string, afterMatchId string, limit int) (*ContractsResponse, error) {
	method := MethodGetContractByChainId

	method.SetParams(
		MethodParam{Key: ":chain", Value: chainId},
		MethodParam{Key: "sort", Value: sort},
		MethodParam{Key: "afterMatchId", Value: afterMatchId},
		MethodParam{Key: "limit", Value: limit},
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

	var toReturn *ContractsResponse
	if jdErr := json.NewDecoder(response).Decode(&toReturn); jdErr != nil {
		return nil, jdErr
	}

	return toReturn, nil
}

// GetContractByChainIdAndAddress retrieves the available verified contract addresses for the given chain ID.
func GetContractByChainIdAndAddress(client *Client, chainId int, address common.Address, fields []string, omit []string) (*ContractResponse, error) {
	method := MethodGetContractByChainIdAndAddress

	// Omit and fields cannot co-exist together
	if len(omit) == 0 && len(fields) == 0 {
		fields = []string{"all"}
	}

	pFields := strings.Join(fields, ",")
	pOmit := strings.Join(omit, ",")

	method.SetParams(
		MethodParam{Key: ":chain", Value: chainId},
		MethodParam{Key: ":address", Value: address.Hex()},
		MethodParam{Key: "fields", Value: pFields},
		MethodParam{Key: "omit", Value: pOmit},
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

	var toReturn ContractResponse
	if jdErr := json.NewDecoder(response).Decode(&toReturn); jdErr != nil {
		return nil, jdErr
	}

	return &toReturn, nil
}
