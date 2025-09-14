# Make targets are optional; on Windows use PowerShell equivalents.

.PHONY: tidy build test run

tidy:
	go mod tidy

build:
	go build ./...

test:
	go test ./...

run:
	set PORT=8080 && go run ./cmd/server
