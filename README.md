
[![Build Status](https://app.travis-ci.com/txpull/sourcify-go.svg?branch=main)](https://app.travis-ci.com/txpull/sourcify-go)
[![Coverage Status](https://coveralls.io/repos/github/txpull/sourcify-go/badge.svg?branch=main)](https://coveralls.io/github/txpull/sourcify-go?branch=main)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# Sourcify Go

Sourcify Go is a Golang package that provides tools for interacting with the [Sourcify API](https://docs.sourcify.dev/docs/api). 

It allows you to access various API endpoints and perform operations such as checking server status, retrieving chains information, obtaining contract information from supported chains and more.

## Installation

To use Sourcify in your Go project, you can simply import the package:

```go
import "github.com/txpull/sourcify-go"
```

## Usage

### Creating a Client

To interact with the Sourcify API, you need to create a client using the `NewClient` function. You can provide optional configuration options using `ClientOption` functions. For example, you can specify a custom base URL or a custom HTTP client and set retry configuration in case sourcify servers are temporairly unavailable.

```go
client := sourcify.NewClient(
	sourcify.WithHTTPClient(httpClient),
	sourcify.WithBaseURL("https://sourcify.dev/server"),
	sourcify.WithRetryOptions(
		sourcify.WithMaxRetries(3),
		sourcify.WithDelay(2*time.Second),
	),
)
``` 

### Calling Raw API Endpoints

Sourcify provides various API endpoints as `Method` objects. You can call these endpoints using the `CallMethod` function on the client and do your own method parsers if you wish to. 

However, if there are methods that we do not provide yet, you can do something like this to extend
package.

```go
customMethod := sourcify.Method{...}

customMethod.SetParams(
	MethodParam{Key: ":chain", Value: chainId},
)

if err := customMethod.Verify(); err != nil {
	return nil, err
}

response, statusCode, err := client.CallMethod(customMethod)
if err != nil {
	return nil, err
}

// Close the io.ReadCloser interface.
// This is important as CallMethod is NOT closing the response body!
// You'll have memory leaks if you don't do this!
defer response.Close()

// Process the response
```

### Supported API Endpoints

Sourcify provides the following API endpoints that you can call that are currently supported by this package:

- `MethodHealth`: Check the server status. [More information](https://docs.sourcify.dev/docs/api/server/check-server-status/)
- `MethodGetChains`: Retrieve the chains (networks) added to Sourcify. [More information](https://docs.sourcify.dev/docs/api/server/retrieve-chains/)
- `MethodCheckByAddresses`: Check if contracts with the desired chain and addresses are verified and in the repository. [More information](https://docs.sourcify.dev/docs/api/server/check-by-addresses/)
- `MethodCheckAllByAddresses`: Check if contracts with the desired chain and addresses are verified and in the repository. [More information](https://docs.sourcify.dev/docs/api/server/check-all-by-addresses/)
- `MethodGetFileTreeFullOrPartialMatch`: Get the file tree with full or partial match for the desired chain and address. [More information](https://docs.sourcify.dev/docs/api/server/get-file-tree-all/)
- `MethodGetFileTreeFullMatch`: Get the file tree with full match for the desired chain and address. [More information](https://docs.sourcify.dev/docs/api/server/get-file-tree-full/)
- `MethodSourceFilesFullOrPartialMatch`: Get the source files with full or partial match for the desired chain and address, including metadata.json. [More information](https://docs.sourcify.dev/docs/api/server/get-source-files-all/)
- `MethodSourceFilesFullMatch`: Get the source files with full match for the desired chain and address, including metadata.json. [More information](https://docs.sourcify.dev/docs/api/server/get-source-files-full/)
- `MethodGetContractAddressesFullOrPartialMatch`: Get the verified contract addresses for the chain with full or partial match. [More information](https://docs.sourcify.dev/docs/api/server/get-contract-addresses-all/)
- `MethodGetFileFromRepositoryFullMatch`: Retrieve statically served files over the server for full match contract. [More information](https://docs.sourcify.dev/docs/api/repository/get-file-static/)
- `MethodGetFileFromRepositoryPartialMatch`: Retrieve statically served files over the server for partial match contract. [More information](https://docs.sourcify.dev/docs/api/repository/get-file-static/)

For more information on each endpoint, including the parameters they require and the expected responses, refer to the [Sourcify API documentation](https://docs.sourcify.dev/docs/api).

## Examples

You can find endpoint examples under the [examples](/examples) directory.

## Contributing

Contributions to Sourcify are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on GitHub.


## License

Sourcify is released under the Apache 2.0 License. See the [LICENSE](LICENSE) file for more details.

