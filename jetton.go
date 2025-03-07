package tonsub

import (
	"errors"
	"fmt"

	"github.com/xssnick/tonutils-go/tlb"
)

// RootJetton represents the root structure for a Jetton transaction,
// including basic transaction details and a nested Jetton structure.
type RootJetton struct {
	IHRDisabled bool   `json:"ihr_disabled"` // Indicator if Instant Hypercube Routing is disabled
	Bounce      bool   `json:"bounce"`       // Indicator if the message should bounce back if it fails
	Bounced     bool   `json:"bounced"`      // Indicator if the message has already bounced
	SrcAddr     string `json:"src_addr"`     // Source address as a string
	DstAddr     string `json:"dst_addr"`     // Destination address as a string
	Amount      string `json:"amount"`       // Amount involved in the transaction as a string
	IHRFee      string `json:"ihr_fee"`      // Instant Hypercube Routing fee as a string
	FwdFee      string `json:"fwd_fee"`      // Forwarding fee as a string
	CreatedLT   uint64 `json:"created_lt"`   // Logical time when the message was created
	CreatedAt   uint32 `json:"created_at"`   // Unix timestamp when the message was created
	Body        Jetton `json:"body"`         // Contained Jetton information
}

// Jetton represents the details of a Jetton transaction body,
// including operation code, query ID, amount, sender, and an optional message.
type Jetton struct {
	OpCode  uint64 `json:"op_code"`  // Operation code for the transaction
	QueryID uint64 `json:"query_id"` // Unique identifier for the query
	Amount  string `json:"amount"`   // Amount of Jetton involved in the transaction
	Sender  string `json:"sender"`   // Address of the sender
	Message string `json:"message"`  // Optional message included in the transaction
}

// JettonBody parses a tlb.InternalMessage containing Jetton transaction data
// into a RootJetton struct. It returns an error if any part of the parsing fails.
func (s *Sub) JettonBody(ti *tlb.InternalMessage) (*RootJetton, error) {
	// Check if the internal message is nil and return an error if it is.
	if ti == nil {
		return nil, errors.New("input internal message is nil")
	}

	// Begin parsing the body of the internal message
	slice := ti.Body.BeginParse()

	// Attempt to load a 32-bit OpCode from the slice
	opCode, err := slice.LoadUInt(32)
	if err != nil {
		// Default the OpCode to 0x00000000 if loading fails
		opCode = 0x00000000
	}

	// Load a 64-bit QueryID from the slice, returning an error if it fails
	queryID, err := slice.LoadUInt(64)
	if err != nil {
		return nil, fmt.Errorf("failed to load QueryID: %v", err)
	}

	// Load BigInt amount from the slice, returning an error if it fails
	amount, err := slice.LoadBigCoins()
	if err != nil {
		return nil, fmt.Errorf("failed to load amount: %v", err)
	}

	// Load the sender's address from the slice, returning an error if it fails
	sender, err := slice.LoadAddr()
	if err != nil {
		return nil, fmt.Errorf("failed to load sender address: %v", err)
	}

	// Load a payload reference, which may be optional
	payload, err := slice.LoadMaybeRef()
	if err != nil {
		return nil, fmt.Errorf("failed to load payload reference: %v", err)
	}

	// Initialize an empty string for the message text
	text := ""
	// If the payload is not nil, attempt to extract additional information
	if payload != nil {
		// Load a 32-bit sumType from the payload
		sumType, err := payload.LoadUInt(32)
		if err == nil && sumType == 0x00000000 {
			// If sumType matches, load the message text in Snake format
			value, err := payload.LoadStringSnake()
			if err != nil {
				return nil, fmt.Errorf("failed to load text comment: %v", err)
			}
			text = value // Store the extracted text
		} else if err != nil {
			return nil, fmt.Errorf("failed to load sumType: %v", err)
		}
	}

	// Return a populated RootJetton struct containing all parsed data
	return &RootJetton{
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
		Body: Jetton{
			OpCode:  opCode,          // Operation code
			QueryID: queryID,         // Query ID
			Amount:  amount.String(), // Amount as string
			Sender:  sender.String(), // Sender as string
			Message: text,            // Extracted message
		},
	}, nil
}
