package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"td_go_mcp/internal/tools"
)

type ServerInfo struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
	Tools     int       `json:"tools_loaded"`
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/healthz", handleHealth)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Load tools count
	loadedTools, err := tools.LoadToolsFromDirectory("tools")
	toolCount := 0
	status := "ok"
	if err != nil {
		status = "warning: " + err.Error()
	} else {
		toolCount = len(loadedTools)
	}

	info := ServerInfo{
		Name:      "td-go-mcp",
		Version:   "0.2.0",
		Timestamp: time.Now().UTC(),
		Status:    status,
		Tools:     toolCount,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
