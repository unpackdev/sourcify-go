package examples

import (
	"fmt"
	"net/http"
	"time"

	"github.com/txpull/sourcify-go"
)

// Example_GetAllAddresses demonstrates how to retrieve all available contract addresses using the Sourcify client.
func Example_GetAllAddresses() {
	// Create a custom HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
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

	// Call the API method to retrieve all available contract addresses
	addresses, err := sourcify.GetAvailableContractAddresses(client, 56)
	if err != nil {
		panic(err)
	}

	// Print the number of full match addresses
	fmt.Printf("Full Match Addresses: %d\n", len(addresses.Full))

	// Print the number of partial match addresses
	fmt.Printf("Partial Match Addresses: %d\n\n", len(addresses.Partial))
}
