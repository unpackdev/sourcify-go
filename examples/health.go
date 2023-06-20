package examples

import (
	"fmt"
	"net/http"
	"time"

	"github.com/txpull/sourcify-go"
)

// Example_GetHealth demonstrates how to check the health status of the Sourcify server using the Sourcify client.
func Example_GetHealth() {
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

	// Call the API method to get the health status
	status, err := sourcify.GetHealth(client)
	if err != nil {
		panic(err)
	}

	// Print the server's health status
	fmt.Printf("Is server alive and ready to receive requests: %v\n", status)
}
