package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"td_go_mcp/internal/tools"

	"golang.org/x/exp/slog"
)

type ServerInfo struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
	Tools     int       `json:"tools_loaded"`
}

func RegisterRoutes(mux *http.ServeMux) {
	slog.Debug("[server] Registering HTTP routes")
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/healthz", handleHealth)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	slog.Info("[server] Handling root endpoint request")
	loadedTools, err := tools.LoadToolsFromDirectory("tools")
	toolCount := 0
	status := "ok"
	if err != nil {
		status = "warning: " + err.Error()
		slog.Error("[server] Error loading tools from directory", "error", err)
	} else {
		toolCount = len(loadedTools)
		slog.Debug("[server] Loaded tools for root endpoint", "count", toolCount)
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
		slog.Error("[server] Failed to encode root endpoint response", "error", err)
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	} else {
		slog.Info("[server] Root endpoint response sent", "info", info)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	slog.Debug("[server] Handling healthz endpoint request")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		slog.Error("[server] Failed to write healthz response", "error", err)
	} else {
		slog.Debug("[server] Healthz response sent")
	}
}
