package mcp

import (
	"encoding/json"

	"golang.org/x/exp/slog"
)

type Request struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      json.RawMessage  `json:"id,omitempty"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params,omitempty"`
}

func (r *Request) Log() {
	slog.Debug("[protocol] MCP Request", "jsonrpc", r.JSONRPC, "id", r.ID, "method", r.Method, "params", r.Params)
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *RespError      `json:"error,omitempty"`
}

func (r *Response) Log() {
	slog.Debug("[protocol] MCP Response", "jsonrpc", r.JSONRPC, "id", r.ID, "result", r.Result, "error", r.Error)
}

type RespError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RespError) Log() {
	slog.Error("[protocol] MCP Error", "code", e.Code, "message", e.Message)
}

type InitializeParams struct {
	ClientInfo struct {
		Name    string `json:"name"`
		Version string `json:"version,omitempty"`
	} `json:"clientInfo"`
}

type InitializeResult struct {
	ServerInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"serverInfo"`
	Capabilities map[string]any `json:"capabilities"`
}

type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

type ToolCallParams struct {
	Name string         `json:"name"`
	Args map[string]any `json:"arguments"`
}

type ToolCallResult struct {
	Content []ToolContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
