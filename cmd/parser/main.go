package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zihaolam/ethereum-parser/internal/api"
	"github.com/zihaolam/ethereum-parser/internal/parser"
)

func getEndpoint(isTestnet bool) string {
	if isTestnet {
		return "https://cloudflare-eth.com"
	}
	return "https://ethereum-rpc.publicnode.com"
}

func main() {
	// Define command-line flags
	addr := flag.String(
		"addr",
		":8080",
		"Address to start the server on, e.g., ':8080' or 'localhost:8080'",
	)
	testnet := flag.Bool("testnet", false, "Use testnet endpoint")

	// Default to 0 means start parsing from the latest block
	initialBlockNumber := flag.Int("initial-block", 0, "Initial block number to start parsing from")
	scanInterval := flag.Int("scan-interval", 10, "Interval in seconds to scan for new blocks")

	flag.Parse()

	// Use default logger for now
	logger := log.Default()

	endpoint := getEndpoint(*testnet)

	p := parser.New(logger, endpoint, *initialBlockNumber)

	// Initialize the API with the parser
	api := api.New(p, logger)

	// Set up a context to handle server shutdown gracefully
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture interrupt and termination signals to gracefully shut down the server
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		logger.Println("Shutting down server...")
		cancel()
	}()

	go func() {
		// Start the scanner
		interval := time.Duration(*scanInterval) * time.Second
		logger.Println("Starting scanner on " + endpoint)
		p.StartScan(ctx, interval)
	}()

	// Start the server
	logger.Printf("Starting server on %s...\n", *addr)
	if err := api.Start(ctx, *addr); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not start server: %v\n", err)
	}

	// Allow time for graceful shutdown
	logger.Println("Server stopped")
	time.Sleep(1 * time.Second)
}
