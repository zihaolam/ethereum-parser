package parser

import (
	"encoding/json"

	"github.com/zihaolam/ethereum-parser/internal/ethclient"
)

func deserializeTxn(txBytes [][]byte) ([]ethclient.Transaction, error) {
	txs := make([]ethclient.Transaction, 0, len(txBytes))
	for _, tx := range txBytes {
		var txn ethclient.Transaction
		if err := json.Unmarshal(tx, &txn); err != nil {
			return nil, err
		}
		txs = append(txs, txn)
	}
	return txs, nil
}

func serializeTxn(txs []ethclient.Transaction) ([][]byte, error) {
	txBytes := make([][]byte, 0, len(txs))
	for _, tx := range txs {
		b, err := json.Marshal(tx)
		if err != nil {
			return nil, err
		}
		txBytes = append(txBytes, b)
	}
	return txBytes, nil
}

func appendDBTxns(txBytes [][]byte, newTxs []ethclient.Transaction) ([][]byte, error) {
	txs, err := deserializeTxn(txBytes)
	if err != nil {
		return nil, err
	}
	txs = append(txs, newTxs...)
	return serializeTxn(txs)
}
