package sourcify

var (
	// MethodGetFileFromRepositoryFullMatch represents the API endpoint for retrieving staticly served files over the server for full match contract in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches only for full matches.
	// More information: https://docs.sourcify.dev/docs/api/repository/get-file-static/
	MethodGetFileFromRepositoryFullMatch = Method{
		Name:           "Retrieve staticly served files over the server for full match contract",
		URI:            "/repository/contracts/full_match/:chain/:address/:filePath",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/repository/get-file-static/",
		Method:         "GET",
		ParamType:      MethodParamTypeUri,
		RequiredParams: []string{":chain", ":address", ":filePath"},
		Params: []MethodParam{
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

	// MethodGetFileFromRepositoryPartialMatch represents the API endpoint for retrieving staticly served files over the server for partial match contract in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches only for partial matches.
	// More information: https://docs.sourcify.dev/docs/api/repository/get-file-static/
	MethodGetFileFromRepositoryPartialMatch = Method{
		Name:           "Retrieve staticly served files over the server for partial match contract",
		URI:            "/repository/contracts/partial_match/:chain/:address/:filePath",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/repository/get-file-static/",
		Method:         "GET",
		ParamType:      MethodParamTypeUri,
		RequiredParams: []string{":chain", ":address", ":filePath"},
		Params: []MethodParam{
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
)
