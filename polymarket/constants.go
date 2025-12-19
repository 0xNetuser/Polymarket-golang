package polymarket

// Access levels
const (
	L0 = 0 // No authentication
	L1 = 1 // Private key authentication
	L2 = 2 // API credentials authentication
)

const (
	CredentialCreationWarning = `🚨🚨🚨
Your credentials CANNOT be recovered after they've been created.
Be sure to store them safely!
🚨🚨🚨`

	L1AuthUnavailable = "A private key is needed to interact with this endpoint!"
	L2AuthUnavailable = "API Credentials are needed to interact with this endpoint!"
	BuilderAuthUnavailable = "Builder API Credentials needed to interact with this endpoint!"

	ZeroAddress = "0x0000000000000000000000000000000000000000"

	// Chain IDs
	Amoy   = 80002
	Polygon = 137

	EndCursor = "LTE="
)

// Order sides
const (
	BUY  = "BUY"
	SELL = "SELL"
)

// Order types
type OrderType string

const (
	OrderTypeGTC OrderType = "GTC" // Good Till Cancel
	OrderTypeFOK OrderType = "FOK" // Fill Or Kill
	OrderTypeGTD OrderType = "GTD" // Good Till Date
	OrderTypeFAK OrderType = "FAK" // Fill And Kill
)

// Tick sizes
type TickSize string

const (
	TickSize01   TickSize = "0.1"
	TickSize001  TickSize = "0.01"
	TickSize0001 TickSize = "0.001"
	TickSize00001 TickSize = "0.0001"
)

// Signature types
const (
	SignatureTypeEOA = 0 // Externally Owned Account
	SignatureTypeEmail = 1 // Email/Magic wallet
	SignatureTypeBrowser = 2 // Browser wallet proxy
)

// Header names
const (
	PolyAddress   = "POLY_ADDRESS"
	PolySignature = "POLY_SIGNATURE"
	PolyTimestamp = "POLY_TIMESTAMP"
	PolyNonce     = "POLY_NONCE"
	PolyAPIKey    = "POLY_API_KEY"
	PolyPassphrase = "POLY_PASSPHRASE"
)

// CLOB domain constants
const (
	CLOBDomainName = "ClobAuthDomain"
	CLOBVersion    = "1"
	MsgToSign      = "This message attests that I control the given wallet"
)

