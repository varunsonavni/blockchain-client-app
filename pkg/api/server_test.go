package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"blockchain-client/pkg/blockchain"
)

// mockBlockchainClient is a mock implementation of the blockchain client for testing
type mockBlockchainClient struct {
	getBlockNumberFunc   func() (string, error)
	getBlockByNumberFunc func(blockNumber string, fullTransactions bool) (*blockchain.Block, error)
}

func (m *mockBlockchainClient) GetBlockNumber() (string, error) {
	return m.getBlockNumberFunc()
}

func (m *mockBlockchainClient) GetBlockByNumber(blockNumber string, fullTransactions bool) (*blockchain.Block, error) {
	return m.getBlockByNumberFunc(blockNumber, fullTransactions)
}

// We need to modify the Server struct in tests to accept the interface instead of the concrete type
type blockchainClient interface {
	GetBlockNumber() (string, error)
	GetBlockByNumber(blockNumber string, fullTransactions bool) (*blockchain.Block, error)
}

// testServer wraps Server for testing
type testServer struct {
	server *Server
	mock   *mockBlockchainClient
}

// newTestServer creates a new Server with a mock client for testing
func newTestServer() *testServer {
	mock := &mockBlockchainClient{}

	server := &Server{
		client: mock,
	}

	return &testServer{server: server, mock: mock}
}

func TestHandleGetBlockNumber(t *testing.T) {
	// Create a test server with mock client
	ts := newTestServer()

	// Set up mock response
	ts.mock.getBlockNumberFunc = func() (string, error) {
		return "0x1234567", nil
	}

	// Create a request
	req, err := http.NewRequest("GET", "/api/blocks/latest", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Create a recorder to record the response
	rec := httptest.NewRecorder()

	// Call the handler
	ts.server.HandleGetBlockNumber(rec, req)

	// Check the status code
	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rec.Code)
	}

	// Check the response body
	var resp BlockNumberResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if resp.BlockNumber != "0x1234567" {
		t.Errorf("expected block number 0x1234567; got %v", resp.BlockNumber)
	}
}

func TestHandleGetBlockByNumber(t *testing.T) {
	// Create a test server with mock client
	ts := newTestServer()

	// Set up mock response with proper json.RawMessage for transactions
	ts.mock.getBlockByNumberFunc = func(blockNumber string, fullTransactions bool) (*blockchain.Block, error) {
		// Create raw JSON for transactions
		var txsJSON json.RawMessage
		if fullTransactions {
			txsJSON = json.RawMessage(`[
				{"hash": "0xtx1", "from": "0xaddr1", "to": "0xaddr2"},
				{"hash": "0xtx2", "from": "0xaddr3", "to": "0xaddr4"}
			]`)
		} else {
			txsJSON = json.RawMessage(`["0xtx1", "0xtx2"]`)
		}

		// Use helper to create a block with transaction count
		return blockchain.CreateMockBlock(
			"0x1234567",
			"0xabcdef1234567890",
			"0x1234567890abcdef",
			"0x123456",
			"0x60123456",
			2,
			txsJSON,
		), nil
	}

	// Create a request
	req, err := http.NewRequest("GET", "/api/blocks?number=0x1234567&full=true", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Create a recorder to record the response
	rec := httptest.NewRecorder()

	// Call the handler
	ts.server.HandleGetBlockByNumber(rec, req)

	// Check the status code
	if rec.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rec.Code)
	}

	// Print raw response for debugging
	respBody, _ := ioutil.ReadAll(rec.Body)
	t.Logf("Raw response: %s", string(respBody))

	// Create a new response recorder with the same body for testing
	newRec := httptest.NewRecorder()
	newRec.Write(respBody)

	// Check the response body
	var resp BlockResponse
	if err := json.NewDecoder(newRec.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if resp.Block == nil {
		t.Fatalf("expected block to not be nil")
	}

	// Debug the received block
	blockBytes, _ := json.Marshal(resp.Block)
	t.Logf("Block: %s", string(blockBytes))

	if resp.Block.Number != "0x1234567" {
		t.Errorf("expected block number 0x1234567; got %v", resp.Block.Number)
	}

	if resp.Block.Hash != "0xabcdef1234567890" {
		t.Errorf("expected block hash 0xabcdef1234567890; got %v", resp.Block.Hash)
	}

	if resp.Block.TransactionCount != 2 {
		t.Errorf("expected 2 transactions; got %v", resp.Block.TransactionCount)
	}
}

