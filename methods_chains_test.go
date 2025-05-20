package sourcify

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testChainsResponse = []byte(`[{"name":"Ethereum Mainnet","chain":"ETH","icon":"ethereum","features":[{"name":"EIP155"},{"name":"EIP1559"}],"faucets":[],"nativeCurrency":{"name":"Ether","symbol":"ETH","decimals":18},"infoURL":"https://ethereum.org","shortName":"eth","chainId":1,"networkId":1,"slip44":60,"ens":{"registry":"0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"},"explorers":[{"name":"etherscan","url":"https://etherscan.io","standard":"EIP3091"}],"supported":true,"monitored":true,"contractFetchAddress":"https://api.etherscan.io/api?module=contract&action=getcontractcreation&contractaddresses=${ADDRESS}&apikey=4K5JKHMHPYK79IQYH8C5E5F5FAVYQS4W5E","rpc":["http://10.10.42.102:8541","https://eth-mainnet.g.alchemy.com/v2/{ALCHEMY_ID}"],"etherscanAPI":"https://api.etherscan.io"}]`)
var testChainsResponseError = []byte(`{"error":"Internal Server Error"}`)

func TestGetChains(t *testing.T) {
	// Create a test server with a mocked handler
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testChainsResponse)
	}))
	defer server.Close()

	// Create a test client with the test server's URL
	client := NewClient(WithBaseURL(server.URL))

	// Call the GetChains function
	chains, err := GetChains(client)
	assert.NoError(t, err)

	// Validate the result
	expectedChains := []Chain{
		{
			Name:  "Ethereum Mainnet",
			Chain: "ETH",
			Icon:  "ethereum",
			Features: []ChainFeature{
				{Name: "EIP155"},
				{Name: "EIP1559"},
			},
			Faucets: []interface{}{},
			NativeCurrency: ChainNativeCurrency{
				Name:     "Ether",
				Symbol:   "ETH",
				Decimals: 18,
			},
			InfoURL:   "https://ethereum.org",
			ShortName: "eth",
			ChainID:   1,
			NetworkID: 1,
			Slip44:    60,
			Ens:       ChainEns{Registry: "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e"},
			Explorers: []ChainExplorer{
				{
					Name:     "etherscan",
					URL:      "https://etherscan.io",
					Standard: "EIP3091",
				},
			},
			Supported:            true,
			Monitored:            true,
			ContractFetchAddress: "https://api.etherscan.io/api?module=contract&action=getcontractcreation&contractaddresses=${ADDRESS}&apikey=4K5JKHMHPYK79IQYH8C5E5F5FAVYQS4W5E",
			RPC: []string{
				"http://10.10.42.102:8541",
				"https://eth-mainnet.g.alchemy.com/v2/{ALCHEMY_ID}",
			},
			EtherscanAPI: "https://api.etherscan.io",
		},
	}

	assert.Equal(t, expectedChains, chains)
}

func TestGetChains_ErrorResponse(t *testing.T) {
	// Create a test server with a mocked error handler
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(testChainsResponseError)
	}))
	defer server.Close()

	// Create a test client with the test server's URL
	client := NewClient(WithBaseURL(server.URL))

	// Call the GetChains function
	_, err := GetChains(client)
	assert.Error(t, err)
}
