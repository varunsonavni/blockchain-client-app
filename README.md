# Blockchain Client Application

A Go-based client application for interacting with the Polygon blockchain using JSON-RPC.

## Features

- Direct connection to Polygon RPC endpoint (https://polygon-rpc.com/)
- JSON-RPC POST endpoint supporting key blockchain operations
- Containerized with Docker for easy deployment
- AWS ECS Fargate deployment using Terraform

## Supported JSON-RPC Methods

### Get Block Number

```
POST /
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "method": "eth_blockNumber",
  "id": 2
}
```

Response:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": "0x134e82a"
}
```

### Get Block By Number

```
POST /
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "method": "eth_getBlockByNumber",
  "params": [
    "0x134e82a",
    true
  ],
  "id": 2
}
```

Response:
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "number": "0x134e82a",
    "hash": "0x...",
    "parentHash": "0x...",
    "nonce": "0x...",
    "timestamp": "0x...",
    "transactions": [...],
    "transactionCount": 123
  }
}
```

## Getting Started

### Prerequisites

- Go 1.24+
- Docker and Docker Compose (for containerized deployment)

### Installation

#### Local Development

1. Clone the repository:
   ```
   git clone <repository-url>
   cd blockchain-client-app
   ```

2. Build & Run the application:
   ```
   make start
   ```

This will build the container and start the service on port 8080.

## Testing

Run the test suite:

```
make test-api
```

## Deployment

### AWS ECS Fargate

The application includes Terraform configuration for deployment to AWS ECS Fargate in the `terraform` directory.

To deploy:

1. Navigate to the `terraform/dev` directory
2. Initialize & apply Terraform:
   ```
   cd vpc
   terraform init
   terraform plan
   terraform apply
   cd ..
   cd ecs
   terraform init
   terraform plan
   terraform apply
   ```


## Production Readiness Recommendations

To make this application production-ready, consider implementing the following improvements:

1. **Authentication/Authorization**:
   - Implement API keys or JWT-based authentication
   - Add rate limiting to prevent abuse

2. **Monitoring and Logging**:
   - Integrate with CloudWatch or similar services
   - Add structured logging with log levels
   - Set up alerting for critical failures

3. **Scalability**:
   - Implement caching for frequently requested blocks
   - Consider connection pooling to the RPC endpoint

4. **Resilience**:
   - Add circuit breakers for calls to the blockchain
   - Implement failover to alternative RPC providers
   - Add request timeouts and retry mechanisms

5. **Security**:
   - Regular security audits and dependency scans
   - Implement TLS/HTTPS
   - Container image vulnerability scanning

6. **CI/CD Pipeline**:
   - Add automated testing and deployment
   - Version tagging for Docker images
   - Blue/green deployment strategies

7. **Documentation**:
   - API documentation using Swagger/OpenAPI
   - Enhanced usage examples and error handling guides

8. **Additional Features**:
   - Support for more blockchain methods
   - Websocket subscription support for real-time updates
   - Enhanced error handling and reporting
