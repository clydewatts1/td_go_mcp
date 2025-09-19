package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ToolDefinition represents a tool loaded from YAML
type ToolDefinition struct {
	Name        string               `yaml:"name" json:"name"`
	Description string               `yaml:"description" json:"description"`
	Parameters  map[string]Parameter `yaml:"parameters" json:"parameters"`
	ReturnType  string               `yaml:"return_type" json:"return_type"`
	SQLTemplate string               `yaml:"sql_template" json:"sql_template"`
	Required    []string             `yaml:"required" json:"required"`
}

// PromptDefinition represents a prompt loaded from YAML
type PromptDefinition struct {
	Type        string               `yaml:"type" json:"type"`
	Name        string               `yaml:"name" json:"name"`
	Prompt      string               `yaml:"prompt" json:"prompt"`
	Description string               `yaml:"description,omitempty" json:"description,omitempty"`
	Parameters  map[string]Parameter `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

// Parameter defines input parameter schema
type Parameter struct {
	Type        string `yaml:"type" json:"type"`
	Description string `yaml:"description" json:"description"`
	Default     any    `yaml:"default,omitempty" json:"default,omitempty"`
}

// LoadToolsFromDirectory loads all YAML files from tools/ directory
func LoadToolsFromDirectory(dir string) ([]ToolDefinition, error) {
	var tools []ToolDefinition

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return tools, nil // No tools directory, return empty
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			// Skip prompt files when loading tools
			if isPromptFile(path) {
				return nil
			}
			
			tool, err := loadToolFromFile(path)
			if err != nil {
				return fmt.Errorf("error loading %s: %w", path, err)
			}
			tools = append(tools, tool)
		}
		return nil
	})

	return tools, err
}

func loadToolFromFile(filepath string) (ToolDefinition, error) {
	var tool ToolDefinition

	data, err := os.ReadFile(filepath)
	if err != nil {
		return tool, err
	}

	err = yaml.Unmarshal(data, &tool)
	if err != nil {
		return tool, err
	}

	// Validate required fields
	if tool.Name == "" {
		return tool, fmt.Errorf("tool name is required")
	}
	if tool.SQLTemplate == "" {
		return tool, fmt.Errorf("sql_template is required")
	}

	return tool, nil
}

// ToMCPTool converts ToolDefinition to MCP Tool format
func (td *ToolDefinition) ToMCPTool() map[string]any {
	properties := make(map[string]any)
	for name, param := range td.Parameters {
		properties[name] = map[string]any{
			"type":        param.Type,
			"description": param.Description,
		}
		if param.Default != nil {
			properties[name].(map[string]any)["default"] = param.Default
		}
	}

	return map[string]any{
		"name":        td.Name,
		"description": td.Description,
		"inputSchema": map[string]any{
			"type":       "object",
			"properties": properties,
			"required":   td.Required,
		},
	}
}

// LoadPromptsFromDirectory loads all prompt YAML files from tools/ directory
func LoadPromptsFromDirectory(dir string) ([]PromptDefinition, error) {
	var prompts []PromptDefinition

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return prompts, nil // No tools directory, return empty
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			// First check if this is a prompt file
			if isPromptFile(path) {
				prompt, err := loadPromptFromFile(path)
				if err != nil {
					return fmt.Errorf("error loading prompt %s: %w", path, err)
				}
				prompts = append(prompts, prompt)
			}
		}
		return nil
	})

	return prompts, err
}

func isPromptFile(filepath string) bool {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return false
	}

	// Quick check if this is a prompt file by looking for type: prompt
	var quickCheck struct {
		Type string `yaml:"type"`
	}
	
	if err := yaml.Unmarshal(data, &quickCheck); err != nil {
		return false
	}
	
	return quickCheck.Type == "prompt"
}

func loadPromptFromFile(filepath string) (PromptDefinition, error) {
	var prompt PromptDefinition

	data, err := os.ReadFile(filepath)
	if err != nil {
		return prompt, err
	}

	err = yaml.Unmarshal(data, &prompt)
	if err != nil {
		return prompt, err
	}

	// Validate required fields
	if prompt.Type != "prompt" {
		return prompt, fmt.Errorf("type must be 'prompt'")
	}
	if prompt.Name == "" {
		return prompt, fmt.Errorf("name is required")
	}
	if prompt.Prompt == "" {
		return prompt, fmt.Errorf("prompt text is required")
	}

	return prompt, nil
}
