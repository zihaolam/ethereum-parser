package parser_test

import (
	"context"
	"log"
	"testing"

	"github.com/zihaolam/ethereum-parser/internal/ethclient"
	"github.com/zihaolam/ethereum-parser/internal/parser"
)

const (
	ethEndpoint  = "https://cloudflare-eth.com"
	initialBlock = 17758000
)

// Test for GetCurrentBlock using Cloudflare Ethereum Gateway
func TestParser(t *testing.T) {
	logger := log.Default()
	p := parser.NewParser(logger, ethEndpoint, initialBlock)

	t.Run("GetCurrentBlock", func(t *testing.T) {
		if p.GetCurrentBlock() != initialBlock {
			t.Fatalf(
				"Expected current block: %d. Got: %d ",
				initialBlock,
				p.GetCurrentBlock(),
			)
		}
	})

	t.Run("GetSubscriptions", func(t *testing.T) {
		p.Subscribe("0xAddress")
		subscribers, err := p.GetSubscriptions()
		if err != nil {
			t.Errorf("Error getting subscribers: %v", err)
		}

		for _, subscriber := range subscribers {
			if subscriber == "0xAddress" {
				return
			}
		}
		t.Errorf("Expected address to be in the subscriber list")
	})

	t.Run("Subscribe", func(t *testing.T) {
		result := p.Subscribe("0xAddress")
		if !result {
			t.Error("Expected subscription to succeed")
		}

		// Verify that the address exists in the database
		subscribers, err := p.GetSubscriptions()
		if err != nil {
			t.Errorf("Error getting subscribers: %v", err)
		}
		for _, subscriber := range subscribers {
			if subscriber == "0xAddress" {
				return
			}
		}
		t.Errorf("Expected address to be in the subscriber list")
	})

	t.Run("ScanBlock", func(t *testing.T) {
		_, err := p.ScanBlock(context.Background(), initialBlock)
		if err != nil {
			t.Errorf("Error scanning block: %v", err)
		}
	})

	t.Run("GetNextBlock", func(t *testing.T) {
		blockNumber, err := p.GetNextBlock(context.Background())
		if err != nil {
			t.Errorf("Error getting next block: %v", err)
		}
		if blockNumber != initialBlock+1 {
			t.Errorf("Expected next block to be %d, got %d", initialBlock+1, blockNumber)
		}
	})

	var txMap map[string][]ethclient.Transaction

	t.Run("SaveTxs", func(t *testing.T) {
		block, err := p.ScanBlock(context.Background(), initialBlock)
		if err != nil {
			t.Errorf("Error scanning block: %v", err)
		}

		for _, tx := range block.Transactions {
			txMap[tx.From] = append(txMap[tx.From], tx)
			txMap[tx.To] = append(txMap[tx.To], tx)
		}

		p.SaveTxs(block.Transactions)

		for addr, txs := range txMap {
			savedTxs := p.GetTransactions(addr)
			if len(savedTxs) != len(txs) {
				t.Errorf("Expected %d transactions, got %d", len(txs), len(savedTxs))
			}
		}
	})

	t.Run("GetTransactions", func(t *testing.T) {
		for addr, txs := range txMap {
			savedTxs := p.GetTransactions(addr)
			if len(savedTxs) != len(txs) {
				t.Errorf("Expected %d transactions, got %d", len(txs), len(savedTxs))
			}
		}
	})

	t.Run("SaveTxsToSubscribers", func(t *testing.T) {
		block, err := p.ScanBlock(context.Background(), initialBlock)
		if err != nil {
			t.Errorf("Error scanning block: %v", err)
		}
		txMap = make(map[string][]ethclient.Transaction)
		var oneAddr string
		for _, tx := range block.Transactions {
			txMap[tx.From] = append(txMap[tx.From], tx)
			txMap[tx.To] = append(txMap[tx.To], tx)
			oneAddr = tx.From
		}
		p.Subscribe(oneAddr)
		p.SaveTxsToSubscribers(block.Transactions)
		txs := p.GetTransactions(oneAddr)
		if len(txs) != len(txMap[oneAddr]) {
			t.Errorf("Expected %d transactions, got %d", len(txMap[oneAddr]), len(txs))
		}
	})
}
