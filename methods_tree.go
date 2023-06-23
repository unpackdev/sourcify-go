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
	// MethodGetFileTreeFullOrPartialMatch represents the API endpoint for getting the file tree with full or partial match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns repository URLs for every file in the source tree for the desired chain and address. Searches for full and partial matches.
	// More information: https://docs.sourcify.dev/docs/api/server/get-file-tree-all/
	MethodGetFileTreeFullOrPartialMatch = Method{
		Name:           "Get File Tree Full or Partial Match",
		URI:            "/files/tree/any/:chain/:address",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/get-file-tree-all/",
		Method:         "GET",
		ParamType:      MethodParamTypeUri,
		RequiredParams: []string{":chain", ":address"},
		Params: []MethodParam{
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

	// MethodGetFileTreeFullMatch represents the API endpoint for getting the file tree with full match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns repository URLs for every file in the source tree for the desired chain and address. Searches only for full matches.
	// More information: https://docs.sourcify.dev/docs/api/server/get-file-tree-full/
	MethodGetFileTreeFullMatch = Method{
		Name:           "Get File Tree Full Match",
		URI:            "/files/tree/:chain/:address",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/get-file-tree-full/",
		Method:         "GET",
		ParamType:      MethodParamTypeUri,
		RequiredParams: []string{":chain", ":address"},
		Params: []MethodParam{
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
)

// FileTree represents the file tree response from the Sourcify service.
type FileTree struct {
	Status string   `json:"status"`
	Files  []string `json:"files"`
}

// GetContractFiles retrieves the repository URLs for every file in the source tree for the given chain ID and contract address.
// The matchType parameter determines whether to search for full matches, partial matches, or any matches.
// It returns the FileTree object containing the status and file URLs, or an error if any.
func GetContractFiles(client *Client, chainId int, contract common.Address, matchType MethodMatchType) (*FileTree, error) {
	var method Method

	switch matchType {
	case MethodMatchTypeFull:
		method = MethodGetFileTreeFullMatch
	case MethodMatchTypePartial:
		method = MethodGetFileTreeFullOrPartialMatch
	case MethodMatchTypeAny:
		method = MethodGetFileTreeFullOrPartialMatch
	default:
		return nil, fmt.Errorf("invalid match type: %s", matchType)
	}

	method.SetParams(
		MethodParam{Key: ":chain", Value: chainId},
		MethodParam{Key: ":address", Value: contract.Hex()},
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

	body, readBodyErr := io.ReadAll(response)
	if readBodyErr != nil {
		return nil, fmt.Errorf("failure to read body: %s", readBodyErr)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	toReturn := &FileTree{}

	if err := json.Unmarshal(body, &toReturn); err != nil {
		// Sometimes, response will not be a JSON object, but an array.
		// In this case, we'll get an error, but we can still return the code.
		// This is a workaround for the Sourcify API.
		// Ugly, but it works.
		if strings.Contains(err.Error(), "cannot unmarshal array into Go value") {
			toReturn.Status = "unknown"
			if err := json.Unmarshal(body, &toReturn.Files); err != nil {
				return nil, err
			}
			return toReturn, nil
		}
		return nil, err
	}

	return toReturn, nil
}
