package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/zihaolam/ethereum-parser/internal/logging"
	"github.com/zihaolam/ethereum-parser/internal/parser"
)

type Api struct {
	parser *parser.Parser
	logger logging.Logger
}

func New(parser *parser.Parser, logger logging.Logger) *Api {
	return &Api{
		parser: parser,
		logger: logger,
	}
}

func (api *Api) Start(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/subscribe", api.loggingMiddleware(api.handleSubscribe()))
	mux.HandleFunc("/transactions", api.loggingMiddleware(api.handleGetTransactions()))
	mux.HandleFunc("/current_block", api.loggingMiddleware(api.handleGetCurrentBlock()))
	mux.HandleFunc("/scan", api.loggingMiddleware(api.handleScanBlock()))
	mux.HandleFunc("/", api.loggingMiddleware(api.handleWildcard()))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

func (api *Api) loggingMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		api.logger.Printf("Started %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		api.logger.Printf("Completed %s %s in %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func (api *Api) handleWildcard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (api *Api) handleSubscribe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address required", http.StatusBadRequest)
			return
		}

		if subscribed := api.parser.Subscribe(address); !subscribed {
			json.NewEncoder(w).
				Encode(map[string]interface{}{"data": false, "message": "Already subscribed" + address})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).
			Encode(map[string]interface{}{"data": true, "message": "Subscribed to address " + address})
		return
	}
}

func (api *Api) handleGetTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address required", http.StatusBadRequest)
			return
		}

		transactions := api.parser.GetTransactions(address)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transactions)
	}
}

func (api *Api) handleGetCurrentBlock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		block := api.parser.GetCurrentBlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"current_block": block})
	}
}

func (api *Api) handleScanBlock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blocknumberQuery := r.URL.Query().Get("blocknumber")
		blocknumber, err := strconv.Atoi(blocknumberQuery)
		if err != nil {
			http.Error(w, "Failed to parse block number", http.StatusBadRequest)
			return
		}
		block, err := api.parser.ScanBlock(r.Context(), blocknumber)
		if err != nil {
			http.Error(w, "Failed to scan block", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(block)
	}
}
