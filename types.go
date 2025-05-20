package sourcify

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
	Outputs         []OutputDetail `json:"outputs"`
	StateMutability string         `json:"stateMutability"`
	Type            string         `json:"type"`
	Constant        bool           `json:"constant"`
	Payable         bool           `json:"payable"`
	Anonymous       bool           `json:"anonymous,omitempty"`
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

// CborAuxData represents CBOR auxiliary data
type CborAuxData struct {
	Offset int64  `json:"offset"`
	Value  string `json:"value"`
}

// Transformation represents a bytecode transformation operation
type Transformation struct {
	ID     string `json:"id,omitempty"`
	Type   string `json:"type"`
	Offset int64  `json:"offset"`
	Reason string `json:"reason"`
}

// TransformationValues represents the values used in bytecode transformations
type TransformationValues struct {
	CborAuxdata          map[string]string `json:"cborAuxdata,omitempty"`
	ConstructorArguments string            `json:"constructorArguments,omitempty"`
}

// Libraries defines contract libraries mapping
type Libraries map[string]string

// Bytecode represents a contract's bytecode information
type Bytecode struct {
	CborAuxdata          map[string]CborAuxData `json:"cborAuxdata"`
	LinkReferences       map[string]interface{} `json:"linkReferences"`
	OnchainBytecode      string                 `json:"onchainBytecode"`
	RecompiledBytecode   string                 `json:"recompiledBytecode"`
	SourceMap            string                 `json:"sourceMap"`
	TransformationValues TransformationValues   `json:"transformationValues"`
	Transformations      []Transformation       `json:"transformations"`
	ImmutableReferences  map[string]interface{} `json:"immutableReferences,omitempty"`
}

// UintType represents a uint type definition
type UintType struct {
	Encoding      string `json:"encoding"`
	Label         string `json:"label"`
	NumberOfBytes string `json:"numberOfBytes"`
}

// StorageLayout represents the storage layout of a contract
type StorageLayout struct {
	Storage []StorageEntry         `json:"storage"`
	Types   map[string]StorageType `json:"types"`
}

// StorageType represents a type definition in the storage layout
type StorageType struct {
	Label         string `json:"label"`
	Encoding      string `json:"encoding"`
	NumberOfBytes string `json:"numberOfBytes"`
}

// StorageEntry represents a storage variable in the contract
type StorageEntry struct {
	Slot     string `json:"slot"`
	Type     string `json:"type"`
	AstId    int    `json:"astId"`
	Label    string `json:"label"`
	Offset   int    `json:"offset"`
	Contract string `json:"contract"`
}

// ContractReference represents a reference to a contract file
type ContractReference struct {
	ID int64 `json:"id"`
}

// ContentReference represents a reference to a contract's content
type ContentReference struct {
	Content string `json:"content"`
}

// SourceURLs represents source URLs information
type SourceURLs struct {
	Keccak256 string   `json:"keccak256"`
	License   string   `json:"license"`
	Urls      []string `json:"urls"`
}

// DevDoc represents developer documentation
type DevDoc struct {
	Details string         `json:"details"`
	Kind    string         `json:"kind"`
	Methods map[string]any `json:"methods"`
	Title   string         `json:"title"`
	Version int64          `json:"version"`
}

// UserDoc represents user documentation
type UserDoc struct {
	Kind    string         `json:"kind"`
	Methods map[string]any `json:"methods"`
	Version int64          `json:"version"`
}

// EVMBytecode represents the EVM bytecode information
type EVMBytecode struct {
	LinkReferences map[string]interface{} `json:"linkReferences"`
	Object         string                 `json:"object"`
	SourceMap      string                 `json:"sourceMap"`
}

// EVMDeployedBytecode represents the EVM deployed bytecode information
type EVMDeployedBytecode struct {
	ImmutableReferences map[string]interface{} `json:"immutableReferences"`
	LinkReferences      map[string]interface{} `json:"linkReferences"`
	Object              string                 `json:"object"`
	SourceMap           string                 `json:"sourceMap"`
}

// EVM represents EVM-related information
type EVM struct {
	Bytecode         EVMBytecode         `json:"bytecode"`
	DeployedBytecode EVMDeployedBytecode `json:"deployedBytecode"`
}

