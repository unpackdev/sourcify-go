package examples

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/txpull/sourcify-go"
)

// Example_GetFiles demonstrates how to retrieve source files for a contract using the Sourcify client.
func Example_GetFiles() {
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

	// Get source files for the Binance Smart Chain with the address of the R3T contract
	files, err := sourcify.GetContractFiles(client, 56, common.HexToAddress("0x054B2223509D430269a31De4AE2f335890be5C8F"), sourcify.MethodMatchTypeAny)
	if err != nil {
		panic(err)
	}

	// Print the status of the response
	fmt.Printf("Status: %+v\n", files.Status)

	// Print the paths of the retrieved files
	for _, file := range files.Files {
		fmt.Printf("Path: %+v\n", file)
	}
}
