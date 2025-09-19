package main

import (
	"log"
	"td_go_mcp/internal/db"
	"td_go_mcp/internal/tools"
)

var (
	loadedTools   []tools.ToolDefinition
	loadedPrompts []tools.PromptDefinition
	processors    map[string]*tools.SQLProcessor
	database      *db.DB
)

func init() {
	var err error
	loadedTools, err = tools.LoadToolsFromDirectory("tools")
	if err != nil {
		log.Printf("Error loading tools: %v", err)
		loadedTools = []tools.ToolDefinition{} // Continue with empty tools
	}

	processors = make(map[string]*tools.SQLProcessor)
	for i := range loadedTools {
		processors[loadedTools[i].Name] = tools.NewSQLProcessor(loadedTools[i])
	}

	// Load prompts from YAML files
	loadedPrompts, err = tools.LoadPromptsFromDirectory("tools")
	if err != nil {
		log.Printf("Error loading prompts: %v", err)
		loadedPrompts = []tools.PromptDefinition{} // Continue with empty prompts
	}

	log.Printf("Loaded %d tools and %d prompts from YAML files", len(loadedTools), len(loadedPrompts))

	// Initialize database connection
	dbConfig := db.LoadConfig()
	database, err = db.Connect(dbConfig)
	if err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
		log.Printf("Continuing without database - SQL preview mode only")
		database = nil
	} else {
		log.Printf("Database connection established successfully")
	}
}
