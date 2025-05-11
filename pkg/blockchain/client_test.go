package blockchain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBlockNumber(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate request
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		var rpcReq RPCRequest
		if err := json.NewDecoder(r.Body).Decode(&rpcReq); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if rpcReq.Method != "eth_blockNumber" {
			t.Errorf("expected eth_blockNumber method, got %s", rpcReq.Method)
		}

		if rpcReq.JSONRPC != "2.0" {
			t.Errorf("expected jsonrpc 2.0, got %s", rpcReq.JSONRPC)
		}

		if rpcReq.ID != 2 {
			t.Errorf("expected id 2, got %d", rpcReq.ID)
		}

		// Write response
		response := RPCResponse{
			JSONRPC: "2.0",
			ID:      rpcReq.ID,
			Result:  json.RawMessage(`"0x1234567"`),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client that uses the mock server
	client := NewClient(server.URL)

	// Test GetBlockNumber
	blockNumber, err := client.GetBlockNumber()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Validate the block number
	if blockNumber != "0x1234567" {
		t.Errorf("expected block number 0x1234567, got %s", blockNumber)
	}
}

func TestGetBlockByNumber(t *testing.T) {
	t.Run("with full transactions", func(t *testing.T) {
		// Create a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate request
			if r.Method != http.MethodPost {
				t.Errorf("expected POST request, got %s", r.Method)
			}

			var rpcReq RPCRequest
			if err := json.NewDecoder(r.Body).Decode(&rpcReq); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}

			if rpcReq.Method != "eth_getBlockByNumber" {
				t.Errorf("expected eth_getBlockByNumber method, got %s", rpcReq.Method)
			}

			if len(rpcReq.Params) != 2 {
				t.Errorf("expected 2 params, got %d", len(rpcReq.Params))
			}

			// Check block number param
			blockNumber, ok := rpcReq.Params[0].(string)
			if !ok || blockNumber != "0x1234567" {
				t.Errorf("expected block number 0x1234567, got %v", rpcReq.Params[0])
			}

			// Check full transactions param
			fullTx, ok := rpcReq.Params[1].(bool)
			if !ok || !fullTx {
				t.Errorf("expected full transactions true, got %v", rpcReq.Params[1])
			}

			// Create a response with transaction objects (not just strings)
			mockResponseStr := `{
				"number": "0x1234567",
				"hash": "0xabcdef1234567890",
				"parentHash": "0x1234567890abcdef",
				"nonce": "0x123456",
				"timestamp": "0x60123456",
				"transactions": [
					{
						"hash": "0xtx1",
						"from": "0xaddr1",
						"to": "0xaddr2"
					},
					{
						"hash": "0xtx2",
						"from": "0xaddr3",
						"to": "0xaddr4"
					}
				]
			}`

			response := RPCResponse{
				JSONRPC: "2.0",
				ID:      rpcReq.ID,
				Result:  json.RawMessage(mockResponseStr),
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		// Create client that uses the mock server
		client := NewClient(server.URL)

		// Test GetBlockByNumber with full transactions
		block, err := client.GetBlockByNumber("0x1234567", true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Validate block data
		if block.Number != "0x1234567" {
			t.Errorf("expected block number 0x1234567, got %s", block.Number)
		}
		if block.Hash != "0xabcdef1234567890" {
			t.Errorf("expected block hash 0xabcdef1234567890, got %s", block.Hash)
		}
		if block.TransactionCount != 2 {
			t.Errorf("expected transaction count 2, got %d", block.TransactionCount)
		}
	})

	t.Run("with transaction hashes only", func(t *testing.T) {
		// Create a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate request
			if r.Method != http.MethodPost {
				t.Errorf("expected POST request, got %s", r.Method)
			}

			var rpcReq RPCRequest
			if err := json.NewDecoder(r.Body).Decode(&rpcReq); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}

			if rpcReq.Method != "eth_getBlockByNumber" {
				t.Errorf("expected eth_getBlockByNumber method, got %s", rpcReq.Method)
			}

			if len(rpcReq.Params) != 2 {
				t.Errorf("expected 2 params, got %d", len(rpcReq.Params))
			}

			// Check block number param
			blockNumber, ok := rpcReq.Params[0].(string)
			if !ok || blockNumber != "0x1234567" {
				t.Errorf("expected block number 0x1234567, got %v", rpcReq.Params[0])
			}

			// Check full transactions param
			fullTx, ok := rpcReq.Params[1].(bool)
			if !ok || fullTx {
				t.Errorf("expected full transactions false, got %v", rpcReq.Params[1])
			}

			// Create a response with transaction hashes (just strings)
			mockResponseStr := `{
				"number": "0x1234567",
				"hash": "0xabcdef1234567890",
				"parentHash": "0x1234567890abcdef",
				"nonce": "0x123456",
				"timestamp": "0x60123456",
				"transactions": ["0xtx1", "0xtx2"]
			}`

			response := RPCResponse{
				JSONRPC: "2.0",
				ID:      rpcReq.ID,
				Result:  json.RawMessage(mockResponseStr),
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		// Create client that uses the mock server
		client := NewClient(server.URL)

		// Test GetBlockByNumber with transaction hashes only
		block, err := client.GetBlockByNumber("0x1234567", false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Validate block data
		if block.Number != "0x1234567" {
			t.Errorf("expected block number 0x1234567, got %s", block.Number)
		}
		if block.Hash != "0xabcdef1234567890" {
			t.Errorf("expected block hash 0xabcdef1234567890, got %s", block.Hash)
		}
		if block.TransactionCount != 2 {
			t.Errorf("expected transaction count 2, got %d", block.TransactionCount)
		}
	})

	t.Run("error handling", func(t *testing.T) {
		// Create a mock HTTP server that returns an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := RPCResponse{
				JSONRPC: "2.0",
				ID:      2,
				Error: &RPCError{
					Code:    -32000,
					Message: "Invalid block number",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		// Create client that uses the mock server
		client := NewClient(server.URL)

		// Test GetBlockByNumber with an invalid block number
		_, err := client.GetBlockByNumber("0xinvalid", true)
		if err == nil {
			t.Errorf("expected error but got nil")
		}
	})
}
