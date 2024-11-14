package parser

import (
	"encoding/json"
	"testing"

	"github.com/zihaolam/ethereum-parser/internal/ethclient"
)

func TestDeserializeTxn(t *testing.T) {
	// Setup sample transaction in JSON format
	txn := ethclient.Transaction{Hash: "0x12345"}
	txnBytes, err := json.Marshal(txn)
	if err != nil {
		t.Fatalf("unexpected error marshaling transaction: %v", err)
	}

	txBytes := [][]byte{txnBytes}

	// Test deserialization
	txs, err := deserializeTxn(txBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(txs) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txs))
	}

	if txs[0].Hash != txn.Hash {
		t.Fatalf("expected transaction hash %s, got %s", txn.Hash, txs[0].Hash)
	}

	// Test error handling for invalid JSON
	invalidTxBytes := [][]byte{[]byte("invalid json")}
	_, err = deserializeTxn(invalidTxBytes)
	if err == nil {
		t.Fatal("expected error for invalid JSON input")
	}
}

func TestSerializeTxn(t *testing.T) {
	// Setup sample transactions
	txs := []ethclient.Transaction{
		{Hash: "0x12345"},
		{Hash: "0x67890"},
	}

	// Test serialization
	txBytes, err := serializeTxn(txs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(txBytes) != 2 {
		t.Fatalf("expected 2 serialized transactions, got %d", len(txBytes))
	}

	// Test deserialization to confirm correct serialization
	deserializedTxs, err := deserializeTxn(txBytes)
	if err != nil {
		t.Fatalf("unexpected error during deserialization: %v", err)
	}

	for i, tx := range txs {
		if deserializedTxs[i].Hash != tx.Hash {
			t.Fatalf("expected transaction hash %s, got %s", tx.Hash, deserializedTxs[i].Hash)
		}
	}
}

func TestAppendDBTxns(t *testing.T) {
	// Setup initial transaction in JSON format
	initialTxn := ethclient.Transaction{Hash: "0x12345"}
	initialBytes, err := json.Marshal(initialTxn)
	if err != nil {
		t.Fatalf("unexpected error marshaling initial transaction: %v", err)
	}
	txBytes := [][]byte{initialBytes}

	// New transaction to append
	newTxn := ethclient.Transaction{Hash: "0x67890"}

	// Test appending transaction
	updatedTxBytes, err := appendDBTxns(txBytes, []ethclient.Transaction{newTxn})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Deserialize to verify contents
	updatedTxs, err := deserializeTxn(updatedTxBytes)
	if err != nil {
		t.Fatalf("unexpected error during deserialization: %v", err)
	}

	if len(updatedTxs) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(updatedTxs))
	}

	if updatedTxs[0].Hash != initialTxn.Hash {
		t.Fatalf("expected first transaction hash %s, got %s", initialTxn.Hash, updatedTxs[0].Hash)
	}

	if updatedTxs[1].Hash != newTxn.Hash {
		t.Fatalf("expected second transaction hash %s, got %s", newTxn.Hash, updatedTxs[0].Hash)
	}
}
