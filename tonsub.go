package tonsub

import (
	"context"
	"fmt"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
)

var core *Sub // Global instance of Sub structure

// Sub represents a subscription manager for different transaction types within the TON blockchain network.
type Sub struct {
	Context context.Context      // Context used for binding requests to a single TON node
	Block   *ton.BlockIDExt      // Current block being processed
	Api     ton.APIClientWrapped // Wrapped API client for interacting with the TON blockchain
	lt      uint64               // Last transaction logical time (LT)

	clbJetton func(t *RootJetton) // Callback function for Jetton transactions
	clbTon    func(t *RootTON)    // Callback function for TON transactions
	clbNFT    func(t *RootNFT)    // Callback function for NFT transactions
}

// OnJetton registers a callback function for Jetton transactions.
func (s *Sub) OnJetton(clb func(t *RootJetton)) {
	s.clbJetton = clb
}

// OnTON registers a callback function for TON transactions.
func (s *Sub) OnTON(clb func(t *RootTON)) {
	s.clbTon = clb
}

// OnNFT registers a callback function for NFT transactions.
func (s *Sub) OnNFT(clb func(t *RootNFT)) {
	s.clbNFT = clb
}

// New initializes a new Sub instance with the provided address and network configuration.
func New(addr string, network string) (*Sub, error) {
	client := liteclient.NewConnectionPool()

	// Create a sticky context to bind all requests to a single TON node.
	ctx := client.StickyContext(context.Background())

	// Retrieve configuration from the given network URL for setting up connections.
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), network)
	if err != nil {
		return nil, err
	}

	// Connect to the mainnet lite servers using the fetched configuration.
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	// Initialize the API client with proof checks enabled.
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	// Fetch the current masterchain block information for integrity checks.
	master, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}

	// Parse the address for processing payments.
	treasuryAddress := address.MustParseAddr(addr)

	// Retrieve account information linked to the parsed address.
	acc, err := api.GetAccount(context.Background(), master, treasuryAddress)
	if err != nil {
		return nil, err
	}

	// Fetch the current masterchain block information again for the subscription setup.
	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}

	// Initialize core Sub instance with default callback functions and set state.
	core = &Sub{
		Context:   ctx,
		Block:     block,
		Api:       api,
		lt:        acc.LastTxLT,
		clbJetton: func(t *RootJetton) {},
		clbTon:    func(t *RootTON) {},
		clbNFT:    func(t *RootNFT) {},
	}

	// Create a channel to receive new transactions.
	transactions := make(chan *tlb.Transaction)

	// Start transaction subscription processing asynchronously.
	go core.subscribe(transactions)

	// Asynchronously start subscription to transactions on the specified address.
	go api.SubscribeOnTransactions(context.Background(), treasuryAddress, core.lt, transactions)

	return core, nil
}

// subscribe listens to the transaction channel and processes transactions based on their type.
func (s *Sub) subscribe(channel chan *tlb.Transaction) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			s.subscribe(channel) // Restart subscription on panic
		}
	}()

	for tx := range channel {
		// Process only internal messages
		if tx.IO.In != nil && tx.IO.In.MsgType == tlb.MsgTypeInternal {
			// Skip if the transaction has successfully bounced
			if dsc, ok := tx.Description.(tlb.TransactionDescriptionOrdinary); ok && dsc.BouncePhase != nil {
				if _, ok = dsc.BouncePhase.Phase.(tlb.BouncePhaseOk); ok {
					continue
				}
			}

			// Parse the internal message body
			ti := tx.IO.In.AsInternal()

			bodyCell := ti.Body.BeginParse()
			if bodyCell != nil {
				// Extract the operation code from the body
				opCode, err := bodyCell.LoadUInt(32)
				if err != nil {
					opCode = 0x00000000 // Default to zero if none found
				}

				// Call the appropriate callback based on the operation code
				switch opCode {
				case 0x05138d91: // NFT transfer opcode
					body, err := s.NFTBody(ti)
					if err != nil {
						continue // Skip on error
					}
					s.clbNFT(body)

				case 0x7362d09c: // Jetton transfer opcode
					body, err := s.JettonBody(ti)
					if err != nil {
						continue // Skip on error
					}
					s.clbJetton(body)

				case 0x00000000: // TON transfer opcode
					body, err := s.TonBody(ti)
					if err != nil {
						continue // Skip on error
					}
					s.clbTon(body)
				}
			}
		}
	}
}
