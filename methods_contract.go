package sourcify

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"net/http"
	"strings"
)

var (
	// MethodGetContractByChainIdAndAddress represents the API endpoint for getting the contract addresses with full or partial match in the Sourcify service.
	// It includes the name, the HTTP method, the URI, and the parameters necessary for the request.
	// Returns all verified sources from the repository for the desired contract address and chain, including metadata.json. Searches only for full matches.
	// More information: https://docs.sourcify.dev/docs/api/#/Contract%20Lookup
	MethodGetContractByChainIdAndAddress = Method{
		Name:           "Get verified contract addresses for the chain full or partial match",
		URI:            "/v2/contract/:chain/:address",
		MoreInfo:       "https://docs.sourcify.dev/docs/api/#/Contract%20Lookup",
		Method:         "GET",
		ParamType:      MethodParamTypeUri,
		RequiredParams: []string{":chain", ":address"},
		Params:         []MethodParam{},
	}
)

// ABIParameter represents a parameter in an ABI function or event
type ABIParameter struct {
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}

// ABIEntry represents a function, event, or error in an ABI
type ABIEntry struct {
	Inputs          []ABIParameter `json:"inputs"`
	Name            string         `json:"name"`
	Outputs         []ABIParameter `json:"outputs"`
	StateMutability string         `json:"stateMutability"`
	Type            string         `json:"type"`
}

// BytecodeHash represents the bytecode hash settings in metadata
type BytecodeHash struct {
	BytecodeHash string `json:"bytecodeHash"`
}

// Optimizer represents compiler optimizer settings
type Optimizer struct {
	Enabled bool  `json:"enabled"`
	Runs    int64 `json:"runs"`
}

// EVMVersion represents EVM version settings used in compiler settings
type EVMVersion struct {
	EvmVersion string        `json:"evmVersion"`
	Libraries  struct{}      `json:"libraries"`
	Metadata   BytecodeHash  `json:"metadata"`
	Optimizer  Optimizer     `json:"optimizer"`
	Remappings []interface{} `json:"remappings"`
}

// CborAuxdataOne represents the CBOR auxdata structure
type CborAuxdataOne struct {
	Offset int64  `json:"offset"`
	Value  string `json:"value"`
}

// CborAuxdata represents CBOR auxiliary data
type CborAuxdata struct {
	One CborAuxdataOne `json:"1"`
}

// Bytecode represents a contract's bytecode information
type Bytecode struct {
	CborAuxdata          CborAuxdata   `json:"cborAuxdata"`
	LinkReferences       struct{}      `json:"linkReferences"`
	OnchainBytecode      string        `json:"onchainBytecode"`
	RecompiledBytecode   string        `json:"recompiledBytecode"`
	SourceMap            string        `json:"sourceMap"`
	TransformationValues struct{}      `json:"transformationValues"`
	Transformations      []interface{} `json:"transformations"`
}

// RuntimeBytecode extends the Bytecode type with immutable references
type RuntimeBytecode struct {
	CborAuxdata          CborAuxdata   `json:"cborAuxdata"`
	ImmutableReferences  struct{}      `json:"immutableReferences"`
	LinkReferences       struct{}      `json:"linkReferences"`
	OnchainBytecode      string        `json:"onchainBytecode"`
	RecompiledBytecode   string        `json:"recompiledBytecode"`
	SourceMap            string        `json:"sourceMap"`
	TransformationValues struct{}      `json:"transformationValues"`
	Transformations      []interface{} `json:"transformations"`
}

// StorageEntry represents a single storage entry
type StorageEntry struct {
	AstID    int64  `json:"astId"`
	Contract string `json:"contract"`
	Label    string `json:"label"`
	Offset   int64  `json:"offset"`
	Slot     string `json:"slot"`
	Type     string `json:"type"`
}

// UintType represents a uint type definition
type UintType struct {
	Encoding      string `json:"encoding"`
	Label         string `json:"label"`
	NumberOfBytes string `json:"numberOfBytes"`
}

// StorageLayout represents the storage layout of a contract
type StorageLayout struct {
	Storage []StorageEntry `json:"storage"`
	Types   struct {
		TUint256 UintType `json:"t_uint256"`
	} `json:"types"`
}

