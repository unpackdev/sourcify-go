package examples

import (
	"fmt"
	"net/http"
	"time"

	"github.com/txpull/sourcify-go"
)

// Example_CheckAddresses demonstrates how to check contract addresses using the Sourcify client.
func Example_CheckAddresses() {
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

	// Define contract addresses to check
	contractAddresses := []string{
		"0x054B2223509D430269a31De4AE2f335890be5C8F",
	}

	// Define chain IDs to check
	chainIds := []int{
		56,  // Binance Smart Chain
		1,   // Ethereum Mainnet
		137, // Polygon
	}

	// Call the API method to check contract addresses
	checks, err := sourcify.CheckContractByAddresses(client, contractAddresses, chainIds, sourcify.MethodMatchTypeAny)
	if err != nil {
		panic(err)
	}

	// Process the response
	for _, check := range checks {
		fmt.Printf("Contract Addresses Check Response: %+v\n", check)
	}
}
