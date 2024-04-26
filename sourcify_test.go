package sourcify

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SourcifySuite is a struct that embeds suite.Suite and contains fields necessary for
// executing tests. This includes a Client to interact with the sourcify service, a list of
// addresses to check, a list of chain IDs to use in testing, and specific instances
// of an address and a chain ID.
type SourcifySuite struct {
	suite.Suite
	client          *Client
	Addresses       []string
	ChainIDs        []int
	SpecificChainID int
	SpecificAddress common.Address
}

// SetupTest initializes the test suite before each test. It creates a new Client with
// specified options and initializes the addresses and chain IDs to be used in the tests.
func (suite *SourcifySuite) SetupTest() {
	suite.client = NewClient(
		WithBaseURL("https://sourcify.dev/server"),
		WithRetryOptions(
			WithMaxRetries(3),
			WithDelay(2*time.Second),
		),
		//WithRateLimit(20, 1*time.Second),
	)
	suite.Addresses = []string{"0x054B2223509D430269a31De4AE2f335890be5C8F"}
	suite.ChainIDs = []int{1, 56, 137}
	suite.SpecificChainID = 56                                                                // Add the specific chain ID for the two methods
	suite.SpecificAddress = common.HexToAddress("0x054B2223509D430269a31De4AE2f335890be5C8F") // Add the specific address for the two methods
}

// TestGetContractMetadata tests the GetContractMetadata function. It asserts that no error
// is returned and the metadata is not nil.
func (suite *SourcifySuite) TestGetContractMetadata() {
	// Act
	metadata, err := GetContractMetadata(suite.client, suite.SpecificChainID, suite.SpecificAddress, MethodMatchTypeFull)

	// Assert
	assert := assert.New(suite.T())
	assert.NoError(err, "GetContractMetadata should not return an error")
	assert.NotNil(metadata, "metadata should not be nil")
}

// TestGetContractSourceCode tests the GetContractSourceCode function. It asserts that no error
// is returned, the source code is not nil, and the length of the source code is 2.
func (suite *SourcifySuite) TestGetContractSourceCode() {
	sourceCode, err := GetContractSourceCode(suite.client, suite.SpecificChainID, suite.SpecificAddress, MethodMatchTypeFull)

	assert := assert.New(suite.T())
	assert.NoError(err, "Expected GetContractSourceCode to run without error")
	assert.NotNil(sourceCode, "source code should not be nil")
	assert.Equal(len(sourceCode.Code), 3, "Expected source code to have 3 files")
}

// TestGetChains tests the GetChains function. It asserts that no error is returned, the chains
// are not nil, and the length of the chains is at least 98.
func (suite *SourcifySuite) TestGetChains() {
	chains, err := GetChains(suite.client)

	assert := assert.New(suite.T())
	assert.NoError(err, "Expected GetChains to run without error")
	assert.NotNil(chains, "source code should not be nil")
	assert.GreaterOrEqual(len(chains), 98, "Expected source code to have at least 98 chains")
}

// TestGetAvailableContractAddresses tests the GetAvailableContractAddresses function. It asserts
// that no error is returned, the addresses are not nil, and the length of the full and partial addresses
// are at least 1000 each.
func (suite *SourcifySuite) TestGetAvailableContractAddresses() {
	// Act
	addresses, err := GetAvailableContractAddresses(suite.client, suite.SpecificChainID)

	// Assert
	assert := assert.New(suite.T())
	assert.NoError(err, "GetAvailableContractAddresses should not return an error")
	assert.NotNil(addresses, "addresses should not be nil")
	assert.GreaterOrEqual(len(addresses.Full), 1000, "Expected source code to have at least 1000 addresses")
	assert.GreaterOrEqual(len(addresses.Partial), 1000, "Expected source code to have at least 1000 addresses")
}

// TestCheckContractByAddresses tests the CheckContractByAddresses function. It asserts that no error
// is returned and the checks are not nil.
func (suite *SourcifySuite) TestCheckContractByAddresses() {
	// Act
	checks, err := CheckContractByAddresses(suite.client, suite.Addresses, suite.ChainIDs, MethodMatchTypeFull)

	// Assert
	assert := assert.New(suite.T())
	assert.NoError(err, "CheckContractByAddresses should not return an error")
	assert.NotNil(checks, "checks should not be nil")
}

// TestGetHealth tests the GetHealth function. It asserts that no error is returned and the status
// returned is true.
func (suite *SourcifySuite) TestGetHealth() {
	// Act
	status, err := GetHealth(suite.client)

	// Assert
	assert := assert.New(suite.T())
	assert.NoError(err, "GetHealth should not return an error")
	assert.True(status, "status should be true")
}

// TestGetContractFiles tests the GetContractFiles function. It asserts that no error is returned,
// the tree is not nil, and the length of the files is 2.
func (suite *SourcifySuite) TestGetContractFiles() {
	tree, err := GetContractFiles(suite.client, suite.SpecificChainID, suite.SpecificAddress, MethodMatchTypeFull)

	assert := assert.New(suite.T())
	assert.NoError(err, "Expected GetContractSourceCode to run without error")
	assert.NotNil(tree, "tree code should not be nil")
	assert.Equal(len(tree.Files), 3, "Expected tree to have 3 files")
}

// TestGetContractFiles tests the GetContractFiles function. It asserts that no error is returned,
// the tree is not nil, and the length of the files is 2.
func TestSourcifySuite(t *testing.T) {
	suite.Run(t, new(SourcifySuite))
}
