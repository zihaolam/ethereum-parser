package parser

import (
	"context"
	"time"

	"github.com/zihaolam/ethereum-parser/internal/datastore"
	"github.com/zihaolam/ethereum-parser/internal/ethclient"
	"github.com/zihaolam/ethereum-parser/internal/logging"
)

type Scanner struct {
	db              datastore.DataStore
	logger          logging.Logger
	ethClient       *ethclient.Client
	lastBlockNumber int
}

func NewScanner(
	db datastore.DataStore,
	ethClient *ethclient.Client,
	logger logging.Logger,
	initialBlockNumber int,
) *Scanner {
	return &Scanner{
		db:              db,
		ethClient:       ethClient,
		logger:          logger,
		lastBlockNumber: initialBlockNumber,
	}
}

func (b *Scanner) SaveTxs(txs []ethclient.Transaction) error {
	txMap := make(map[string][]ethclient.Transaction)

	for _, tx := range txs {
		txMap[tx.From] = append(txMap[tx.From], tx)
		txMap[tx.To] = append(txMap[tx.To], tx)
	}

	for addr, txs := range txMap {
		if err := b.db.Update(addr, func(oldTxs [][]byte) ([][]byte, error) {
			return appendDBTxns(oldTxs, txs)
		}); err != nil {
			return err
		}
	}

	return nil
}

func (b *Scanner) FilterSubscribedTxs(txs []ethclient.Transaction) []ethclient.Transaction {
	subscribedTxs := make([]ethclient.Transaction, 0)
	for _, tx := range txs {
		if b.db.Has(tx.To) || b.db.Has(tx.From) {
			subscribedTxs = append(subscribedTxs, tx)
		}
	}

	return subscribedTxs
}

func (b *Scanner) SaveTxsToSubscribers(txs []ethclient.Transaction) error {
	subscribedTxs := b.FilterSubscribedTxs(txs)
	return b.SaveTxs(subscribedTxs)
}

func (b *Scanner) ScanBlock(ctx context.Context, blockNumber int) (ethclient.Block, error) {
	block, err := b.ethClient.GetBlockByNumber(ctx, blockNumber)
	if err != nil {
		return ethclient.Block{}, err
	}

	return block, nil
}

// returns 0 if there are no more new blocks else returns the next block number
func (b *Scanner) GetNextBlock(ctx context.Context) (int, error) {
	currBlockNumber, err := b.ethClient.GetCurrentBlockNumber(ctx)
	if err != nil {
		b.logger.Printf("error getting current block number: %v", err)
		return 0, err
	}

	if b.lastBlockNumber == currBlockNumber {
		return 0, nil
	}

	if b.lastBlockNumber == 0 {
		return currBlockNumber, nil
	}

	return b.lastBlockNumber + 1, nil
}

// Scan checks for new blocks and saves transactions to the datastore.
func (b *Scanner) ScanAll(ctx context.Context) error {
	// scan and save new blocks
	for {
		nextBlock, err := b.GetNextBlock(ctx)
		if err != nil {
			return err
		}
		if nextBlock == 0 {
			break
		}

		b.logger.Printf("Scanning block %d\n", nextBlock)
		block, err := b.ScanBlock(ctx, nextBlock)
		if err != nil {
			return err
		}
		b.logger.Println("Saving transactions")
		err = b.SaveTxsToSubscribers(block.Transactions)
		if err != nil {
			return err
		}

		b.lastBlockNumber = nextBlock
	}

	return nil
}

// Interval scan for new blocks
func (b *Scanner) StartScan(ctx context.Context, interval time.Duration) {
	timer1 := time.NewTimer(interval)
	for {
		select {
		case <-ctx.Done():
			b.logger.Println("stopping scanner")
			return
		case <-timer1.C:
			go func() {
				err := b.ScanAll(ctx)
				b.logger.Printf("Scanning for new blocks\n")
				if err != nil {
					b.logger.Printf("error scanning: %v", err)
				}
				timer1.Reset(10 * time.Second)
			}()
		}
	}
}
