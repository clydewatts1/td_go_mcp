package main

import (
	"context"
	"fmt"

	"td_go_mcp/internal/tools"

	mcp "github.com/metoro-io/mcp-golang"
	"golang.org/x/exp/slog"
)

func createPromptHandler(promptDef tools.PromptDefinition) func(context.Context, map[string]any) (*mcp.PromptResponse, error) {
	return func(ctx context.Context, args map[string]any) (*mcp.PromptResponse, error) {
		slog.Info("[prompt_handlers] Handling prompt call", "prompt", promptDef.Name)

		responseText := fmt.Sprintf("Received prompt '%s' with arguments: %v", promptDef.Name, args)

		// Correctly create a prompt response.
		// NewPromptResponse takes a description string and a variadic of *PromptMessage.
		// NewPromptMessage takes content and a role.
		response := mcp.NewPromptResponse(
			promptDef.Description,
			mcp.NewPromptMessage(mcp.NewTextContent(responseText), mcp.RoleAssistant),
		)

		return response, nil
	}
}
