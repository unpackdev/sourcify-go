package examples

import (
	"fmt"
	"net/http"
	"time"

	"github.com/txpull/sourcify-go"
)

func GetChains() {
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
	chains, err := sourcify.GetChains(client)
	if err != nil {
		panic(err)
	}

	// Process the response
	for _, chain := range chains {
		fmt.Printf("Chain: %+v\n\n", chain)
	}
}
