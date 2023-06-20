package sourcify

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
