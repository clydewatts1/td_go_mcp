package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/glebarez/sqlite"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Initialize database connection (using SQLite for development)
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create sample tables and data
	if err := initDatabase(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create a new MCP server
	mcpServer := server.NewMCPServer("td-go-mcp-fresh", "1.0.0")

	// Add database tools to the server
	addDatabaseTools(mcpServer, db)

	// Start the stdio server
	log.Println("Starting fresh MCP server...")
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// initDatabase creates sample tables and inserts test data
func initDatabase(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE,
			status TEXT DEFAULT 'active'
		)
	`)
	if err != nil {
		return err
	}

	// Create sessions table
	_, err = db.Exec(`
		CREATE TABLE sessions (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			start_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			end_time DATETIME,
			status TEXT DEFAULT 'active',
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	// Insert sample data
	_, err = db.Exec(`
		INSERT INTO users (name, email, status) VALUES
		('Alice Johnson', 'alice@example.com', 'active'),
		('Bob Smith', 'bob@example.com', 'active'),
		('Charlie Brown', 'charlie@example.com', 'inactive')
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO sessions (user_id, status) VALUES
		(1, 'active'),
		(2, 'active'),
		(1, 'ended')
	`)
	if err != nil {
		return err
	}

	return nil
}

// addDatabaseTools registers the database operation tools with the MCP server
func addDatabaseTools(mcpServer *server.MCPServer, db *sql.DB) {
	// Count records tool
	countRecordsTool := server.ServerTool{
		Tool: mcp.Tool{
			Name:        "count_records",
			Description: "Count records in a specified table with optional filtering",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"table_name": map[string]any{
						"type":        "string",
						"description": "Name of the table to count records from",
					},
					"date_filter": map[string]any{
						"type":        "string",
						"description": "Optional date filter (YYYY-MM-DD format)",
					},
					"status": map[string]any{
						"type":        "string",
						"description": "Optional status filter",
					},
				},
				Required: []string{"table_name"},
			},
		},
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleCountRecords(db, request)
		},
	}

	// Get user by ID tool
	getUserTool := server.ServerTool{
		Tool: mcp.Tool{
			Name:        "get_user_by_id",
			Description: "Retrieve user information by ID",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"user_id": map[string]any{
						"type":        "string",
						"description": "The unique identifier for the user",
					},
					"include_details": map[string]any{
						"type":        "boolean",
						"description": "Include detailed user information",
						"default":     false,
					},
				},
				Required: []string{"user_id"},
			},
		},
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleGetUserByID(db, request)
		},
	}

	// List active sessions tool
	listSessionsTool := server.ServerTool{
		Tool: mcp.Tool{
			Name:        "list_active_sessions",
			Description: "List all active user sessions",
			InputSchema: mcp.ToolInputSchema{
				Type: "object",
				Properties: map[string]any{
					"limit": map[string]any{
						"type":        "integer",
						"description": "Maximum number of sessions to return",
						"default":     100,
					},
					"user_type": map[string]any{
						"type":        "string",
						"description": "Filter by user type (admin, user, guest)",
						"default":     "user",
					},
				},
				Required: []string{},
			},
		},
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleListActiveSessions(db, request)
		},
	}

	mcpServer.AddTools(countRecordsTool, getUserTool, listSessionsTool)
	log.Println("Database tools registered")
}

// handleCountRecords handles the count_records tool call
func handleCountRecords(db *sql.DB, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid arguments format")
	}

	tableName, ok := args["table_name"].(string)
	if !ok {
		return nil, fmt.Errorf("table_name is required and must be a string")
	}

	// Build query with optional filters
	query := "SELECT COUNT(*) FROM " + tableName + " WHERE 1=1"
	var params []interface{}

	if dateFilter, ok := args["date_filter"].(string); ok && dateFilter != "" {
		query += " AND DATE(created_at) = ?"
		params = append(params, dateFilter)
	}

	if status, ok := args["status"].(string); ok && status != "" {
		query += " AND status = ?"
		params = append(params, status)
	}

	var count int
	err := db.QueryRow(query, params...).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to count records: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(fmt.Sprintf("%d", count)),
		},
	}, nil
}

// handleGetUserByID handles the get_user_by_id tool call
func handleGetUserByID(db *sql.DB, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid arguments format")
	}

	userID, ok := args["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("user_id is required and must be a string")
	}

	includeDetails := false
	if details, ok := args["include_details"].(bool); ok {
		includeDetails = details
	}

	// Build query based on include_details flag
	var query string
	if includeDetails {
		query = "SELECT id, name, email, status FROM users WHERE id = ?"
	} else {
		query = "SELECT id, name, email FROM users WHERE id = ?"
	}

	var id int
	var name, email string
	var status string

	if includeDetails {
		err := db.QueryRow(query, userID).Scan(&id, &name, &email, &status)
		if err != nil {
			if err == sql.ErrNoRows {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						mcp.NewTextContent("User not found"),
					},
				}, nil
			}
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		result := map[string]interface{}{
			"id":     id,
			"name":   name,
			"email":  email,
			"status": status,
		}

		// Convert to JSON string for response
		resultJSON, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	} else {
		err := db.QueryRow(query, userID).Scan(&id, &name, &email)
		if err != nil {
			if err == sql.ErrNoRows {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						mcp.NewTextContent("User not found"),
					},
				}, nil
			}
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		result := map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
		}

		// Convert to JSON string for response
		resultJSON, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}
}

// handleListActiveSessions handles the list_active_sessions tool call
func handleListActiveSessions(db *sql.DB, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		args = make(map[string]interface{})
	}

	limit := 100
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	// Note: user_type filter not implemented in this simplified version
	// userType := "user"
	// if ut, ok := args["user_type"].(string); ok && ut != "" {
	//     userType = ut
	// }

	// Query active sessions (status = 'active')
	query := "SELECT id, user_id, start_time FROM sessions WHERE status = ? ORDER BY start_time DESC"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := db.Query(query, "active")
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()

	var sessions []map[string]interface{}
	for rows.Next() {
		var id, userID int
		var startTime string

		err := rows.Scan(&id, &userID, &startTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		session := map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"start_time": startTime,
			"status":     "active",
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	// Convert to JSON string for response
	resultJSON, err := json.Marshal(sessions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(string(resultJSON)),
		},
	}, nil
}
