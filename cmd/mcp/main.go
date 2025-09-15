package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"td_go_mcp/internal/mcp"
)

func main() {
	// Setup logging
	logDir := "logging"
	os.MkdirAll(logDir, 0755)

	currentLogPath := filepath.Join(logDir, "current.log")
	historyLogPath := filepath.Join(logDir, "history-"+time.Now().Format("20060102")+".log")

	currentLog, err := os.OpenFile(currentLogPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open current log: %v\n", err)
		os.Exit(1)
	}
	defer currentLog.Close()

	historyLog, err := os.OpenFile(historyLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open history log: %v\n", err)
		os.Exit(1)
	}
	defer historyLog.Close()

	log.SetOutput(currentLog)

	logRequest := func(label string, data any) {
		b, _ := json.MarshalIndent(data, "", "  ")
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		entry := fmt.Sprintf("[%s] %s: %s\n", timestamp, label, string(b))
		currentLog.WriteString(entry)
		historyLog.WriteString(entry)
	}

	in := bufio.NewReader(os.Stdin)
	for {
		body, err := mcp.ReadFrame(in)
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			_ = mcp.WriteJSON(os.Stdout, mcp.Response{
				JSONRPC: "2.0",
				Error:   &mcp.RespError{Code: -32700, Message: "read error: " + err.Error()},
			})
			return
		}

		var req mcp.Request
		if err := json.Unmarshal(body, &req); err != nil {
			logRequest("parse error", string(body))
			_ = mcp.WriteJSON(os.Stdout, mcp.Response{
				JSONRPC: "2.0",
				Error:   &mcp.RespError{Code: -32700, Message: "parse error: " + err.Error()},
			})
			continue
		}

		logRequest("request", req)

		var resp any
		switch req.Method {
		case "initialize":
			resp = handleInitialize(req)
		case "tools/list":
			resp = handleToolsList(req)
		case "tools/call":
			resp = handleToolsCall(req)
		default:
			resp = mcp.Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &mcp.RespError{Code: -32601, Message: "method not found"},
			}
			_ = mcp.WriteJSON(os.Stdout, resp)
		}

		if resp != nil {
			logRequest("response", resp)
		}
	}
}

func handleInitialize(req mcp.Request) any {
	var params mcp.InitializeParams
	if req.Params != nil {
		_ = json.Unmarshal(*req.Params, &params)
	}

	res := mcp.InitializeResult{
		Capabilities: map[string]any{
			"tools": map[string]any{},
		},
	}
	res.ServerInfo.Name = "td-go-mcp"
	res.ServerInfo.Version = "0.1.0"

	response := mcp.Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  res,
	}
	_ = mcp.WriteJSON(os.Stdout, response)
	return response
}

func toolDefs() []mcp.Tool {
	return []mcp.Tool{
		{
			Name:        "ping",
			Description: "Echo back the input text prefixed with 'pong: '",
			InputSchema: map[string]any{
				"type":       "object",
				"required":   []string{"text"},
				"properties": map[string]any{"text": map[string]any{"type": "string"}},
			},
		},
		{
			Name:        "time",
			Description: "Return current UTC time in RFC3339 format",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "upper",
			Description: "Uppercase the provided text",
			InputSchema: map[string]any{
				"type":       "object",
				"required":   []string{"text"},
				"properties": map[string]any{"text": map[string]any{"type": "string"}},
			},
		},
		{
			Name:        "sum",
			Description: "Sum an array of numbers",
			InputSchema: map[string]any{
				"type":     "object",
				"required": []string{"numbers"},
				"properties": map[string]any{
					"numbers": map[string]any{
						"type":  "array",
						"items": map[string]any{"type": "number"},
					},
				},
			},
		},
		{
			Name:        "uuid",
			Description: "Generate a UUID v4",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
	}
}

func handleToolsList(req mcp.Request) any {
	response := mcp.Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  mcp.ToolsListResult{Tools: toolDefs()},
	}
	_ = mcp.WriteJSON(os.Stdout, response)
	return response
}

func handleToolsCall(req mcp.Request) any {
	var params mcp.ToolsCallParams
	if req.Params != nil {
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			_ = mcp.WriteJSON(os.Stdout, mcp.Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &mcp.RespError{Code: -32602, Message: "invalid params"},
			})
			return
		}
	}

	writeText := func(text string) any {
		response := mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  mcp.ToolsCallResult{Content: []mcp.ContentItem{{Type: "text", Text: text}}},
		}
		_ = mcp.WriteJSON(os.Stdout, response)
		return response
	}

	switch params.Name {
	case "ping":
		text := getString(params.Arguments, "text")
		return writeText("pong: " + text)
	case "time":
		return writeText(time.Now().UTC().Format(time.RFC3339))
	case "upper":
		text := getString(params.Arguments, "text")
		return writeText(strings.ToUpper(text))
	case "sum":
		total := sumNumbers(params.Arguments["numbers"])
		return writeText(fmt.Sprintf("sum: %g", total))
	case "uuid":
		return writeText(uuidV4())
	default:
		response := mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &mcp.RespError{Code: -32601, Message: "unknown tool: " + params.Name},
		}
		_ = mcp.WriteJSON(os.Stdout, response)
		return response
	}
}

func getString(m map[string]any, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func sumNumbers(v any) float64 {
	arr, ok := v.([]any)
	if !ok {
		return 0
	}
	var total float64
	for _, x := range arr {
		switch n := x.(type) {
		case float64:
			total += n
		case float32:
			total += float64(n)
		case int:
			total += float64(n)
		case int64:
			total += float64(n)
		case json.Number:
			if f, err := n.Float64(); err == nil {
				total += f
			}
		}
	}
	return total
}

func uuidV4() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 10
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16],
	)
}