// MetadataDetail provides additional metadata.
type MetadataDetail struct {
	BytecodeHash string `json:"bytecodeHash"` // Hash of the bytecode
}

// Settings includes details about the compiler settings used.
type Settings struct {
	CompilationTarget CompilationTarget `json:"compilationTarget"` // CompilationTarget represents the compilation target details.
	EvmVersion        string            `json:"evmVersion"`        // EVM version used
	Libraries         Libraries         `json:"libraries"`         // Libraries used in the source code
	Metadata          MetadataDetail    `json:"metadata"`          // MetadataDetail represents additional metadata.
	Optimizer         Optimizer         `json:"optimizer"`         // Optimizer represents the compiler optimization details.
	Remappings        []string          `json:"remappings"`        // Remappings used in the source code
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

// Compiler contains information about the compiler used to compile the smart contract.
type Compiler struct {
	Version string `json:"version"` // Version of the compiler
}

// Output represents details of the compiled code in the metadata.
type Output struct {
	Abi     []ABIEntry `json:"abi,omitempty"`
	Devdoc  DevDoc     `json:"devdoc,omitempty"`
	Userdoc UserDoc    `json:"userdoc,omitempty"`
}

// MetadataSource represents a source file in the metadata.
type MetadataSource struct {
	Keccak256 string   `json:"keccak256,omitempty"`
	Urls      []string `json:"urls,omitempty"`
	License   string   `json:"license,omitempty"`
}

// Metadata represents the top-level structure for compiler metadata
// for Ethereum smart contracts.
type Metadata struct {
	Compiler Compiler                  `json:"compiler"` // Compiler contains information about the compiler used.
	Language string                    `json:"language"` // Language of the source code
	Output   Output                    `json:"output"`   // Output represents details of the compiled code.
	Settings Settings                  `json:"settings"` // Settings represent the compiler settings used.
	Sources  map[string]MetadataSource `json:"sources"`  // Sources represents the details of the source code.
	Version  int                       `json:"version"`  // Version of the metadata.
}

// Compilation contains information about the contract compilation process
type Compilation struct {
	Compiler           string     `json:"compiler"`
	CompilerSettings   EVMVersion `json:"compilerSettings"`
	CompilerVersion    string     `json:"compilerVersion"`
	FullyQualifiedName string     `json:"fullyQualifiedName"`
	Language           string     `json:"language"`
	Name               string     `json:"name"`
}

// ProxyResolution contains information about proxy contract resolution
type ProxyResolution struct {
	Implementations []string `json:"implementations"`
	IsProxy         bool     `json:"isProxy"`
	ProxyType       string   `json:"proxyType,omitempty"`
}

// StdJSONInput represents the standard JSON input format for the compiler
type StdJSONInput struct {
	Language string                      `json:"language"`
	Settings Settings                    `json:"settings"`
	Sources  map[string]ContentReference `json:"sources"`
}

// SourceIdReference represents the ID reference for a source file
type SourceIdReference struct {
	ID int `json:"id"`
}

// SourceIds maps file paths to their corresponding source IDs
type SourceIds map[string]SourceIdReference

// StdJSONOutput represents the standard JSON output format from the compiler
type StdJSONOutput struct {
	Contracts map[string]map[string]ContractOutput `json:"contracts"`
	Sources   SourceIds                            `json:"sources"`
}

// ContractOutput contains the compiled information of a contract
type ContractOutput struct {
	Abi           []ABIEntry    `json:"abi"`
	Devdoc        DevDoc        `json:"devdoc"`
	Evm           EVM           `json:"evm"`
	Metadata      string        `json:"metadata"`
	StorageLayout StorageLayout `json:"storageLayout"`
	Userdoc       UserDoc       `json:"userdoc"`
}

// SourceContent represents the content of a source file
type SourceContent struct {
	Content string `json:"content"`
}

// Sources maps file names to their content
type Sources map[string]SourceContent

// OutputDetail holds information about the output parameters of the functions.
type OutputDetail struct {
	InternalType string `json:"internalType"` // Internal type of the parameter
	Name         string `json:"name"`         // Name of the parameter
	Type         string `json:"type"`         // Type of the parameter
}
