package examples

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/txpull/sourcify-go"
)

func Example_GetSourceCode() {
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
	sources, err := sourcify.GetContractSourceCode(client, 56, common.HexToAddress("0x054B2223509D430269a31De4AE2f335890be5C8F"), sourcify.MethodMatchTypeAny)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Status: %+v\n", sources.Status)

	for _, source := range sources.Code {
		fmt.Printf("Name: %+v\n", source.Name)
		fmt.Printf("Path: %+v\n", source.Path)
		// Uncomment the following line to print the content of the source code
		//fmt.Printf("Content: %+v\n\n", source.Content)
	}
}
