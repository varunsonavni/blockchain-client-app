FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod separately to cache dependencies
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /blockchain-client ./cmd/client

# Use Google's distroless image as base
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Copy only the binary from the builder stage
COPY --from=builder /blockchain-client .

EXPOSE 8080

ENV BLOCKCHAIN_RPC_URL="https://polygon-rpc.com/"
ENV API_PORT=":8080"

ENTRYPOINT [ "./blockchain-client" ] 