package tonsub

import (
	"fmt"

	"github.com/xssnick/tonutils-go/tlb"
)

// RootTON represents a struct for a TON transaction, including core transaction details
// and a nested Ton struct for additional message information.
type RootTON struct {
	OpCode      uint64 `json:"op_code"`      // Operation code indicating the type of transaction
	IHRDisabled bool   `json:"ihr_disabled"` // Flag indicating if Instant Hypercube Routing is disabled
	Bounce      bool   `json:"bounce"`       // Flag indicating if bounces are enabled for this transaction
	Bounced     bool   `json:"bounced"`      // Flag indicating if this message has already bounced
	SrcAddr     string `json:"src_addr"`     // Source address of the transaction, represented as a string
	DstAddr     string `json:"dst_addr"`     // Destination address of the transaction, represented as a string
	Amount      string `json:"amount"`       // Amount involved in the transaction, formatted as a string
	IHRFee      string `json:"ihr_fee"`      // Fee for Instant Hypercube Routing, formatted as a string
	FwdFee      string `json:"fwd_fee"`      // Forward fee in string format
	CreatedLT   uint64 `json:"created_lt"`   // Logical time for when the transaction was created
	CreatedAt   uint32 `json:"created_at"`   // Unix timestamp for the creation of the transaction
	Body        Ton    `json:"body"`         // Nested struct for message-related data
}

// Ton holds additional message details within the transaction,
// specifically capturing any messages encoded in the payload.
type Ton struct {
	Message string `json:"message"` // Optional message extracted from the transaction's payload
}

// TonBody processes a tlb.InternalMessage, extracting its payload
// and mapping the result to a RootTON struct, returning an error if parsing issues occur.
func (s *Sub) TonBody(ti *tlb.InternalMessage) (*RootTON, error) {

	// Begin parsing the payload section of the internal message body
	payload := ti.Body.BeginParse()

	// Initialize a variable to store any extracted message text from the payload
	text := ""
	// Check if the payload exists to proceed with text extraction
	if payload != nil {
		// Attempt to load a 32-bit sumType from the payload
		sumType, err := payload.LoadUInt(32)
		// Continue parsing if the sumType matches the expected value (0x00000000)
		if err == nil && sumType == 0x00000000 {
			// Attempt to load the message text, assuming it's encoded using the Snake format
			value, err := payload.LoadStringSnake()
			if err != nil {
				// Return an error if message text extraction fails
				return nil, fmt.Errorf("failed to load text comment: %v", err)
			}
			text = value // Store the successfully extracted text
		} else if err != nil {
			// Return an error if the sumType cannot be loaded correctly
			return nil, fmt.Errorf("failed to load sumType: %v", err)
		}
	}

	// Construct and return a new RootTON struct populated with extracted data
	return &RootTON{
		OpCode:      0x00000000,          // Default OpCode value
		IHRDisabled: ti.IHRDisabled,      // Copy the IHRDisabled status from internal message
		Bounce:      ti.Bounce,           // Copy the Bounce flag status
		Bounced:     ti.Bounced,          // Copy the Bounced status
		SrcAddr:     ti.SrcAddr.String(), // Convert source address to string
		DstAddr:     ti.DstAddr.String(), // Convert destination address to string
		Amount:      ti.Amount.String(),  // Convert transaction amount to string
		IHRFee:      ti.IHRFee.String(),  // Convert IHR fee to string
		FwdFee:      ti.FwdFee.String(),  // Convert forwarding fee to string
		CreatedLT:   ti.CreatedLT,        // Copy the logical creation time
		CreatedAt:   ti.CreatedAt,        // Copy the creation Unix timestamp
		Body: Ton{
			Message: text, // Include the extracted message text in the Body
		},
	}, nil
}
