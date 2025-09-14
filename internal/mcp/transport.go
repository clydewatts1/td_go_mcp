package mcp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Transport struct {
	r *bufio.Reader
	w *bufio.Writer
}

func NewTransport(r io.Reader, w io.Writer) *Transport {
	return &Transport{r: bufio.NewReader(r), w: bufio.NewWriter(w)}
}

// Read reads a single JSON-RPC message framed with Content-Length headers (LSP-style).
func (t *Transport) Read() ([]byte, error) {
	// Read headers
	contentLength := -1
	for {
		line, err := t.r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" { // end of headers
			break
		}
		if cl, ok := parseContentLength(line); ok {
			contentLength = cl
		}
	}
	if contentLength < 0 {
		return nil, errors.New("missing Content-Length header")
	}
	buf := make([]byte, contentLength)
	if _, err := io.ReadFull(t.r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (t *Transport) Write(payload []byte) error {
	if _, err := t.w.WriteString(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(payload))); err != nil {
		return err
	}
	if _, err := t.w.Write(payload); err != nil {
		return err
	}
	return t.w.Flush()
}

func parseContentLength(line string) (int, bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return 0, false
	}
	if !strings.EqualFold(strings.TrimSpace(parts[0]), "Content-Length") {
		return 0, false
	}
	n, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, false
	}
	return n, true
}

// Helper to build framed message (useful for testing)
func Frame(payload []byte) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "Content-Length: %d\r\n\r\n", len(payload))
	_, _ = b.Write(payload)
	out := b.Bytes()
	// Return a copy to avoid mutations if caller holds reference
	cp := make([]byte, len(out))
	copy(cp, out)
	return cp
}
