package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"td_go_mcp/internal/tools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Place prompt handler functions here, e.g. createPromptHandler(...)

func addPromptToServer(mcpServer *server.MCPServer, promptDef tools.PromptDefinition) {
	mcpPrompt := convertPromptDefinition(promptDef)
	handler := createPromptHandler(promptDef)
	mcpServer.AddPrompt(mcpPrompt, handler)
	log.Printf("Registered prompt: %s", promptDef.Name)
}

func convertPromptDefinition(promptDef tools.PromptDefinition) mcp.Prompt {
	opts := []mcp.PromptOption{
		mcp.WithPromptDescription(promptDef.Description),
	}
	if len(promptDef.Parameters) > 0 {
		for paramName, param := range promptDef.Parameters {
			argumentOpts := []mcp.ArgumentOption{
				mcp.ArgumentDescription(param.Description),
			}
			if param.Default == nil {
				argumentOpts = append(argumentOpts, mcp.RequiredArgument())
			}
			opts = append(opts, mcp.WithArgument(paramName, argumentOpts...))
		}
	}
	return mcp.NewPrompt(promptDef.Name, opts...)
}

func createPromptHandler(promptDef tools.PromptDefinition) func(context.Context, mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		log.Printf("Handling prompt request: %s", promptDef.Name)
		args := req.Params.Arguments
		if args == nil {
			args = make(map[string]string)
		}
		processedPrompt := promptDef.Prompt
		for paramName, paramValue := range args {
			placeholder := fmt.Sprintf("{{%s}}", paramName)
			processedPrompt = strings.ReplaceAll(processedPrompt, placeholder, paramValue)
		}
		for paramName, param := range promptDef.Parameters {
			placeholder := fmt.Sprintf("{{%s}}", paramName)
			if strings.Contains(processedPrompt, placeholder) {
				if param.Default != nil {
					if defaultStr, ok := param.Default.(string); ok {
						processedPrompt = strings.ReplaceAll(processedPrompt, placeholder, defaultStr)
					}
				} else {
					processedPrompt = strings.ReplaceAll(processedPrompt, placeholder, "")
				}
			}
		}
		messages := []mcp.PromptMessage{
			{
				Role:    mcp.RoleUser,
				Content: mcp.TextContent{Type: "text", Text: processedPrompt},
			},
		}
		return mcp.NewGetPromptResult(promptDef.Description, messages), nil
	}
}
