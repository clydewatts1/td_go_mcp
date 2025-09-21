package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

type Request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type Response struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./cmd/td_go_mcp_test_client <mcp_request.json>")
		os.Exit(1)
	}

	requestFile := os.Args[1]
	slog.Info("Reading MCP request from file", "file", requestFile)

	file, err := os.Open(requestFile)
	if err != nil {
		slog.Error("Failed to open request file", "error", err)
		os.Exit(1)
	}
	defer file.Close()

	requestBytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read request file", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting MCP server process")
	cmd := exec.Command("go", "run", "./cmd/mcp")
	cmd.Stderr = os.Stderr // Redirect MCP server's stderr to our stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		slog.Error("Failed to get stdin pipe", "error", err)
		os.Exit(1)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Failed to get stdout pipe", "error", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		slog.Error("Failed to start MCP server process", "error", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	// Write to the server's stdin synchronously first
	slog.Info("Sending request to MCP server")

	// Parse the JSON to ensure it's valid and compact it
	var jsonData interface{}
	if err := json.Unmarshal(requestBytes, &jsonData); err != nil {
		slog.Error("Failed to parse JSON request", "error", err)
		os.Exit(1)
	}

	// Marshal back to compact JSON
	compactJSON, err := json.Marshal(jsonData)
	if err != nil {
		slog.Error("Failed to marshal compact JSON", "error", err)
		os.Exit(1)
	}

	requestStr := string(compactJSON) + "\n"

	// Write to a temp file for debugging
	tempFile := "debug_request.txt"
	if err := os.WriteFile(tempFile, []byte(requestStr), 0644); err != nil {
		slog.Error("Failed to write debug file", "error", err)
	} else {
		slog.Info("Wrote request to debug file", "file", tempFile)
	}

	// Use bufio.Writer to ensure proper flushing
	writer := bufio.NewWriter(stdin)
	if _, err := writer.WriteString(requestStr); err != nil {
		slog.Error("Failed to write request to stdin", "error", err)
		os.Exit(1)
	}
	if err := writer.Flush(); err != nil {
		slog.Error("Failed to flush stdin", "error", err)
		os.Exit(1)
	}

	// Small delay to ensure data is fully transmitted
	time.Sleep(200 * time.Millisecond)

	// Close stdin to signal EOF to the server
	stdin.Close()
	slog.Info("Finished writing to stdin")

	wg.Add(1)
	// Goroutine to read from the server's stdout
	go func() {
		defer wg.Done()
		slog.Info("Reading response from MCP server")
		reader := bufio.NewReader(stdout)

		// For stdio transport, read line by line
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					slog.Error("Failed to read line from stdout", "error", err)
				}
				break
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			slog.Info("Received response line", "line", line)

			var response Response
			if err := json.Unmarshal([]byte(line), &response); err != nil {
				slog.Error("Failed to parse JSON response", "error", err)
				// Still print the raw line for debugging
				fmt.Println(line)
				continue
			}

			prettyJSON, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(prettyJSON))
			break // Exit after receiving one response
		}
	}()

	wg.Add(1)
	// Goroutine to read from the server's stdout
	go func() {
		defer wg.Done()
		slog.Info("Reading response from MCP server")
		reader := bufio.NewReader(stdout)

		// For stdio transport, read line by line
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					slog.Error("Failed to read line from stdout", "error", err)
				}
				break
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			slog.Info("Received response line", "line", line)

			var response Response
			if err := json.Unmarshal([]byte(line), &response); err != nil {
				slog.Error("Failed to parse JSON response", "error", err)
				// Still print the raw line for debugging
				fmt.Println(line)
				continue
			}

			prettyJSON, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(prettyJSON))
			break // Exit after receiving one response
		}
	}()

	wg.Wait()
	cmd.Wait()
	slog.Info("Test client finished")
}
