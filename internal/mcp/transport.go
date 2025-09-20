package mcp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/exp/slog"
)

type Transport struct {
	r *bufio.Reader
	w *bufio.Writer
}

func NewTransport(r io.Reader, w io.Writer) *Transport {
	slog.Debug("[transport] Creating new Transport")
	return &Transport{r: bufio.NewReader(r), w: bufio.NewWriter(w)}
}

// Read reads a single JSON-RPC message framed with Content-Length headers (LSP-style).
func (t *Transport) Read() ([]byte, error) {
	slog.Debug("[transport] Reading message from transport")
	contentLength := -1
	for {
		line, err := t.r.ReadString('\n')
		if err != nil {
			slog.Error("[transport] Error reading header line", "error", err)
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" { // end of headers
			break
		}
		if cl, ok := parseContentLength(line); ok {
			contentLength = cl
			slog.Debug("[transport] Parsed Content-Length header", "length", cl)
		}
	}
	if contentLength < 0 {
		slog.Error("[transport] Missing Content-Length header")
		return nil, errors.New("missing Content-Length header")
	}
	buf := make([]byte, contentLength)
	if _, err := io.ReadFull(t.r, buf); err != nil {
		slog.Error("[transport] Error reading message body", "error", err)
		return nil, err
	}
	slog.Debug("[transport] Successfully read message", "length", contentLength)
	return buf, nil
}

func (t *Transport) Write(payload []byte) error {
	slog.Debug("[transport] Writing message to transport", "length", len(payload))
	if _, err := t.w.WriteString(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(payload))); err != nil {
		slog.Error("[transport] Error writing Content-Length header", "error", err)
		return err
	}
	if _, err := t.w.Write(payload); err != nil {
		slog.Error("[transport] Error writing payload", "error", err)
		return err
	}
	err := t.w.Flush()
	if err != nil {
		slog.Error("[transport] Error flushing writer", "error", err)
	} else {
		slog.Debug("[transport] Successfully wrote message")
	}
	return err
}

func parseContentLength(line string) (int, bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		slog.Debug("[transport] Invalid header line for Content-Length", "line", line)
		return 0, false
	}
	if !strings.EqualFold(strings.TrimSpace(parts[0]), "Content-Length") {
		return 0, false
	}
	n, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		slog.Error("[transport] Error parsing Content-Length value", "error", err, "line", line)
		return 0, false
	}
	return n, true
}

// Helper to build framed message (useful for testing)
func Frame(payload []byte) []byte {
	slog.Debug("[transport] Framing payload for test", "length", len(payload))
	var b bytes.Buffer
	fmt.Fprintf(&b, "Content-Length: %d\r\n\r\n", len(payload))
	_, _ = b.Write(payload)
	out := b.Bytes()
	// Return a copy to avoid mutations if caller holds reference
	cp := make([]byte, len(out))
	copy(cp, out)
	return cp
}