func TestHandleJSONRPC(t *testing.T) {
	t.Run("eth_blockNumber", func(t *testing.T) {
		// Create a test server with mock client
		ts := newTestServer()

		// Set up mock response
		ts.mock.getBlockNumberFunc = func() (string, error) {
			return "0x1234567", nil
		}

		// Create a JSON-RPC request
		reqBody := `{
			"jsonrpc": "2.0",
			"method": "eth_blockNumber",
			"id": 2
		}`

		// Create a request
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(reqBody))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create a recorder to record the response
		rec := httptest.NewRecorder()

		// Call the handler
		ts.server.HandleJSONRPC(rec, req)

		// Check the status code
		if rec.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", rec.Code)
		}

		// Check the response body
		var resp RPCResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}

		if resp.JSONRPC != "2.0" {
			t.Errorf("expected jsonrpc 2.0; got %v", resp.JSONRPC)
		}

		if resp.ID != 2 {
			t.Errorf("expected id 2; got %v", resp.ID)
		}

		if resp.Error != nil {
			t.Errorf("expected no error; got %v", resp.Error)
		}

		// Extract the block number
		var blockNumber string
		if err := json.Unmarshal(resp.Result, &blockNumber); err != nil {
			t.Fatalf("could not unmarshal result: %v", err)
		}

		if blockNumber != "0x1234567" {
			t.Errorf("expected block number 0x1234567; got %v", blockNumber)
		}
	})

	t.Run("eth_getBlockByNumber", func(t *testing.T) {
		// Create a test server with mock client
		ts := newTestServer()

		// Set up mock response
		ts.mock.getBlockByNumberFunc = func(blockNumber string, fullTransactions bool) (*blockchain.Block, error) {
			if blockNumber != "0x1234567" {
				t.Errorf("expected block number 0x1234567; got %v", blockNumber)
			}
			if !fullTransactions {
				t.Errorf("expected full transactions true; got false")
			}

			// Create raw JSON for transactions
			txsJSON := json.RawMessage(`[
				{"hash": "0xtx1", "from": "0xaddr1", "to": "0xaddr2"},
				{"hash": "0xtx2", "from": "0xaddr3", "to": "0xaddr4"}
			]`)

			return &blockchain.Block{
				Number:           "0x1234567",
				Hash:             "0xabcdef1234567890",
				ParentHash:       "0x1234567890abcdef",
				Nonce:            "0x123456",
				Timestamp:        "0x60123456",
				Transactions:     txsJSON,
				TransactionCount: 2,
			}, nil
		}

		// Create a JSON-RPC request
		reqBody := `{
			"jsonrpc": "2.0",
			"method": "eth_getBlockByNumber",
			"params": [
				"0x1234567",
				true
			],
			"id": 2
		}`

		// Create a request
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(reqBody))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create a recorder to record the response
		rec := httptest.NewRecorder()

		// Call the handler
		ts.server.HandleJSONRPC(rec, req)

		// Check the status code
		if rec.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", rec.Code)
		}

		// Check the response body
		var resp RPCResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}

		if resp.JSONRPC != "2.0" {
			t.Errorf("expected jsonrpc 2.0; got %v", resp.JSONRPC)
		}

		if resp.ID != 2 {
			t.Errorf("expected id 2; got %v", resp.ID)
		}

		if resp.Error != nil {
			t.Errorf("expected no error; got %v", resp.Error)
		}

		// Extract the block
		var block map[string]interface{}
		if err := json.Unmarshal(resp.Result, &block); err != nil {
			t.Fatalf("could not unmarshal result: %v", err)
		}

		if block["number"] != "0x1234567" {
			t.Errorf("expected block number 0x1234567; got %v", block["number"])
		}

		if block["hash"] != "0xabcdef1234567890" {
			t.Errorf("expected block hash 0xabcdef1234567890; got %v", block["hash"])
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		// Create a test server with mock client
		ts := newTestServer()

		// Create a JSON-RPC request with invalid method
		reqBody := `{
			"jsonrpc": "2.0",
			"method": "invalid_method",
			"id": 2
		}`

		// Create a request
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(reqBody))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create a recorder to record the response
		rec := httptest.NewRecorder()

		// Call the handler
		ts.server.HandleJSONRPC(rec, req)

		// Check the status code
		if rec.Code != http.StatusOK {
			t.Errorf("expected status OK; got %v", rec.Code)
		}

		// Check the response body
		var resp RPCResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("could not decode response: %v", err)
		}

		if resp.JSONRPC != "2.0" {
			t.Errorf("expected jsonrpc 2.0; got %v", resp.JSONRPC)
		}

		if resp.ID != 2 {
			t.Errorf("expected id 2; got %v", resp.ID)
		}

		if resp.Error == nil {
			t.Errorf("expected error; got nil")
		} else if resp.Error.Code != -32601 {
			t.Errorf("expected error code -32601; got %v", resp.Error.Code)
		}
	})
}

func TestDirectBlockJson(t *testing.T) {
	// Create a Block with transaction count explicitly set
	txsJSON := json.RawMessage(`["0xtx1", "0xtx2"]`)

	// Create a Block via helper
	block := blockchain.CreateMockBlock(
		"0x1234567",
		"0xabcdef1234567890",
		"0x1234567890abcdef",
		"0x123456",
		"0x60123456",
		2,
		txsJSON,
	)

	// Check transaction count directly in the struct
	if block.TransactionCount != 2 {
		t.Errorf("expected TransactionCount 2 in original block; got %v", block.TransactionCount)
	}

	// Test marshaling and unmarshaling with block response
	resp := BlockResponse{Block: block}

	// Marshal to JSON
	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal block response: %v", err)
	}

	t.Logf("Marshaled JSON: %s", string(jsonData))

	// Unmarshal back
	var newResp BlockResponse
	if err := json.Unmarshal(jsonData, &newResp); err != nil {
		t.Fatalf("Failed to unmarshal block response: %v", err)
	}

	// Check TransactionCount after round trip
	if newResp.Block.TransactionCount != 2 {
		t.Errorf("expected TransactionCount 2 after unmarshal; got %v", newResp.Block.TransactionCount)
	}
}
