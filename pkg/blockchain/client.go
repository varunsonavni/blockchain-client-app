package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	PolygonRPC = "https://polygon-rpc.com/"
)

// Client represents a blockchain client
type Client struct {
	httpClient *http.Client
	rpcURL     string
}

// RPCRequest represents a JSON-RPC request
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
	ID      int           `json:"id"`
}

// RPCResponse represents a JSON-RPC response
type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
	ID      int             `json:"id"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewClient creates a new blockchain client
func NewClient(rpcURL string) *Client {
	if rpcURL == "" {
		rpcURL = PolygonRPC
	}
	return &Client{
		httpClient: &http.Client{},
		rpcURL:     rpcURL,
	}
}

// call makes an RPC call to the blockchain
func (c *Client) call(method string, params []interface{}) (*RPCResponse, error) {
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      2,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.rpcURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rpcResp RPCResponse
	if err := json.Unmarshal(bodyBytes, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error: %s (code: %d)", rpcResp.Error.Message, rpcResp.Error.Code)
	}

	return &rpcResp, nil
}

// GetBlockNumber returns the latest block number
func (c *Client) GetBlockNumber() (string, error) {
	resp, err := c.call("eth_blockNumber", nil)
	if err != nil {
		return "", err
	}

	var blockNumber string
	if err := json.Unmarshal(resp.Result, &blockNumber); err != nil {
		return "", fmt.Errorf("failed to unmarshal block number: %w", err)
	}

	return blockNumber, nil
}

// Block represents an Ethereum block
type Block struct {
	Number           string          `json:"number"`
	Hash             string          `json:"hash"`
	ParentHash       string          `json:"parentHash"`
	Nonce            string          `json:"nonce"`
	Timestamp        string          `json:"timestamp"`
	Transactions     json.RawMessage `json:"transactions"`
	TransactionCount int             `json:"transactionCount"`
}

// GetBlockByNumber returns the block information by block number
func (c *Client) GetBlockByNumber(blockNumber string, fullTransactions bool) (*Block, error) {
	resp, err := c.call("eth_getBlockByNumber", []interface{}{blockNumber, fullTransactions})
	if err != nil {
		return nil, err
	}

	var block Block
	if err := json.Unmarshal(resp.Result, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}

	var txCount int
	if fullTransactions {
		var txArray []interface{}
		if err := json.Unmarshal(block.Transactions, &txArray); err == nil {
			txCount = len(txArray)
		}
	} else {
		var txArray []string
		if err := json.Unmarshal(block.Transactions, &txArray); err == nil {
			txCount = len(txArray)
		}
	}

	block.TransactionCount = txCount
	return &block, nil
}

// Create a helper method for tests
// CreateMockBlock creates a block with the given data for testing purposes
func CreateMockBlock(number, hash, parentHash, nonce, timestamp string, txCount int, txData json.RawMessage) *Block {
	return &Block{
		Number:           number,
		Hash:             hash,
		ParentHash:       parentHash,
		Nonce:            nonce,
		Timestamp:        timestamp,
		Transactions:     txData,
		TransactionCount: txCount,
	}
}
