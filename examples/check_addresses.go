package examples

import (
	"fmt"
	"net/http"
	"time"

	"github.com/txpull/sourcify-go"
)

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

	contractAddresses := []string{
		"0x054B2223509D430269a31De4AE2f335890be5C8F",
	}

	chainIds := []int{
		56,  // Binance Smart Chain
		1,   // Ethereum Mainnet
		137, // Polygon
	}

	checks, err := sourcify.CheckContractByAddresses(client, contractAddresses, chainIds, sourcify.MethodMatchTypeAny)
	if err != nil {
		panic(err)
	}

	for _, check := range checks {
		fmt.Printf("Contract Addresses Check Response: %+v\n", check)
	}
}
