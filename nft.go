package tonsub

import (
	"fmt"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/nft"
)

// RootNFT represents the root structure for an NFT transaction,
// containing essential transaction details and a nested NFT structure.
type RootNFT struct {
	IHRDisabled bool   `json:"ihr_disabled"` // Flag indicating if Instant Hypercube Routing is disabled
	Bounce      bool   `json:"bounce"`       // Flag for bounce status of the transaction
	Bounced     bool   `json:"bounced"`      // Flag indicating if the transaction has already bounced
	SrcAddr     string `json:"src_addr"`     // Source address of the transaction as a string
	DstAddr     string `json:"dst_addr"`     // Destination address of the transaction as a string
	Amount      string `json:"amount"`       // Amount involved in the transaction represented as a string
	IHRFee      string `json:"ihr_fee"`      // Instant Hypercube Routing fee as a string
	FwdFee      string `json:"fwd_fee"`      // Forward fee as a string
	CreatedLT   uint64 `json:"created_lt"`   // Logical time indicating when the transaction was created
	CreatedAt   uint32 `json:"created_at"`   // Unix timestamp indicating when the transaction was created
	Body        NFT    `json:"body"`         // Nested NFT structure containing detailed NFT data
}

// NFT represents details specific to the NFT itself, including
// operational codes, ownership, and associated addresses.
type NFT struct {
	OpCode            uint64 `json:"op_code"`            // Operation code of the NFT transaction
	Initialized       bool   `json:"initialized"`        // Indicates if the NFT item is initialized
	Index             string `json:"index"`              // Index of the NFT in the collection
	Address           string `json:"nft_address"`        // Address of the NFT
	OwnerAddress      string `json:"owner_address"`      // Address of the current owner of the NFT
	CollectionAddress string `json:"collection_address"` // Address of the collection to which the NFT belongs
	Meta              string `json:"meta"`
	Message           string `json:"message"` // Optional message extracted from the transaction
}

// NFTBody processes a tlb.InternalMessage and extracts NFT transaction data,
// returning a filled RootNFT structure or an error if processing fails.
func (s *Sub) NFTBody(ti *tlb.InternalMessage) (*RootNFT, error) {

	// Parse the source address of the transaction to create an NFT item client.
	nftAddr := address.MustParseAddr(ti.SrcAddr.String())
	item := nft.NewItemClient(s.Api, nftAddr)

	// Retrieve NFT data, which contains various information related to the NFT.
	nftData, err := item.GetNFTData(s.Context)
	if err != nil {
		return nil, err
	}

	var meta string
	if nftData.CollectionAddress.Type() != address.NoneAddress {
		// get info about our nft's collection
		collection := nft.NewCollectionClient(s.Api, nftData.CollectionAddress)

		// get full nft's content url using collection method that will merge base url with nft's data
		nftContent, err := collection.GetNFTContent(s.Context, nftData.Index, nftData.Content)
		if err != nil {
			return nil, err
		}

		if off, ok := nftContent.(*nft.ContentOffchain); ok {
			meta = off.URI
		}
	}

	// Initialize payload parsing to extract additional transaction details.
	payload := ti.Body.BeginParse()

	// Initialize the message text to an empty string.
	text := ""
	// Check if there is a payload to process for message extraction.
	if payload != nil {
		// Attempt to load the sumType from the payload to determine the message format.
		sumType, err := payload.LoadUInt(32)
		if err == nil && sumType == 0x00000000 {
			// If sumType matches, attempt to load the message in Snake format.
			value, err := payload.LoadStringSnake()
			if err != nil {
				// Return an error if the text loading fails.
				return nil, fmt.Errorf("failed to load text comment: %v", err)
			}
			// Assign the loaded text to the message field.
			text = value
		} else if err != nil {
			// Return an error if loading sumType fails.
			return nil, fmt.Errorf("failed to load sumType: %v", err)
		}
	}

	// Return composed RootNFT structure with all processed transaction data.
	return &RootNFT{
		IHRDisabled: ti.IHRDisabled,      // Indicates if IHR is disabled
		Bounce:      ti.Bounce,           // Bounce status
		Bounced:     ti.Bounced,          // Bounced status
		SrcAddr:     ti.SrcAddr.String(), // Source address as string
		DstAddr:     ti.DstAddr.String(), // Destination address as string
		Amount:      ti.Amount.String(),  // Amount as string
		IHRFee:      ti.IHRFee.String(),  // IHR fee as string
		FwdFee:      ti.FwdFee.String(),  // Forwarding fee as string
		CreatedLT:   ti.CreatedLT,        // Created logical time
		CreatedAt:   ti.CreatedAt,        // Unix creation timestamp
		Body: NFT{
			OpCode:            0x7362d09c, // Operation code
			Initialized:       nftData.Initialized,
			Index:             nftData.Index.String(),
			Address:           nftAddr.String(),
			OwnerAddress:      nftData.OwnerAddress.String(),
			CollectionAddress: nftData.CollectionAddress.String(),
			Meta:              meta,
			Message:           text,
		},
	}, nil
}