// ContractReference represents a reference to a contract file
type ContractReference struct {
	ID int64 `json:"id"`
}

// ContentReference represents a reference to a contract's content
type ContentReference struct {
	Content string `json:"content"`
}

// FileReferences represents contract file references in sources
type FileReferences struct {
	Contracts_1Storage_sol ContractReference `json:"contracts/1_Storage.sol"`
}

// ContentReferences represents contract content references in sources
type ContentReferences struct {
	Contracts_1Storage_sol ContentReference `json:"contracts/1_Storage.sol"`
}

// SourceURLs represents source URLs information
type SourceURLs struct {
	Keccak256 string   `json:"keccak256"`
	License   string   `json:"license"`
	Urls      []string `json:"urls"`
}

// SourceURLReferences represents source URL references
type SourceURLReferences struct {
	Contracts_1Storage_sol SourceURLs `json:"contracts/1_Storage.sol"`
}

// RetrieveMethodReturn represents the return value documentation for the retrieve method
type RetrieveMethodReturn struct {
	Zero string `json:"_0"`
}

// RetrieveMethod represents the documentation for the retrieve method
type RetrieveMethod struct {
	Details string               `json:"details"`
	Returns RetrieveMethodReturn `json:"returns"`
}

// StoreMethodParams represents the parameter documentation for the store method
type StoreMethodParams struct {
	Num string `json:"num"`
}

// StoreMethod represents the documentation for the store method
type StoreMethod struct {
	Details string            `json:"details"`
	Params  StoreMethodParams `json:"params"`
}

// DocMethods represents method documentation
type DocMethods struct {
	Retrieve      RetrieveMethod `json:"retrieve()"`
	Store_uint256 StoreMethod    `json:"store(uint256)"`
}

// DevDoc represents developer documentation
type DevDoc struct {
	Details string     `json:"details"`
	Kind    string     `json:"kind"`
	Methods DocMethods `json:"methods"`
	Title   string     `json:"title"`
	Version int64      `json:"version"`
}

// UserDoc represents user documentation
type UserDoc struct {
	Kind    string   `json:"kind"`
	Methods struct{} `json:"methods"`
	Version int64    `json:"version"`
}

// EVMBytecode represents the EVM bytecode information
type EVMBytecode struct {
	LinkReferences struct{} `json:"linkReferences"`
	Object         string   `json:"object"`
	SourceMap      string   `json:"sourceMap"`
}

// EVMDeployedBytecode represents the EVM deployed bytecode information
type EVMDeployedBytecode struct {
	ImmutableReferences struct{} `json:"immutableReferences"`
	LinkReferences      struct{} `json:"linkReferences"`
	Object              string   `json:"object"`
	SourceMap           string   `json:"sourceMap"`
}

// EVM represents EVM-related information
type EVM struct {
	Bytecode         EVMBytecode         `json:"bytecode"`
	DeployedBytecode EVMDeployedBytecode `json:"deployedBytecode"`
}

// Settings includes details about the compiler settings used.
type Settings struct {
	CompilationTarget CompilationTarget `json:"compilationTarget"` // CompilationTarget represents the compilation target details.
	EvmVersion        string            `json:"evmVersion"`        // EVM version used
	Libraries         Libraries         `json:"libraries"`         // Libraries used in the source code
	Metadata          MetadataDetail    `json:"metadata"`          // MetadataDetail represents additional metadata.
	Optimizer         Optimizer         `json:"optimizer"`         // Optimizer represents the compiler optimization details.
	Remappings        []any             `json:"remappings"`        // Remappings used in the source code
}

// CompilationTarget holds the details of the compilation target.
type CompilationTarget map[string]string

// Deployment contains information about the contract deployment transaction
type Deployment struct {
	TransactionHash  string `json:"transactionHash"`  // The hash of the transaction that deployed the contract
	BlockNumber      string `json:"blockNumber"`      // The block number in which the contract was deployed
	TransactionIndex string `json:"transactionIndex"` // The index of the transaction in the block
	Deployer         string `json:"deployer"`         // The address that deployed the contract
}

