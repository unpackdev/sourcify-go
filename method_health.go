package sourcify

import (
	"net/http"
)

// MethodHealth represents the API endpoint for checking the server status in the Sourcify service.
// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
// Ping the server and see if it is alive and ready for requests.
// More information: https://docs.sourcify.dev/docs/api/health/
var MethodHealth = Method{
	Name:           "Show Server Status",
	URI:            "/health",
	MoreInfo:       "https://docs.sourcify.dev/docs/api/health/",
	Method:         "GET",
	ParamType:      MethodParamTypeUri,
	RequiredParams: []string{},
	Params:         []MethodParam{},
}

// GetHealth checks the server status by calling the MethodHealth endpoint using the provided client.
// It returns a boolean indicating if the server is healthy and an error if any occurred during the request.
func GetHealth(client *Client) (bool, error) {
	response, statusCode, err := client.CallMethod(MethodHealth)
	if err != nil {
		return false, err
	}

	// Close the io.ReadCloser interface
	// This is important as CallMethod is NOT closing the response body!
	// You'll have memory leaks if you don't do this!
	defer response.Close()

	if statusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}
