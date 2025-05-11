package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"blockchain-client/pkg/blockchain"
)

// BlockchainClient interface for blockchain operations
type BlockchainClient interface {
	GetBlockNumber() (string, error)
	GetBlockByNumber(blockNumber string, fullTransactions bool) (*blockchain.Block, error)
}

// Server represents the API server
type Server struct {
	client BlockchainClient
}

// NewServer creates a new API server
func NewServer(rpcURL string) *Server {
	return &Server{
		client: blockchain.NewClient(rpcURL),
	}
}

// BlockNumberResponse represents the response for block number endpoint
type BlockNumberResponse struct {
	BlockNumber string `json:"blockNumber"`
}

// BlockResponse represents the response for block details endpoint
type BlockResponse struct {
	Block *blockchain.Block `json:"block"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
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

// writeJSONResponse writes a JSON response
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// HandleGetBlockNumber handles the /blocks/latest endpoint
func (s *Server) HandleGetBlockNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONResponse(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	blockNumber, err := s.client.GetBlockNumber()
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, BlockNumberResponse{BlockNumber: blockNumber})
}

// HandleGetBlockByNumber handles the /blocks/{blockNumber} endpoint
func (s *Server) HandleGetBlockByNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONResponse(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	blockNumber := r.URL.Query().Get("number")
	if blockNumber == "" {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Error: "block number is required"})
		return
	}

	fullTx := r.URL.Query().Get("full") == "true"

	block, err := s.client.GetBlockByNumber(blockNumber, fullTx)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSONResponse(w, http.StatusOK, BlockResponse{Block: block})
}

// HandleJSONRPC handles JSON-RPC requests directly
func (s *Server) HandleJSONRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONResponse(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Error: "failed to read request body"})
		return
	}
	defer r.Body.Close()

	// Parse JSON-RPC request
	var request RPCRequest
	if err := json.Unmarshal(body, &request); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON-RPC request"})
		return
	}

	// Check JSON-RPC version
	if request.JSONRPC != "2.0" {
		writeJSONResponse(w, http.StatusBadRequest, RPCResponse{
			JSONRPC: "2.0",
			Error: &RPCError{
				Code:    -32600,
				Message: "invalid JSON-RPC version, expected 2.0",
			},
			ID: request.ID,
		})
		return
	}

	// Process request based on method
	var result interface{}
	var rpcError *RPCError

	switch request.Method {
	case "eth_blockNumber":
		blockNumber, err := s.client.GetBlockNumber()
		if err != nil {
			rpcError = &RPCError{
				Code:    -32603,
				Message: err.Error(),
			}
		} else {
			result = blockNumber
		}

	case "eth_getBlockByNumber":
		if len(request.Params) < 2 {
			rpcError = &RPCError{
				Code:    -32602,
				Message: "invalid params for eth_getBlockByNumber",
			}
			break
		}

		// Get block number from params
		blockNumberParam, ok := request.Params[0].(string)
		if !ok {
			rpcError = &RPCError{
				Code:    -32602,
				Message: "invalid block number parameter",
			}
			break
		}

		// Get full transactions from params
		fullTransactions, ok := request.Params[1].(bool)
		if !ok {
			rpcError = &RPCError{
				Code:    -32602,
				Message: "invalid full transactions parameter",
			}
			break
		}

		// Get block
		block, err := s.client.GetBlockByNumber(blockNumberParam, fullTransactions)
		if err != nil {
			rpcError = &RPCError{
				Code:    -32603,
				Message: err.Error(),
			}
		} else {
			result = block
		}

	default:
		rpcError = &RPCError{
			Code:    -32601,
			Message: "method not found",
		}
	}

	// Prepare response
	response := RPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Error:   rpcError,
	}

	if result != nil && rpcError == nil {
		resultBytes, err := json.Marshal(result)
		if err != nil {
			response.Error = &RPCError{
				Code:    -32603,
				Message: "failed to marshal result",
			}
		} else {
			response.Result = resultBytes
		}
	}

	// Send response
	writeJSONResponse(w, http.StatusOK, response)
}

// SetupRoutes sets up the API routes
func (s *Server) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Original REST endpoints
	mux.HandleFunc("/api/blocks/latest", s.HandleGetBlockNumber)
	mux.HandleFunc("/api/blocks", s.HandleGetBlockByNumber)

	// New JSON-RPC endpoint
	mux.HandleFunc("/", s.HandleJSONRPC)

	return mux
}

// Start starts the API server
func (s *Server) Start(addr string) error {
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("Starting API server on %s", addr)
	return http.ListenAndServe(addr, s.SetupRoutes())
}
