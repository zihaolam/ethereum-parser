package parser

import (
	"github.com/zihaolam/ethereum-parser/internal/ethclient"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []ethclient.Transaction
	// get existing subscriptions
}
