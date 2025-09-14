package tools

import (
	"testing"
)

func TestLoadToolsFromDirectory(t *testing.T) {
	// Test loading from tools directory
	tools, err := LoadToolsFromDirectory("../../tools")
	if err != nil {
		t.Fatalf("Failed to load tools: %v", err)
	}

	if len(tools) == 0 {
		t.Fatal("Expected to load some tools, got none")
	}

	// Check that get_user_by_id tool exists
	var userTool *ToolDefinition
	for i := range tools {
		if tools[i].Name == "get_user_by_id" {
			userTool = &tools[i]
			break
		}
	}

	if userTool == nil {
		t.Fatal("Expected to find get_user_by_id tool")
	}

	if userTool.Description == "" {
		t.Error("Tool description should not be empty")
	}

	if userTool.SQLTemplate == "" {
		t.Error("SQL template should not be empty")
	}

	if len(userTool.Required) == 0 {
		t.Error("Expected at least one required parameter")
	}
}

func TestSQLProcessor(t *testing.T) {
	tool := ToolDefinition{
		Name:        "test_tool",
		SQLTemplate: "SELECT * FROM users WHERE id = '{{.user_id}}' {{if .active}}AND active = 1{{end}}",
		Parameters: map[string]Parameter{
			"user_id": {Type: "string", Description: "User ID"},
			"active":  {Type: "boolean", Description: "Filter active users", Default: false},
		},
		Required: []string{"user_id"},
	}

	processor := NewSQLProcessor(tool)

	// Test with required parameter
	params := map[string]any{"user_id": "123"}
	sql, err := processor.ProcessTemplate(params)
	if err != nil {
		t.Fatalf("Template processing failed: %v", err)
	}

	expected := "SELECT * FROM users WHERE id = '123'"
	if sql != expected {
		t.Errorf("Expected SQL: %s, got: %s", expected, sql)
	}

	// Test with optional parameter
	params["active"] = true
	sql, err = processor.ProcessTemplate(params)
	if err != nil {
		t.Fatalf("Template processing failed: %v", err)
	}

	expectedWithActive := "SELECT * FROM users WHERE id = '123' AND active = 1"
	if sql != expectedWithActive {
		t.Errorf("Expected SQL: %s, got: %s", expectedWithActive, sql)
	}

	// Test validation - missing required parameter
	err = processor.ValidateParameters(map[string]any{})
	if err == nil {
		t.Error("Expected validation error for missing required parameter")
	}
}
