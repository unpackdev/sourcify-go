package examples

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/unpackdev/sourcify-go"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"time"
)

// Example_GetContractByChainIdAndAddress demonstrates how to retrieve all available information for a contract using the Sourcify client.
func Example_GetContractByChainIdAndAddress() {
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
	sources, err := sourcify.GetContractByChainIdAndAddress(client, 56, common.HexToAddress("0x054B2223509D430269a31De4AE2f335890be5C8F"), []string{"all"}, []string{})
	if err != nil {
		panic(err)
	}

	spew.Dump(sources)
}
