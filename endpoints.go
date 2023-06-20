// Package sourcify provides tools for interacting with the Sourcify API.
package sourcify

import (
	"fmt"
	"reflect"
	"strings"
)

type EndpointParamType int

const (
	// EndpointParamTypeUri denotes the type of parameter which is part of the URI.
	EndpointParamTypeUri EndpointParamType = iota // 0

	// EndpointParamTypeQueryString denotes the type of parameter which is part of the query string.
	EndpointParamTypeQueryString // 1
)

// Endpoint represents an API endpoint in the Sourcify service.
// It includes the name, the HTTP method, the URI, and any necessary parameters for requests to this endpoint.
type Endpoint struct {
	Name           string
	Method         string
	URI            string
	MoreInfo       string
	ParamType      EndpointParamType
	RequiredParams []string
	Params         []EndpointParam
}

// EndpointParam represents a parameter key-value pair.
type EndpointParam struct {
	Key   string
	Value interface{}
}

// GetParams returns a slice of the parameters for the API endpoint.
func (e Endpoint) GetParams() []EndpointParam {
	return e.Params
}

// SetParams allows setting parameters for the API endpoint using a variadic list of EndpointParam values.
func (e *Endpoint) SetParams(params ...EndpointParam) {
	e.Params = params
}

// Verify checks if all the required parameters for the API endpoint are provided.
// It returns an error if any of the required parameters is missing.
func (e Endpoint) Verify() error {
	for _, param := range e.RequiredParams {
		found := false
		for _, endpointParam := range e.Params {
			if param == endpointParam.Key {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing required parameter: %s", param)
		}
	}
	return nil
}

// ParseUri constructs a URL for the API endpoint, including any necessary query parameters or URI replacements.
// It can handle parameters of type string, int, []string, and []int. For []string and []int,
// the individual elements are joined with commas for the query string.
// Other types of parameters will trigger an error.
func (e Endpoint) ParseUri() (string, error) {
	switch e.ParamType {
	case EndpointParamTypeQueryString:
		var toReturn string

		// Add the parameters to the URL
		params := []string{}
		for _, param := range e.Params {
			switch v := param.Value.(type) {
			case []string:
				if len(v) > 0 {
					params = append(params, fmt.Sprintf("%s=%s", param.Key, strings.Join(v, ",")))
				}
			case []int:
				if len(v) > 0 {
					strs := []string{}
					for _, i := range v {
						strs = append(strs, fmt.Sprintf("%d", i))
					}
					params = append(params, fmt.Sprintf("%s=%s", param.Key, strings.Join(strs, ",")))
				}
			case string:
				if v != "" {
					params = append(params, fmt.Sprintf("%s=%s", param.Key, v))
				}
			case int:
				params = append(params, fmt.Sprintf("%s=%d", param.Key, v))
			default:
				// Return an error when encountering unsupported parameter type
				return "", ErrInvalidParamType(reflect.TypeOf(v).String())
			}
		}

		if len(params) > 0 {
			toReturn = fmt.Sprintf("%s?%s", toReturn, strings.Join(params, "&"))
		}

		return toReturn, nil

	case EndpointParamTypeUri:
		toReturn := e.URI
		for _, param := range e.Params {
			switch v := param.Value.(type) {
			case string, int:
				toReturn = strings.ReplaceAll(toReturn, ":"+param.Key, fmt.Sprintf("%v", v))
			default:
				// Return an error when encountering unsupported parameter type
				return "", ErrInvalidParamType(reflect.TypeOf(v).String())
			}
		}
		return toReturn, nil

	default:
		return "", fmt.Errorf("invalid EndpointParamType: %v", e.ParamType)
	}
}

var (
	// EndpointCheckByAddresses represents the API endpoint for checking by addresses in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Checks if contract with the desired chain and address is verified and in the repository.
	// More information: https://docs.sourcify.dev/docs/api/server/check-by-addresses/
	EndpointCheckByAddresses = Endpoint{
		Name:           "Check By Addresses",
		URI:            "/check-by-addresses",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/check-by-addresses/",
		Method:         "GET",
		ParamType:      EndpointParamTypeQueryString,
		RequiredParams: []string{"addresses", "chainIds"},
		Params: []EndpointParam{
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

	// EndpointCheckAllByAddresses represents the API endpoint for checking all addresses in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Checks if contract with the desired chain and address is verified and in the repository.
	// More information: https://docs.sourcify.dev/docs/api/server/check-all-by-addresses/
	EndpointCheckAllByAddresses = Endpoint{
		Name:           "Check All By Addresses",
		URI:            "/check-all-by-addresses",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/check-all-by-addresses/",
		Method:         "GET",
		ParamType:      EndpointParamTypeQueryString,
		RequiredParams: []string{"addresses", "chainIds"},
		Params: []EndpointParam{
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

	// EndpointGetFileTreeFullOrPartialMatch represents the API endpoint for getting the file tree with full or partial match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns repository URLs for every file in the source tree for the desired chain and address. Searches for full and partial matches.
	// More information: https://docs.sourcify.dev/docs/api/server/get-file-tree-all/
	EndpointGetFileTreeFullOrPartialMatch = Endpoint{
		Name:           "Get File Tree Full or Partial Match",
		URI:            "/files/tree/any/:chain/:address",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/get-file-tree-all/",
		Method:         "GET",
		ParamType:      EndpointParamTypeUri,
		RequiredParams: []string{":chain", ":address"},
		Params: []EndpointParam{
			{
				Key:   ":chain",
				Value: "",
			},
			{
				Key:   ":address",
				Value: "",
			},
		},
	}

	// EndpointGetFileTreeFullMatch represents the API endpoint for getting the file tree with full match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns repository URLs for every file in the source tree for the desired chain and address. Searches only for full matches.
	// More information: https://docs.sourcify.dev/docs/api/server/get-file-tree-full/
	EndpointGetFileTreeFullMatch = Endpoint{
		Name:           "Get File Tree Full Match",
		URI:            "/files/tree/:chain/:address",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/get-file-tree-full/",
		Method:         "GET",
		ParamType:      EndpointParamTypeUri,
		RequiredParams: []string{":chain", ":address"},
		Params: []EndpointParam{
			{
				Key:   ":chain",
				Value: "",
			},
			{
				Key:   ":address",
				Value: "",
			},
		},
	}

	// EndpointSourceFilesFullOrPartialMatch represents the API endpoint for getting the source files with full or partial match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches for full and partial matches.
	// More information: https://docs.sourcify.dev/docs/api/server/get-source-files-all/
	EndpointSourceFilesFullOrPartialMatch = Endpoint{
		Name:           "Get source files for the address full or partial match",
		URI:            "/files/any/:chain/:address",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/get-source-files-all/",
		Method:         "GET",
		ParamType:      EndpointParamTypeUri,
		RequiredParams: []string{":chain", ":address"},
		Params: []EndpointParam{
			{
				Key:   ":chain",
				Value: "",
			},
			{
				Key:   ":address",
				Value: "",
			},
		},
	}

	// EndpointSourceFilesFullMatch represents the API endpoint for getting the source files with full match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches only for full matches.
	// More information: https://docs.sourcify.dev/docs/api/server/get-source-files-full/
	EndpointSourceFilesFullMatch = Endpoint{
		Name:           "Get source files for the address full match",
		URI:            "/files/:chain/:address",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/get-source-files-full/",
		Method:         "GET",
		ParamType:      EndpointParamTypeUri,
		RequiredParams: []string{":chain", ":address"},
		Params: []EndpointParam{
			{
				Key:   ":chain",
				Value: "",
			},
			{
				Key:   ":address",
				Value: "",
			},
		},
	}

	// EndpointGetContractAddressesFullOrPartialMatch represents the API endpoint for getting the contract addresses with full or partial match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches only for full matches.
	// More information: https://docs.sourcify.dev/docs/api/server/get-contract-addresses-all/
	EndpointGetContractAddressesFullOrPartialMatch = Endpoint{
		Name:           "Get verified contract addresses for the chain full or partial match",
		URI:            "/files/contracts/:chain",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/get-contract-addresses-all/",
		Method:         "GET",
		ParamType:      EndpointParamTypeUri,
		RequiredParams: []string{":chain"},
		Params: []EndpointParam{
			{
				Key:   ":chain",
				Value: "",
			},
		},
	}

	// EndpointGetFileFromRepositoryFullMatch represents the API endpoint for retrieving staticly served files over the server for full match contract in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches only for full matches.
	// More information: https://docs.sourcify.dev/docs/api/repository/get-file-static/
	EndpointGetFileFromRepositoryFullMatch = Endpoint{
		Name:           "Retrieve staticly served files over the server for full match contract",
		URI:            "/repository/contracts/full_match/:chain/:address/:filePath",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/repository/get-file-static/",
		Method:         "GET",
		ParamType:      EndpointParamTypeUri,
		RequiredParams: []string{":chain", ":address", ":filePath"},
		Params: []EndpointParam{
			{
				Key:   ":chain",
				Value: "",
			},
			{
				Key:   ":address",
				Value: "",
			},
			{
				Key:   ":filePath",
				Value: "",
			},
		},
	}

	// EndpointGetFileFromRepositoryPartialMatch represents the API endpoint for retrieving staticly served files over the server for partial match contract in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches only for partial matches.
	// More information: https://docs.sourcify.dev/docs/api/repository/get-file-static/
	EndpointGetFileFromRepositoryPartialMatch = Endpoint{
		Name:           "Retrieve staticly served files over the server for partial match contract",
		URI:            "/repository/contracts/partial_match/:chain/:address/:filePath",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/repository/get-file-static/",
		Method:         "GET",
		ParamType:      EndpointParamTypeUri,
		RequiredParams: []string{":chain", ":address", ":filePath"},
		Params: []EndpointParam{
			{
				Key:   ":chain",
				Value: "",
			},
			{
				Key:   ":address",
				Value: "",
			},
			{
				Key:   ":filePath",
				Value: "",
			},
		},
	}

	// EndpointGetChains represents the API endpoint for getting the chains (networks) added to Sourcify in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns the chains (networks) added to Sourcify. Contains both supported, unsupported, monitored, and unmonitored chains.
	// More information: https://docs.sourcify.dev/docs/api/chains/
	EndpointGetChains = Endpoint{
		Name:           "Retrieve staticly served files over the server for partial match contract",
		URI:            "/chains",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/chains/",
		Method:         "GET",
		ParamType:      EndpointParamTypeUri,
		RequiredParams: []string{},
		Params:         []EndpointParam{},
	}

	// EndpointHealth represents the API endpoint for checking the server status in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Ping the server and see if it is alive and ready for requests.
	// More information: https://docs.sourcify.dev/docs/api/health/
	EndpointHealth = Endpoint{
		Name:           "Show Server Status",
		URI:            "/health",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/health/",
		Method:         "GET",
		ParamType:      EndpointParamTypeUri,
		RequiredParams: []string{},
		Params:         []EndpointParam{},
	}
)
