package ethclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

const (
	ApiVersion                  = "2.0"
	GetCurrentBlocknumberMethod = "eth_blockNumber"
	GetBlockByNumberMethod      = "eth_getBlockByNumber"
)

type Client struct {
	endpoint string
}

type RequestBody struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type ResponseBody[T any] struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  T      `json:"result"`
}

type Block struct {
	Number       string        `json:"number"`
	Hash         string        `json:"hash"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	ChainID          string            `json:"chainId"`
	BlockNumber      string            `json:"blockNumber"`
	BlockHash        string            `json:"-"`
	Hash             string            `json:"hash"`
	Nonce            string            `json:"nonce"`
	From             string            `json:"from"`
	To               string            `json:"to"`
	Value            string            `json:"value"`
	Gas              string            `json:"gas"`
	GasPrice         string            `json:"gasPrice"`
	Input            string            `json:"input"`
	Type             string            `json:"-"`
	R                string            `json:"-"`
	S                string            `json:"-"`
	V                string            `json:"-"`
	TransactionIndex string            `json:"-"`
	AccessList       []AccessListEntry `json:"-"`
}

type AccessListEntry struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

func New(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

func sendRPC[T any](
	ctx context.Context,
	endpoint string,
	rbody RequestBody,
) (ResponseBody[T], error) {
	body, err := json.Marshal(rbody)
	if err != nil {
		return ResponseBody[T]{}, fmt.Errorf("error marshaling json: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return ResponseBody[T]{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return ResponseBody[T]{}, fmt.Errorf("error making request: %v", err)
	}

	if r.StatusCode != http.StatusOK {
		return ResponseBody[T]{}, fmt.Errorf("error response status code: %v", r.StatusCode)
	}

	var responseBody ResponseBody[T]

	if err := json.NewDecoder(r.Body).Decode(&responseBody); err != nil {
		return ResponseBody[T]{}, fmt.Errorf("error decoding response body: %v", err)
	}

	return responseBody, nil
}

// Returns the current block number.
// It calls the JSON-RPC eth_blockNumber method.
func (c Client) GetCurrentBlockNumber(ctx context.Context) (int, error) {
	body := makeRequestBody(GetCurrentBlocknumberMethod, []string{})
	res, err := sendRPC[string](ctx, c.endpoint, body)

	if err != nil {
		return 0, fmt.Errorf("error sending rpc: %v", err)
	}

	blocknumber, err := strconv.ParseInt(res.Result[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing response body: %v", err)
	}

	return int(blocknumber), nil
}

// Returns the block information for the given block number.
// It calls the JSON-RPC eth_getBlockByNumber method.
func (c Client) GetBlockByNumber(ctx context.Context, blocknumber int) (Block, error) {
	body := makeRequestBody(
		GetBlockByNumberMethod,
		[]string{fmt.Sprintf("0x%x", blocknumber), "true"},
	)
	res, err := sendRPC[Block](ctx, c.endpoint, body)

	if err != nil {
		return Block{}, fmt.Errorf("error sending rpc: %v", err)
	}

	return res.Result, nil
}

func makeRequestBody(method string, params interface{}) RequestBody {
	return RequestBody{
		Jsonrpc: ApiVersion,
		ID:      rand.Int(),
		Method:  method,
		Params:  params,
	}
}
