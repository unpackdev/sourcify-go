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
	// MethodSourceFilesFullOrPartialMatch represents the API endpoint for getting the source files with full or partial match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches for full and partial matches.
	// More information: https://docs.sourcify.dev/docs/api/server/get-source-files-all/
	MethodSourceFilesFullOrPartialMatch = Method{
		Name:           "Get source files for the address full or partial match",
		URI:            "/files/any/:chain/:address",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/get-source-files-all/",
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

	// MethodSourceFilesFullMatch represents the API endpoint for getting the source files with full match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches only for full matches.
	// More information: https://docs.sourcify.dev/docs/api/server/get-source-files-full/
	MethodSourceFilesFullMatch = Method{
		Name:           "Get source files for the address full match",
		URI:            "/files/:chain/:address",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/server/get-source-files-full/",
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

// SourceCode represents the source code details for a file.
type SourceCode struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

// SourceCodes represents the source code details for multiple files.
type SourceCodes struct {
	Status string       `json:"status"`
	Code   []SourceCode `json:"files"`
}

// GetContractSourceCode retrieves the source code files for a contract with the given chain ID and address, based on the match type.
// It makes an API request to the Sourcify service and returns the source code details as a SourceCodes object.
func GetContractSourceCode(client *Client, chainId int, contract common.Address, matchType MethodMatchType) (*SourceCodes, error) {
	var method Method

	switch matchType {
	case MethodMatchTypeFull:
		method = MethodSourceFilesFullMatch
	case MethodMatchTypePartial:
		method = MethodSourceFilesFullOrPartialMatch
	case MethodMatchTypeAny:
		method = MethodSourceFilesFullOrPartialMatch
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

	toReturn := &SourceCodes{}

	if err := json.Unmarshal(body, &toReturn); err != nil {
		// Sometimes, response will not be a JSON object, but an array.
		// In this case, we'll get an error, but we can still return the code.
		// This is a workaround for the Sourcify API.
		// Ugly, but it works.
		if strings.Contains(err.Error(), "cannot unmarshal array into Go value") {
			toReturn.Status = "unknown"
			if err := json.Unmarshal(body, &toReturn.Code); err != nil {
				return nil, err
			}
			return toReturn, nil
		}
		return nil, err
	}

	return toReturn, nil
}
