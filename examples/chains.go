package examples

import (
	"fmt"
	"net/http"
	"time"
)

// Example_GetChains demonstrates how to retrieve chains using the Sourcify client.
func Example_GetChains() {
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

	// Call get chains to retrieve all chains from Sourcify
	chains, err := sourcify.GetChains(client)
	if err != nil {
		panic(err)
	}

	for _, chain := range chains {
		fmt.Printf("Chain: %+v\n\n", chain)
	}
}
