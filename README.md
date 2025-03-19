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
    OpCode      uint64 `json:"op_code"`      // Operation code for the transaction type
    IHRDisabled bool   `json:"ihr_disabled"` // Indicates if IHR is disabled
    Bounce      bool   `json:"bounce"`       // Determines if bounces are enabled
    Bounced     bool   `json:"bounced"`      // Indicates if the message has bounced back
    SrcAddr     string `json:"src_addr"`     // Source address as a string
    DstAddr     string `json:"dst_addr"`     // Destination address as a string
    Amount      string `json:"amount"`       // Transaction amount formatted as a string
    IHRFee      string `json:"ihr_fee"`      // IHR fee as a string
    FwdFee      string `json:"fwd_fee"`      // Forwarding fee in string format
    CreatedLT   uint64 `json:"created_lt"`   // Logical time of transaction creation
    CreatedAt   uint32 `json:"created_at"`   // UNIX timestamp of transaction creation
    Body        Ton    `json:"body"`         // Encapsulates message-related data
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
    IHRDisabled bool   `json:"ihr_disabled"` // Determines if IHR is disabled
    Bounce      bool   `json:"bounce"`       // Sets the bounce option for the message
    Bounced     bool   `json:"bounced"`      // Indicates if already bounced
    SrcAddr     string `json:"src_addr"`     // Source address in string format
    DstAddr     string `json:"dst_addr"`     // Destination address in string format
    Amount      string `json:"amount"`       // Amount involved in the transaction
    IHRFee      string `json:"ihr_fee"`      // Fee for Instant Hypercube Routing
    FwdFee      string `json:"fwd_fee"`      // Fee for forwarding the message
    CreatedLT   uint64 `json:"created_lt"`   // Logical time of message creation
    CreatedAt   uint32 `json:"created_at"`   // Timestamp when message was created
    Body        Jetton `json:"body"`         // Details specific to the Jetton transaction
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
    IHRDisabled bool   `json:"ihr_disabled"` // Indicates if IHR is disabled
    Bounce      bool   `json:"bounce"`       // Bounce flag for transaction messages
    Bounced     bool   `json:"bounced"`      // Indicates if the message has bounced
    SrcAddr     string `json:"src_addr"`     // Source address in string form
    DstAddr     string `json:"dst_addr"`     // Destination address in string form
    Amount      string `json:"amount"`       // Amount involved in the transaction
    IHRFee      string `json:"ihr_fee"`      // IHR fee in string form
    FwdFee      string `json:"fwd_fee"`      // Forwarding fee in string form
    CreatedLT   uint64 `json:"created_lt"`   // Logical time for transaction creation
    CreatedAt   uint32 `json:"created_at"`   // Creation timestamp
    Body        NFT    `json:"body"`         // Detailed NFT transaction data
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