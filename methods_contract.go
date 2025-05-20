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
	// MethodGetContractByChainIdAndAddress represents the API endpoint for retrieving contract details by chain ID and address.
	// It returns all verified sources from the repository for the specified contract address and chain, including metadata.json.
	// This endpoint searches only for full matches.
	// HTTP Method: GET
	// URI: /v2/contract/:chain/:address
	// Documentation: https://docs.sourcify.dev/docs/api/#/Contract%20Lookup/get-contract
	MethodGetContractByChainIdAndAddress = Method{
		Name:           "Get contract by chain id and address",
		URI:            "/v2/contract/:chain/:address",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/#/Contract%20Lookup/get-contract",
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

	// MethodGetContractByChainId represents the API endpoint for listing all verified contracts for a specific chain.
	// It returns all verified contract addresses and basic details for the specified chain ID.
	// HTTP Method: GET
	// URI: /v2/contracts/:chain
	// Documentation: https://docs.sourcify.dev/docs/api/#/Contract%20Lookup/get-v2-contracts-chainId
	MethodGetContractByChainId = Method{
		Name:           "List through all contracts for the chain",
		URI:            "/v2/contracts/:chain",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/#/Contract%20Lookup/get-v2-contracts-chainId",
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

// ContractResponse represents the detailed response from the Sourcify API when retrieving complete contract information.
// Contains full contract data including ABI, bytecode, sources, and metadata.
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

// ContractBaseResponse represents the basic response from the Sourcify API when retrieving contract information.
// Contains only essential contract verification details without the full source code or ABI.
type ContractBaseResponse struct {
	Address       string    `json:"address"`
	ChainID       string    `json:"chainId"`
	CreationMatch string    `json:"creationMatch"`
	Match         string    `json:"match"`
	MatchID       string    `json:"matchId"`
	RuntimeMatch  string    `json:"runtimeMatch"`
	VerifiedAt    time.Time `json:"verifiedAt"`
}

// ContractsResponse wraps a collection of ContractBaseResponse objects returned when listing multiple contracts.
type ContractsResponse struct {
	Results []ContractBaseResponse `json:"results"`
}

// GetContractsByChainId retrieves a paginated list of verified contract addresses for the given chain ID.
// Parameters:
//   - client: The Sourcify API client
//   - chainId: The blockchain network ID
//   - sort: Sorting option for results
//   - afterMatchId: Pagination parameter; returns results after this match ID
//   - limit: Maximum number of results to return
// Returns a ContractsResponse containing basic information about each contract or an error.
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

// GetContractByChainIdAndAddress retrieves the complete details of a verified contract by its chain ID and address.
// Parameters:
//   - client: The Sourcify API client
//   - chainId: The blockchain network ID
//   - address: The Ethereum contract address
//   - fields: Specific fields to include in the response (use []string{"all"} for complete data)
//   - omit: Fields to exclude from the response
// Note: fields and omit parameters are mutually exclusive; if both are empty, fields defaults to ["all"].
// Returns a ContractResponse containing detailed contract information or an error.
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
