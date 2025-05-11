package main

import (
	"flag"
	"log"
	"os"

	"blockchain-client/pkg/api"
	"blockchain-client/pkg/blockchain"
)

func main() {
	rpcURL := flag.String("rpc", blockchain.PolygonRPC, "Blockchain RPC URL")
	port := flag.String("port", ":8080", "API server port")
	flag.Parse()

	if envRPC := os.Getenv("BLOCKCHAIN_RPC_URL"); envRPC != "" {
		*rpcURL = envRPC
	}

	if envPort := os.Getenv("API_PORT"); envPort != "" {
		*port = envPort
	}

	server := api.NewServer(*rpcURL)

	log.Printf("Blockchain client connecting to %s", *rpcURL)
	log.Printf("Starting API server on %s", *port)

	if err := server.Start(*port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
