package main

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/exp/slog"

	"td_go_mcp/internal/tools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Place prompt handler functions here, e.g. createPromptHandler(...)

func addPromptToServer(mcpServer *server.MCPServer, promptDef tools.PromptDefinition) {
	slog.Debug("[prompt_handlers] Adding prompt to server", "prompt", promptDef.Name)
	mcpPrompt := convertPromptDefinition(promptDef)
	handler := createPromptHandler(promptDef)
	mcpServer.AddPrompt(mcpPrompt, handler)
	slog.Info("[prompt_handlers] Registered prompt", "prompt", promptDef.Name)
}

func convertPromptDefinition(promptDef tools.PromptDefinition) mcp.Prompt {
	slog.Debug("[prompt_handlers] Converting prompt definition", "prompt", promptDef.Name)
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
			slog.Debug("[prompt_handlers] Adding argument to prompt", "prompt", promptDef.Name, "param", paramName)
			opts = append(opts, mcp.WithArgument(paramName, argumentOpts...))
		}
	}
	slog.Debug("[prompt_handlers] Prompt definition converted", "prompt", promptDef.Name)
	return mcp.NewPrompt(promptDef.Name, opts...)
}

func createPromptHandler(promptDef tools.PromptDefinition) func(context.Context, mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		slog.Info("[prompt_handlers] Handling prompt request", "prompt", promptDef.Name)
		args := req.Params.Arguments
		if args == nil {
			slog.Debug("[prompt_handlers] No arguments provided, initializing empty map", "prompt", promptDef.Name)
			args = make(map[string]string)
		}
		slog.Debug("[prompt_handlers] Prompt arguments", "prompt", promptDef.Name, "args", args)
		processedPrompt := promptDef.Prompt
		for paramName, paramValue := range args {
			placeholder := fmt.Sprintf("{{%s}}", paramName)
			processedPrompt = strings.ReplaceAll(processedPrompt, placeholder, paramValue)
			slog.Debug("[prompt_handlers] Replaced argument in prompt", "param", paramName, "value", paramValue)
		}
		for paramName, param := range promptDef.Parameters {
			placeholder := fmt.Sprintf("{{%s}}", paramName)
			if strings.Contains(processedPrompt, placeholder) {
				if param.Default != nil {
					if defaultStr, ok := param.Default.(string); ok {
						processedPrompt = strings.ReplaceAll(processedPrompt, placeholder, defaultStr)
						slog.Debug("[prompt_handlers] Used default value for prompt parameter", "param", paramName, "default", defaultStr)
					}
				} else {
					processedPrompt = strings.ReplaceAll(processedPrompt, placeholder, "")
					slog.Debug("[prompt_handlers] No value for prompt parameter, replaced with empty string", "param", paramName)
				}
			}
		}
		messages := []mcp.PromptMessage{
			{
				Role:    mcp.RoleUser,
				Content: mcp.TextContent{Type: "text", Text: processedPrompt},
			},
		}
		slog.Debug("[prompt_handlers] Returning prompt result", "prompt", promptDef.Name, "content", processedPrompt)
		return mcp.NewGetPromptResult(promptDef.Description, messages), nil
	}
}
