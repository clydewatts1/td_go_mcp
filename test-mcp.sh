#!/bin/bash

# Comprehensive test script for MCP server
echo "=== Testing MCP Server ==="

# Test 1: Initialize
echo "1. Testing initialize..."
req='{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"test"}}}'
bytes=`echo -n "$req" | wc -c`
response=$(printf "Content-Length: $bytes\r\n\r\n$req" | timeout 5s go run ./cmd/mcp)
if echo "$response" | grep -q "td-go-mcp"; then
    echo "✓ Initialize test passed"
else
    echo "✗ Initialize test failed"
    echo "Response: $response"
fi

# Test 2: Tools list
echo "2. Testing tools/list..."
req='{"jsonrpc":"2.0","id":2,"method":"tools/list"}'
bytes=`echo -n "$req" | wc -c`  
response=$(printf "Content-Length: $bytes\r\n\r\n$req" | timeout 5s go run ./cmd/mcp)
if echo "$response" | grep -q '"name":"ping"' && echo "$response" | grep -q '"name":"uuid"'; then
    echo "✓ Tools list test passed"
else
    echo "✗ Tools list test failed"
    echo "Response: $response"
fi

# Test 3: Ping tool
echo "3. Testing ping tool..."
req='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"ping","arguments":{"text":"test"}}}'
bytes=`echo -n "$req" | wc -c`
response=$(printf "Content-Length: $bytes\r\n\r\n$req" | timeout 5s go run ./cmd/mcp)
if echo "$response" | grep -q "pong: test"; then
    echo "✓ Ping tool test passed"
else
    echo "✗ Ping tool test failed"
    echo "Response: $response"
fi

# Test 4: Time tool
echo "4. Testing time tool..."
req='{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"time","arguments":{}}}'
bytes=`echo -n "$req" | wc -c`
response=$(printf "Content-Length: $bytes\r\n\r\n$req" | timeout 5s go run ./cmd/mcp)
if echo "$response" | grep -qE "20[0-9]{2}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z"; then
    echo "✓ Time tool test passed"
else
    echo "✗ Time tool test failed" 
    echo "Response: $response"
fi

# Test 5: Sum tool
echo "5. Testing sum tool..."
req='{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"sum","arguments":{"numbers":[1,2,3]}}}'
bytes=`echo -n "$req" | wc -c`
response=$(printf "Content-Length: $bytes\r\n\r\n$req" | timeout 5s go run ./cmd/mcp)
if echo "$response" | grep -q "sum: 6"; then
    echo "✓ Sum tool test passed"
else
    echo "✗ Sum tool test failed"
    echo "Response: $response"
fi

# Test 6: UUID tool
echo "6. Testing uuid tool..."
req='{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"uuid","arguments":{}}}'
bytes=`echo -n "$req" | wc -c`
response=$(printf "Content-Length: $bytes\r\n\r\n$req" | timeout 5s go run ./cmd/mcp)
if echo "$response" | grep -qE "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"; then
    echo "✓ UUID tool test passed"
else
    echo "✗ UUID tool test failed"
    echo "Response: $response"
fi

echo "=== All tests completed ==="