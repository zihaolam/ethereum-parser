# Parser Service

This is a simple Go-based HTTP server for handling Ethereum address subscriptions, fetching transactions, and retrieving the current block. The server is designed to use an Ethereum blockchain parser, allowing users to subscribe to Ethereum addresses and retrieve transactions for those addresses.

## Features

- Subscribe: Subscribe to Ethereum addresses to track transactions.
- Get Transactions: Retrieve a list of transactions for subscribed addresses.
- Get Current Block: Get the latest block number.
- Wildcard Route Logging: Logs requests to unhandled routes, providing a custom 404 response.

## Requirements

- Go 1.21+
- Make

## Getting Started

### Installation

Clone the repository:

```bash
git clone https://github.com/zihaolam/ethereum-parser.git
cd ethereum-parser
```

Install dependencies:

```bash
go mod tidy
```

### Build and Run

#### Build the Project

To build the project, run:

```bash
make build
```

This will compile the project and place the binary in the `bin/` directory as `parser`.

#### Run the Project

To run the server, use:

```bash
make run
```

This will build and run the server binary.

By default, the server listens on `localhost:8080`. You can change the address by modifying the `addr` flag in `cmd/parser/main.go`.

### API Endpoints

- Subscribe to an Address:

  ```bash
  POST /subscribe?address=<ethereum_address>
  ```

- Get Transactions for an Address:

  ```bash
  GET /transactions?address=<ethereum_address>
  ```

- Get Current Block:

  ```bash
  GET /current_block
  ```

- Wildcard Unhandled Route: Any unhandled routes are logged and return a `404 - Not Found` response.

### Logging

The server uses a custom logging middleware to log request details. Each request logs the method, path, and processing time.

### Testing

The Makefile includes several test commands for running tests:

Run all tests:

```bash
make test-all
```

Run memory database tests:

```bash
make test-memorydb
```

Run parser tests:

```bash
make test-parser
```

### Cleanup

To remove the generated binary and clean up:

```bash
make clean
```

## Makefile Commands Summary

- build: Builds the project binary to `bin/parser`.
- run: Builds and runs the project.
- test-all: Runs all tests.
- test-memorydb: Runs tests for `memorydb`.
- test-parser: Runs tests for the `parser` package.
- clean: Cleans up the binary and other generated files.
