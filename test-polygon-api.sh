#!/bin/bash

echo "Testing Polygon-compatible RPC API endpoints..."

# Check if container is running
if ! docker ps | grep blockchain-client > /dev/null; then
  echo "Error: blockchain-client container is not running."
  echo "Start it using: make start"
  exit 1
fi

# Wait for service to be ready
echo "Waiting for service to be ready..."
sleep 3

# Step 1: Get the latest block number using eth_blockNumber
echo -e "\n========== STEP 1: Testing eth_blockNumber ==========\n"

BLOCK_NUMBER_RESPONSE=$(curl -s -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
  "jsonrpc": "2.0",
  "method": "eth_blockNumber",
  "id": 2
}')

echo "Response from eth_blockNumber:"
echo "$BLOCK_NUMBER_RESPONSE" | jq .

# Extract block number from the response
BLOCK_NUM=$(echo "$BLOCK_NUMBER_RESPONSE" | jq -r .result)

if [ -z "$BLOCK_NUM" ] || [ "$BLOCK_NUM" == "null" ]; then
  echo "Error: Could not get a valid block number from the response."
  exit 1
fi

echo -e "\nSuccessfully got block number: $BLOCK_NUM"

# Step 2: Get block details using eth_getBlockByNumber with the block number from Step 1
echo -e "\n========== STEP 2: Testing eth_getBlockByNumber with full=true ==========\n"

BLOCK_DETAILS_RESPONSE=$(curl -s -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d "{
  \"jsonrpc\": \"2.0\",
  \"method\": \"eth_getBlockByNumber\",
  \"params\": [
    \"$BLOCK_NUM\",
    true
  ],
  \"id\": 2
}")

echo "Response from eth_getBlockByNumber (full=true):"
echo "$BLOCK_DETAILS_RESPONSE" | jq .

# Check if we got a valid response with block data
BLOCK_HASH=$(echo "$BLOCK_DETAILS_RESPONSE" | jq -r '.result.hash')
if [ -z "$BLOCK_HASH" ] || [ "$BLOCK_HASH" == "null" ]; then
  echo "Error: Could not get a valid block hash from the response."
  exit 1
fi

echo -e "\nSuccessfully got block with hash: $BLOCK_HASH"

# Step 3: Get block with transaction hashes only (full=false)
echo -e "\n========== STEP 3: Testing eth_getBlockByNumber with full=false ==========\n"

BLOCK_HASH_RESPONSE=$(curl -s -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d "{
  \"jsonrpc\": \"2.0\",
  \"method\": \"eth_getBlockByNumber\",
  \"params\": [
    \"$BLOCK_NUM\",
    false
  ],
  \"id\": 2
}")

echo "Response from eth_getBlockByNumber (full=false):"
echo "$BLOCK_HASH_RESPONSE" | jq .

# Check if we got a valid response with block data
BLOCK_HASH2=$(echo "$BLOCK_HASH_RESPONSE" | jq -r '.result.hash')
if [ -z "$BLOCK_HASH2" ] || [ "$BLOCK_HASH2" == "null" ]; then
  echo "Error: Could not get a valid block hash from the response."
  exit 1
fi

echo -e "\nSuccessfully got block with hash: $BLOCK_HASH2"