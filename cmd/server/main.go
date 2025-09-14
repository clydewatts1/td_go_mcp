package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"td_go_mcp/internal/server"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	mux := http.NewServeMux()
	server.RegisterRoutes(mux)

	log.Printf("Server listening on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Errorf("server error: %w", err))
	}
}
