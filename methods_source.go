package sourcify

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
