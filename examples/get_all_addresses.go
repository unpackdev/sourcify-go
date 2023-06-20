package examples

import (
	"fmt"
	"net/http"
	"time"

	"github.com/txpull/sourcify-go"
)

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

	// Call the API method
	addresses, err := sourcify.GetAvailableContractAddresses(client, 56)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Full Match Addresses: %d\n", len(addresses.Full))
	fmt.Printf("Partial Match Addresses: %d\n\n", len(addresses.Partial))
}
