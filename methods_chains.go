package sourcify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MethodGetChains represents the API endpoint for getting the chains (networks) added to Sourcify in the Sourcify service.
// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
// Returns the chains (networks) added to Sourcify. Contains both supported, unsupported, monitored, and unmonitored chains.
// More information: https://docs.sourcify.dev/docs/api/chains/
var MethodGetChains = Method{
	Name:           "Retrieve staticly served files over the server for partial match contract",
	URI:            "/chains",
	MoreInfo:       "https://docs.sourcify.dev/docs/api/chains/",
	Method:         "GET",
	ParamType:      MethodParamTypeUri,
	RequiredParams: []string{},
	Params:         []MethodParam{},
}

// ChainFeature represents a feature of a chain.
type ChainFeature struct {
	Name string `json:"name"`
}

// ChainNativeCurrency represents the native currency of a chain.
type ChainNativeCurrency struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

// ChainEns represents the ENS (Ethereum Name Service) configuration of a chain.
type ChainEns struct {
	Registry string `json:"registry"`
}

// ChainExplorer represents an explorer configuration for a chain.
type ChainExplorer struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Standard string `json:"standard"`
}

// ChainParent represents the parent chain information for a chain.
type ChainParent struct {
	Type    string `json:"type"`
	Chain   string `json:"chain"`
	Bridges []struct {
		URL string `json:"url"`
	} `json:"bridges"`
}

// Chain represents a chain (network) added to Sourcify.
type Chain struct {
	Name                 string              `json:"name"`
	Chain                string              `json:"chain"`
	Icon                 string              `json:"icon,omitempty"`
	Features             []ChainFeature      `json:"features,omitempty"`
	Faucets              []any               `json:"faucets"`
	NativeCurrency       ChainNativeCurrency `json:"nativeCurrency"`
	InfoURL              string              `json:"infoURL"`
	ShortName            string              `json:"shortName"`
	ChainID              int                 `json:"chainId"`
	NetworkID            int                 `json:"networkId"`
	Slip44               int                 `json:"slip44,omitempty"`
	Ens                  ChainEns            `json:"ens,omitempty"`
	Explorers            []ChainExplorer     `json:"explorers,omitempty"`
	Supported            bool                `json:"supported"`
	Monitored            bool                `json:"monitored"`
	ContractFetchAddress string              `json:"contractFetchAddress,omitempty"`
	RPC                  []string            `json:"rpc"`
	EtherscanAPI         string              `json:"etherscanAPI,omitempty"`
	Title                string              `json:"title,omitempty"`
	TxRegex              string              `json:"txRegex,omitempty"`
	RedFlags             []string            `json:"redFlags,omitempty"`
	Status               string              `json:"status,omitempty"`
	Parent               ChainParent         `json:"parent,omitempty"`
	GraphQLFetchAddress  string              `json:"graphQLFetchAddress,omitempty"`
}

// GetChains gets the chains (networks) added to Sourcify by calling the MethodGetChains endpoint using the provided client.
// It returns the chains and an error if any occurred during the request.
func GetChains(client *Client) ([]Chain, error) {
	response, statusCode, err := client.CallMethod(MethodGetChains)
	if err != nil {
		return nil, err
	}

	// Close the io.ReadCloser interface
	// This is important as CallMethod is NOT closing the response body!
	// You'll have memory leaks if you don't do this!
	defer response.Close()

	if statusCode != http.StatusOK {
		if rErr := ToErrorResponse(response); rErr != nil {
			return nil, rErr
		}

		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	var chains []Chain
	if err := json.NewDecoder(response).Decode(&chains); err != nil {
		return nil, err
	}

	return chains, nil
}
