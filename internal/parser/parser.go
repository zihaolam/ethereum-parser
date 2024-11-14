package parser

import (
	"github.com/zihaolam/ethereum-parser/internal/datastore/memorydb"
	"github.com/zihaolam/ethereum-parser/internal/ethclient"
	"github.com/zihaolam/ethereum-parser/internal/logging"
)

type Parser struct {
	ethClient *ethclient.Client
	logger    logging.Logger
	db        *memorydb.MemoryDB
	*Scanner
}

func NewParser(logger logging.Logger, ethEndpoint string, initialBlockNumber int) *Parser {
	db := memorydb.New()
	ethClient := ethclient.New(ethEndpoint)
	scanner := NewScanner(db, ethClient, logger, initialBlockNumber)
	return &Parser{
		ethClient: ethClient,
		db:        db,
		logger:    logger,
		Scanner:   scanner,
	}
}

func (p *Parser) GetCurrentBlock() int {
	return p.lastBlockNumber
}

func (p *Parser) Subscribe(address string) bool {
	err := p.db.Put(address, [][]byte{})
	if err != nil {
		p.logger.Printf("failed to subscribe to address %s: %v", address, err)
		return false
	}
	return true
}

func (p *Parser) GetTransactions(address string) []ethclient.Transaction {
	v, err := p.db.Get(address)
	if err != nil {
		p.logger.Printf("failed to get transactions for address %s: %v", address, err)
		return nil
	}
	txs, err := deserializeTxn(v)
	return txs
}

func (p *Parser) GetSubscriptions() ([]string, error) {
	addresses, err := p.db.List()
	if err != nil {
		return nil, err
	}
	return addresses, nil
}
