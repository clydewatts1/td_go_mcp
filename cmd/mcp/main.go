package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"td_go_mcp/internal/mcp"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	for {
		body, err := mcp.ReadFrame(in)
		if err != nil {
			// Exit cleanly on EOF/pipe close
			if err.Error() == "EOF" {
				return
			}
			// Best-effort error response if we can't parse a request id
			_ = mcp.WriteJSON(os.Stdout, mcp.Response{
				JSONRPC: "2.0",
				Error:   &mcp.RespError{Code: -32700, Message: "read error: " + err.Error()},
			})
			return
		}

		var req mcp.Request
		if err := json.Unmarshal(body, &req); err != nil {
			_ = mcp.WriteJSON(os.Stdout, mcp.Response{
				JSONRPC: "2.0",
				Error:   &mcp.RespError{Code: -32700, Message: "parse error: " + err.Error()},
			})
			continue
		}

		switch req.Method {
		case "initialize":
			handleInitialize(req)
		case "tools/list":
			handleToolsList(req)
		case "tools/call":
			handleToolsCall(req)
		default:
			_ = mcp.WriteJSON(os.Stdout, mcp.Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &mcp.RespError{Code: -32601, Message: "method not found"},
			})
		}
	}
}

func handleInitialize(req mcp.Request) {
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

	_ = mcp.WriteJSON(os.Stdout, mcp.Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  res,
	})
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

func handleToolsList(req mcp.Request) {
	_ = mcp.WriteJSON(os.Stdout, mcp.Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  mcp.ToolsListResult{Tools: toolDefs()},
	})
}

func handleToolsCall(req mcp.Request) {
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

	writeText := func(text string) {
		_ = mcp.WriteJSON(os.Stdout, mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  mcp.ToolsCallResult{Content: []mcp.ContentItem{{Type: "text", Text: text}}},
		})
	}

	switch params.Name {
	case "ping":
		text := getString(params.Arguments, "text")
		writeText("pong: " + text)
	case "time":
		writeText(time.Now().UTC().Format(time.RFC3339))
	case "upper":
		text := getString(params.Arguments, "text")
		writeText(strings.ToUpper(text))
	case "sum":
		total := sumNumbers(params.Arguments["numbers"])
		writeText(fmt.Sprintf("sum: %g", total))
	case "uuid":
		writeText(uuidV4())
	default:
		_ = mcp.WriteJSON(os.Stdout, mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &mcp.RespError{Code: -32601, Message: "unknown tool: " + params.Name},
		})
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
