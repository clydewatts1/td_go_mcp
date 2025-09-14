package tools

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// SQLProcessor handles SQL template processing and parameter substitution
type SQLProcessor struct {
	tool ToolDefinition
}

// NewSQLProcessor creates a new SQL processor for the given tool
func NewSQLProcessor(tool ToolDefinition) *SQLProcessor {
	return &SQLProcessor{tool: tool}
}

// ProcessTemplate fills the SQL template with provided parameters
func (p *SQLProcessor) ProcessTemplate(params map[string]any) (string, error) {
	// Create template with custom functions
	tmpl, err := template.New(p.tool.Name).Funcs(template.FuncMap{
		"escape": escapeSQL,
	}).Parse(p.tool.SQLTemplate)
	if err != nil {
		return "", fmt.Errorf("invalid SQL template: %w", err)
	}

	// Merge parameters with defaults
	processedParams := make(map[string]any)

	// Set defaults first
	for name, param := range p.tool.Parameters {
		if param.Default != nil {
			processedParams[name] = param.Default
		}
	}

	// Override with provided parameters
	for key, value := range params {
		processedParams[key] = value
	}

	// Execute template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, processedParams)
	if err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	// Clean up the SQL (remove extra whitespace)
	sql := strings.TrimSpace(buf.String())
	sql = strings.ReplaceAll(sql, "\n\n", "\n")

	return sql, nil
}

// ValidateParameters checks if required parameters are provided and types match
func (p *SQLProcessor) ValidateParameters(params map[string]any) error {
	// Check required parameters
	for _, required := range p.tool.Required {
		if _, exists := params[required]; !exists {
			return fmt.Errorf("missing required parameter: %s", required)
		}
	}

	// Basic type validation
	for name, value := range params {
		if paramDef, exists := p.tool.Parameters[name]; exists {
			if !isValidType(value, paramDef.Type) {
				return fmt.Errorf("parameter %s: expected %s, got %T", name, paramDef.Type, value)
			}
		}
	}

	return nil
}

// isValidType performs basic type checking
func isValidType(value any, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "integer":
		switch value.(type) {
		case int, int32, int64, float64: // JSON numbers come as float64
			return true
		default:
			return false
		}
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "number":
		switch value.(type) {
		case int, int32, int64, float32, float64:
			return true
		default:
			return false
		}
	default:
		return true // Unknown types pass validation
	}
}

// escapeSQL provides basic SQL escaping (simple quote doubling)
func escapeSQL(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
