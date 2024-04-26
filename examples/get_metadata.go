package examples

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/unpackdev/sourcify-go"
	"net/http"
	"time"
)

// Example_GetMetadata demonstrates how to retrieve full metadata for a contract using the Sourcify client.
func Example_GetMetadata() {
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

	// Get full metadata for the Binance Smart Chain with the address of the R3T contract
	fullMetadata, err := sourcify.GetContractMetadata(client, 56, common.HexToAddress("0x054B2223509D430269a31De4AE2f335890be5C8F"), sourcify.MethodMatchTypeFull)
	if err != nil {
		panic(err)
	}

	// Print the full match metadata
	fmt.Printf("Full Match Metadata: %+v\n", fullMetadata)
}
