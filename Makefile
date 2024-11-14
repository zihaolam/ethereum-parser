BINARY_NAME := bin/parser

build:
	go build -o $(BINARY_NAME) cmd/parser/main.go

test-all:
	go test ./...

test-memorydb:
	go test -v ./internal/datastore/memorydb/memorydb_test.go

test-parser:
	go test -v ./.../internal/parser	

clean:
	go clean
	rm -f $(BINARY_NAME)
