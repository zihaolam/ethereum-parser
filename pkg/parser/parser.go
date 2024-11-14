package parser

import (
	"context"
	"time"

	"github.com/zihaolam/ethereum-parser/internal/ethclient"
	"github.com/zihaolam/ethereum-parser/internal/logging"
	"github.com/zihaolam/ethereum-parser/internal/parser"
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

func NewParser(
	logger logging.Logger,
	ethEndpoint string,
	initialBlockNumber int,
) (Parser, func(ctx context.Context, interval time.Duration)) {
	p := parser.NewParser(logger, ethEndpoint, initialBlockNumber)
	return p, p.StartScan
}
