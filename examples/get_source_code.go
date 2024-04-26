package examples

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/unpackdev/sourcify-go"
	"net/http"
	"time"
)

// Example_GetSourceCode demonstrates how to retrieve source code for a contract using the Sourcify client.
func Example_GetSourceCode() {
	// Create a custom HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create a new Sourcify client with custom options
	client := sourcify.NewClient(
		sourcify.WithHTTPClient(httpClient),
		sourcify.WithBaseURL("https://sourcify.dev/server"),
		sourcify.WithRetryOptions(
			sourcify.WithMaxRetries(3),
			sourcify.WithDelay(2*time.Second),
		),
	)

	// Get source code for the Binance Smart Chain with the address of the R3T contract
	sources, err := sourcify.GetContractSourceCode(client, 56, common.HexToAddress("0x054B2223509D430269a31De4AE2f335890be5C8F"), sourcify.MethodMatchTypeAny)
	if err != nil {
		panic(err)
	}

	// Print the status of the source code retrieval
	fmt.Printf("Status: %+v\n", sources.Status)

	// Iterate over the source code files and print their names and paths
	for _, source := range sources.Code {
		fmt.Printf("Name: %+v\n", source.Name)
		fmt.Printf("Path: %+v\n", source.Path)
		// Uncomment the following line to print the content of the source code
		//fmt.Printf("Content: %+v\n\n", source.Content)
	}
}