// ContractResponse represents the response from the Sourcify API when retrieving contract information
type ContractResponse struct {
	Abi         []ABIEntry `json:"abi"`
	Address     string     `json:"address"`
	ChainID     string     `json:"chainId"`
	Compilation struct {
		Compiler           string     `json:"compiler"`
		CompilerSettings   EVMVersion `json:"compilerSettings"`
		CompilerVersion    string     `json:"compilerVersion"`
		FullyQualifiedName string     `json:"fullyQualifiedName"`
		Language           string     `json:"language"`
		Name               string     `json:"name"`
	} `json:"compilation"`
	CreationBytecode Bytecode   `json:"creationBytecode"`
	CreationMatch    string     `json:"creationMatch"`
	Deployment       Deployment `json:"deployment"`
	Devdoc           DevDoc     `json:"devdoc"`
	Match            string     `json:"match"`
	MatchID          string     `json:"matchId"`
	Metadata         struct {
		Compiler struct {
			Version string `json:"version"`
		} `json:"compiler"`
		Language string `json:"language"`
		Output   struct {
			Abi     []ABIEntry `json:"abi"`
			Devdoc  DevDoc     `json:"devdoc"`
			Userdoc UserDoc    `json:"userdoc"`
		} `json:"output"`
		Settings Settings            `json:"settings"`
		Sources  SourceURLReferences `json:"sources"`
		Version  int64               `json:"version"`
	} `json:"metadata"`
	ProxyResolution struct {
		Implementations []interface{} `json:"implementations"`
		IsProxy         bool          `json:"isProxy"`
		ProxyType       interface{}   `json:"proxyType"`
	} `json:"proxyResolution"`
	RuntimeBytecode RuntimeBytecode   `json:"runtimeBytecode"`
	RuntimeMatch    string            `json:"runtimeMatch"`
	SourceIds       FileReferences    `json:"sourceIds"`
	Sources         ContentReferences `json:"sources"`
	StdJSONInput    struct {
		Language string            `json:"language"`
		Settings Settings          `json:"settings"`
		Sources  ContentReferences `json:"sources"`
	} `json:"stdJsonInput"`
	StdJSONOutput struct {
		Contracts struct {
			Contracts_1Storage_sol struct {
				Storage struct {
					Abi           []ABIEntry    `json:"abi"`
					Devdoc        DevDoc        `json:"devdoc"`
					Evm           EVM           `json:"evm"`
					Metadata      string        `json:"metadata"`
					StorageLayout StorageLayout `json:"storageLayout"`
					Userdoc       UserDoc       `json:"userdoc"`
				} `json:"Storage"`
			} `json:"contracts/1_Storage.sol"`
		} `json:"contracts"`
		Sources FileReferences `json:"sources"`
	} `json:"stdJsonOutput"`
	StorageLayout StorageLayout `json:"storageLayout"`
	Userdoc       UserDoc       `json:"userdoc"`
	VerifiedAt    string        `json:"verifiedAt"`
}

// GetContractByChainIdAndAddress retrieves the available verified contract addresses for the given chain ID.
func GetContractByChainIdAndAddress(client *Client, chainId int, address common.Address, fields []string, omit []string) (*ContractResponse, error) {
	method := MethodGetContractByChainIdAndAddress

	if len(omit) > 0 && len(fields) > 0 {
		return nil, fmt.Errorf("fields and omit parameters cannot be used together")
	}

	// Omit and fields cannot co-exist together
	if len(omit) == 0 && len(fields) == 0 {
		fields = []string{"all"}
	}

	pFields := strings.Join(fields, ",")
	pOmit := strings.Join(omit, ",")

	method.SetParams(
		MethodParam{Key: ":chain", Value: chainId},
		MethodParam{Key: ":address", Value: address.Hex()},
		MethodParam{Key: "fields", Value: pFields},
		MethodParam{Key: "omit", Value: pOmit},
	)

	if err := method.Verify(); err != nil {
		return nil, err
	}

	response, statusCode, err := client.CallMethod(method)
	if err != nil {
		return nil, err
	}

	// Close the io.ReadCloser interface.
	// This is important as CallMethod is NOT closing the response body!
	// You'll have memory leaks if you don't do this!
	defer response.Close()

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	var toReturn ContractResponse
	if err := json.NewDecoder(response).Decode(&toReturn); err != nil {
		return nil, err
	}

	return &toReturn, nil
}
