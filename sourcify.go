// Package sourcify provides a Go client for interacting with the Sourcify API.
//
// The Sourcify API allows developers to fetch metadata and verify contract addresses for sourcify supported smart contracts.
// This package includes methods for retrieving contract metadata, checking contract verification status by addresses,
// and accessing compiler information and output details.
//
// To use this package, create a new client using the NewClient function, configure it with optional options,
// and then call the desired methods to interact with the Sourcify API.
//
// The package also includes various types and structures representing contract metadata, compiler information,
// output details, and client options. These types can be used to parse and work with the data returned by the API.
//
// Example usage:
//
//	client := sourcify.NewClient(
//		sourcify.WithHTTPClient(httpClient),
//		sourcify.WithBaseURL("https://sourcify.dev/server"),
//		sourcify.WithRetryOptions(
//			sourcify.WithMaxRetries(3),
//			sourcify.WithDelay(2*time.Second),
//		),
//
// )
//
//	metadata, err := sourcify.GetContractMetadata(client, 1, common.HexToAddress("0x1234567890abcdef"), sourcify.MethodMatchTypeFull)
//	if err != nil {
//	  log.Fatal(err)
//	}
//	fmt.Println("Contract Metadata:", metadata)
//
// For more information on the Sourcify API and its endpoints, refer to the official documentation:
//   - API Documentation: https://pkg.go.dev/github.com/txpull/sourcify-go
//   - GitHub Repository: https://github.com/txpull/sourcify-go
package sourcify
