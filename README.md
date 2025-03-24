# TONSub

**tonSub** is a Go package aimed at interacting with transactions within the TON network. The package supports Jetton, NFT, and TON transactions, providing transaction details via a convenient event-driven interface.

## Installation

To install the `tonSub` package with Go Modules, execute the following command:

```bash
go get github.com/elum-utils/tonsub
```

## Initialization

Initialize the subscription by connecting to the TON network:

```go
subs, err := tonsub.New(
    "wallet_address",
    "https://ton.org/global.config.json",
)

if err != nil {
    panic(err.Error())
}
```

## Transaction Handling

### Handling TON Transactions

Register a handler for standard TON transactions:

```go
subs.OnTON(func(t *tonsub.RootTON) {
    jsonData, err := json.MarshalIndent(t, "", "  ")
    if err != nil {
        fmt.Println("Error encoding JSON:", err)
        return
    }

    fmt.Println("TON Transaction")
    fmt.Println(string(jsonData))
})
```

**Data Structures:**

```go
// RootTON defines the structure for TON transactions, encapsulating core transaction data.
type RootTON struct {
	OpCode      uint64 `json:"op_code"`      // Operation code indicating the type of transaction
	IHRDisabled bool   `json:"ihr_disabled"` // Flag indicating if Instant Hypercube Routing is disabled
	Bounce      bool   `json:"bounce"`       // Flag indicating if bounces are enabled for this transaction
	Bounced     bool   `json:"bounced"`      // Flag indicating if this message has already bounced
	SndrAddr    string `json:"snr_addr"`     // Sender address of the transaction as a string
	SrcAddr     string `json:"src_addr"`     // Source address of the transaction, represented as a string
	DstAddr     string `json:"dst_addr"`     // Destination address of the transaction, represented as a string
	Amount      string `json:"amount"`       // Amount involved in the transaction, formatted as a string
	IHRFee      string `json:"ihr_fee"`      // Fee for Instant Hypercube Routing, formatted as a string
	FwdFee      string `json:"fwd_fee"`      // Forward fee in string format
	CreatedLT   uint64 `json:"created_lt"`   // Logical time for when the transaction was created
	CreatedAt   uint32 `json:"created_at"`   // Unix timestamp for the creation of the transaction
	Body        Ton    `json:"body"`         // Nested struct for message-related data
}

// Ton holds additional transaction message details.
type Ton struct {
    Message string `json:"message"` // Extracted message from the transaction payload
    TxHash  string `json:"tx_hash"` // Transaction hash
}
```

### Handling Jetton Transactions

To process Jetton transactions, use:

```go
subs.OnJetton(func(t *tonsub.RootJetton) {
    jsonData, err := json.MarshalIndent(t, "", "  ")
    if err != nil {
        fmt.Println("Error encoding JSON:", err)
        return
    }

    fmt.Println("Jetton Transaction")
    fmt.Println(string(jsonData))
})
```

**Data Structures:**

```go
// RootJetton defines the main structure for Jetton transactions.
type RootJetton struct {
	IHRDisabled bool   `json:"ihr_disabled"` // Indicator if Instant Hypercube Routing is disabled
	Bounce      bool   `json:"bounce"`       // Indicator if the message should bounce back if it fails
	Bounced     bool   `json:"bounced"`      // Indicator if the message has already bounced
	SndrAddr    string `json:"snr_addr"`     // Sender address of the transaction as a string
	SrcAddr     string `json:"src_addr"`     // Source address as a string
	DstAddr     string `json:"dst_addr"`     // Destination address as a string
	Amount      string `json:"amount"`       // Amount involved in the transaction as a string
	IHRFee      string `json:"ihr_fee"`      // Instant Hypercube Routing fee as a string
	FwdFee      string `json:"fwd_fee"`      // Forwarding fee as a string
	CreatedLT   uint64 `json:"created_lt"`   // Logical time when the message was created
	CreatedAt   uint32 `json:"created_at"`   // Unix timestamp when the message was created
	Body        Jetton `json:"body"`         // Contained Jetton information
}

// Jetton holds information about Jetton transaction bodies.
type Jetton struct {
    OpCode  uint64 `json:"op_code"`  // Operation code for the transaction
    QueryID uint64 `json:"query_id"` // Unique identifier for the query
    Amount  string `json:"amount"`   // Jetton amount involved
    Sender  string `json:"sender"`   // Sender's address
    Message string `json:"message"`  // Optional message
    TxHash  string `json:"tx_hash"`  // Transaction hash
}
```

### Handling NFT Transactions

For NFT transactions, define the following handler:

```go
subs.OnNFT(func(t *tonsub.RootNFT) {
    jsonData, err := json.MarshalIndent(t, "", "  ")
    if err != nil {
        fmt.Println("Error encoding JSON:", err)
        return
    }

    fmt.Println("NFT Transaction")
    fmt.Println(string(jsonData))
})
```

**Data Structures:**

```go
// RootNFT represents the main structure for NFT transactions.
type RootNFT struct {
	IHRDisabled bool   `json:"ihr_disabled"` // Flag indicating if Instant Hypercube Routing is disabled
	Bounce      bool   `json:"bounce"`       // Flag for bounce status of the transaction
	Bounced     bool   `json:"bounced"`      // Flag indicating if the transaction has already bounced
	SndrAddr    string `json:"snr_addr"`     // Sender address of the transaction as a string
	SrcAddr     string `json:"src_addr"`     // Source address of the transaction as a string
	DstAddr     string `json:"dst_addr"`     // Destination address of the transaction as a string
	Amount      string `json:"amount"`       // Amount involved in the transaction represented as a string
	IHRFee      string `json:"ihr_fee"`      // Instant Hypercube Routing fee as a string
	FwdFee      string `json:"fwd_fee"`      // Forward fee as a string
	CreatedLT   uint64 `json:"created_lt"`   // Logical time indicating when the transaction was created
	CreatedAt   uint32 `json:"created_at"`   // Unix timestamp indicating when the transaction was created
	Body        NFT    `json:"body"`         // Nested NFT structure containing detailed NFT data
}

// NFT captures details specific to the NFT data.
type NFT struct {
    OpCode            uint64 `json:"op_code"`            // Operation code for handling
    Initialized       bool   `json:"initialized"`        // Initialization status
    Index             string `json:"index"`              // NFT's collection index
    Address           string `json:"nft_address"`        // Address of the NFT entity
    OwnerAddress      string `json:"owner_address"`      // Current owner's address
    CollectionAddress string `json:"collection_address"` // Address of NFT collection
    Meta              string `json:"meta"`               // Metadata associated with the NFT
    Message           string `json:"message"`            // Optional transaction message
    TxHash            string `json:"tx_hash"`            // Transaction hash value
}
```

## Support the Project

If this project aids your development and you would like to support its growth, please consider making a donation:

**TON Wallet:** 
```
UQActNkydex6WaHNO55qMHkzNjLKsKN_oF_dKMpZ3K7lAAAA
```

Your contributions help in enhancing and progressing the project. Thank you for your support!